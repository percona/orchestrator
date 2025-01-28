package inst

import (
	"sync"
	"time"

	"github.com/openark/golib/log"
	"github.com/openark/orchestrator/go/config"
	"github.com/rcrowley/go-metrics"
)

// The behavior depends on settings:
// 1. DeadInstanceDiscoveryMaxConcurrency > 0 and DeadInstancePollSecondsMultiplyFactor > 1:
//    The separate discovery queue for dead instances is created and dead instances
// 	  are checked by dedicated pool of go workers
//    and the instance is checked with exponential backoff mechanism time
// 2. DeadInstanceDiscoveryMaxConcurrency = 0 and DeadInstancePollSecondsMultiplyFactor > 1:
//    No separate discovery queue for dead instances is created and dead instances
//    are checked by the same pool of go workers as healthy instances, however
//    an exponential backoff mechanism is applied for dead instances
// 3. DeadInstanceDiscoveryMaxConcurrency > 0 and DeadInstancePollSecondsMultiplyFactor = 1:
//    The separate discovery queue for dead instances is created and dead instances
//    are checked by dedicated pool of go workers. No exponential backoff mechanism
//    is applied for dead instances
// 4. DeadInstanceDiscoveryMaxConcurrency = 0 and DeadInstancePollSecondsMultiplyFactor = 1:
//    No separate discovery queue for dead instances, no dedicated go workers,
//    no backoff mechanism. This is the default working mode.
//
// We register a dead instance always. It shouldn't be a big overhead,
// and we will get the info about the dead instances count.

type deadInstance struct {
	DelayFactor   float32
	NextCheckTime time.Time
	TryCnt        int
}

type deadInstancesFilter struct {
	deadInstances      map[InstanceKey]deadInstance
	deadInstancesMutex sync.RWMutex
}

var DeadInstancesFilter deadInstancesFilter

var deadInstancesCounter = metrics.NewCounter()

func init() {
	metrics.Register("discoveries.dead_instances", deadInstancesCounter)
	DeadInstancesFilter.deadInstances = make(map[InstanceKey]deadInstance)
	DeadInstancesFilter.deadInstancesMutex = sync.RWMutex{}
}

// RegisterInstance registers a given instance in a dead instances cache.
// Once the instance is registered its discovery can be delayed with exponential
// backoff mechanism according to DeadInstancePollSecondsMultiplyFactor value.
// During the registration, next desired check time is calculated basing on
// the current delay factor, DeadInstancePollSecondsMultiplyFactor and
// DeadInstancePollSecondsMax parameters.
func (f *deadInstancesFilter) RegisterInstance(instanceKey *InstanceKey) {
	delayFactor := float32(1)
	previousTry := 0

	f.deadInstancesMutex.Lock()
	defer f.deadInstancesMutex.Unlock()

	instance, exists := f.deadInstances[*instanceKey]
	if exists {
		delayFactor = config.Config.DeadInstancePollSecondsMultiplyFactor * instance.DelayFactor
		previousTry = instance.TryCnt
	} else {
		deadInstancesCounter.Inc(1)
	}

	maxDelay := time.Duration(config.Config.DeadInstancePollSecondsMax) * time.Second
	currentDelay := time.Duration(delayFactor*float32(config.Config.InstancePollSeconds)) * time.Second

	// needed only for the debug log below
	delayFactorTmp := delayFactor

	if currentDelay > maxDelay {
		// saturation
		currentDelay = maxDelay
		delayFactor = instance.DelayFactor // back to previous one
	}
	nextCheck := time.Now().Add(currentDelay)

	instance = deadInstance{
		DelayFactor:   delayFactor,
		NextCheckTime: nextCheck,
		TryCnt:        previousTry + 1,
	}
	f.deadInstances[*instanceKey] = instance

	if config.Config.DeadInstanceDiscoveryLogsEnabled {
		log.Debugf("Dead instance registered %v:%v. Iteration: %v. Current delay factor: %v (next check in %v (on %v))",
			instanceKey.Hostname, instanceKey.Port, instance.TryCnt, delayFactorTmp, currentDelay, instance.NextCheckTime)
	}
}

// UnregisterInstace removes the given instance from dead instances cache.
func (f *deadInstancesFilter) UnregisterInstance(instanceKey *InstanceKey) {
	f.deadInstancesMutex.Lock()
	defer f.deadInstancesMutex.Unlock()

	instance, exists := f.deadInstances[*instanceKey]
	if exists {
		if config.Config.DeadInstanceDiscoveryLogsEnabled {
			log.Debugf("Dead instance unregistered: %v:%v after iteration: %v",
				instanceKey.Hostname, instanceKey.Port, instance.TryCnt)
		}
		deadInstancesCounter.Dec(1)
		delete(f.deadInstances, *instanceKey)
	}
}

// InstanceRecheckNeeded checks if a given instance is registered in a dead instances
// cache and if it is, is it time to rediscover it.
// It returns two boolean values:
// - The first boolean indicates if the instance is registered.
// - The second boolean, indicates if it is time to rediscover the node.
func (f *deadInstancesFilter) InstanceRecheckNeeded(instanceKey *InstanceKey) (bool, bool) {
	f.deadInstancesMutex.RLock()
	defer f.deadInstancesMutex.RUnlock()

	instance, exists := f.deadInstances[*instanceKey]

	if !exists {
		return exists, false
	}

	if instance.NextCheckTime.After(time.Now()) {
		// recheck time still in the future
		return exists, false
	}

	if config.Config.DeadInstanceDiscoveryLogsEnabled {
		log.Debugf("Dead instance recheck: %v:%v. Iteration: %v",
			instanceKey.Hostname, instanceKey.Port, instance.TryCnt)
	}
	return exists, true
}
