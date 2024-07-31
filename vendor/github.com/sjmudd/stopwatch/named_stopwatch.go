/*
Copyright (c) 2016, Simon J Mudd
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package stopwatch implements simple stopwatch functionality
package stopwatch

import (
	"fmt"
	"sync"
	"time"
)

type logger interface {
	Warning(args ...interface{})
}

// NamedStopwatch holds a map of string named stopwatches. Intended
// to be used when several Stopwatches are being used at once, and
// easy to use as they are name based.
type NamedStopwatch struct {
	sync.RWMutex
	stopwatches map[string](*Stopwatch)
	logger      logger
}

// NewNamedStopwatch creates an empty Stopwatch list
func NewNamedStopwatch() *NamedStopwatch {
	return new(NamedStopwatch)
}

// Add adds a single Stopwatch name with the given name.
func (ns *NamedStopwatch) Add(name string) error {
	ns.Lock()
	defer ns.Unlock()

	return ns.add(name)
}

// Add adds a single Stopwatch name with the given name.
// The caller is assumed to have locked the structure.
func (ns *NamedStopwatch) add(name string) error {
	if ns.stopwatches == nil {
		// create structure
		ns.stopwatches = make(map[string](*Stopwatch))
	} else {
		// check for existing name
		if _, ok := ns.stopwatches[name]; ok {
			return fmt.Errorf("NamedStopwatch.add() Stopwatch name %q already exists", name)
		}
	}
	ns.stopwatches[name] = New(nil)

	return nil
}

// AddMany adds several named stopwatches in one go
func (ns *NamedStopwatch) AddMany(names []string) error {
	ns.Lock()
	defer ns.Unlock()

	for _, name := range names {
		if err := ns.add(name); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a Stopwatch with the given name (if it exists)
func (ns *NamedStopwatch) Delete(name string) {
	ns.Lock()
	defer ns.Unlock()

	if ns.stopwatches == nil {
		return
	}

	delete(ns.stopwatches, name) // check if it exists in case the user did the wrong thing
}

// Exists returns true if the NamedStopwatch exists
func (ns *NamedStopwatch) Exists(name string) bool {
	ns.RLock()
	defer ns.RUnlock()

	if ns == nil {
		return false
	}

	_, found := ns.stopwatches[name]

	return found
}

// Start starts a NamedStopwatch if it exists
func (ns *NamedStopwatch) Start(name string) {
	if ns == nil {
		return // if we're not using stopwatches we just do nothing
	}
	ns.Lock()
	defer ns.Unlock()

	ns.start(name)
}

// start starts a NamedStopwatch if it exists. The structure is expected to be locked.
func (ns *NamedStopwatch) start(name string) {
	if ns == nil {
		return
	}
	if s, ok := ns.stopwatches[name]; ok {
		s.Start()
	}
}

// StartMany allows you to start several stopwatches in one go
func (ns *NamedStopwatch) StartMany(names []string) {
	if ns == nil {
		return // if we're not using stopwatches we just do nothing
	}
	ns.Lock()
	defer ns.Unlock()

	for _, name := range names {
		ns.start(name)
	}
}

// Stop stops a NamedStopwatch if it exists
func (ns *NamedStopwatch) Stop(name string) {
	if ns == nil {
		return // if we're not using stopwatches we just do nothing
	}
	ns.Lock()
	defer ns.Unlock()

	ns.stop(name)
}

// stop stops a NamedStopwatch if it exists and expects the structure to be locked.
func (ns *NamedStopwatch) stop(name string) {
	if ns == nil {
		return
	}
	if s, ok := ns.stopwatches[name]; ok {
		if s.IsRunning() {
			s.Stop()
		} else if ns.logger != nil {
			ns.logger.Warning("NamedStopwatch.Stop(%q) IsRunning is false\n", name)
		}
	}
}

// StopMany allows you to stop several stopwatches in one go
func (ns *NamedStopwatch) StopMany(names []string) {
	if ns == nil {
		return // if we're not using stopwatches we just do nothing
	}
	ns.Lock()
	defer ns.Unlock()

	for _, name := range names {
		ns.stop(name)
	}
}

// Reset resets a NamedStopwatch if it exists
func (ns *NamedStopwatch) Reset(name string) {
	if ns == nil {
		return // if we're not using stopwatches we just do nothing
	}
	ns.Lock()
	defer ns.Unlock()

	if ns == nil {
		return
	}
	if s, ok := ns.stopwatches[name]; ok {
		s.Reset()
	}
}

// Keys returns the known names of Stopwatches
func (ns *NamedStopwatch) Keys() []string {
	if ns == nil {
		return nil
	}

	ns.RLock()
	defer ns.RUnlock()

	keys := []string{}
	for k := range ns.stopwatches {
		keys = append(keys, k)
	}
	return keys
}

// Elapsed returns the elapsed time.Duration of the named stopwatch if it exists or 0
func (ns *NamedStopwatch) Elapsed(name string) time.Duration {
	if ns == nil {
		return time.Duration(0)
	}
	ns.RLock()
	defer ns.RUnlock()

	if s, ok := ns.stopwatches[name]; ok {
		return s.Elapsed()
	}
	return time.Duration(0)
}

// ElapsedSeconds returns the elapsed time in seconds of the named
// stopwatch if it exists or 0.
func (ns *NamedStopwatch) ElapsedSeconds(name string) float64 {
	if ns == nil {
		return float64(0)
	}
	ns.RLock()
	defer ns.RUnlock()

	if s, ok := ns.stopwatches[name]; ok {
		return s.ElapsedSeconds()
	}
	return float64(0)
}

// ElapsedMilliSeconds returns the elapsed time in milliseconds of
// the named stopwatch if it exists or 0.
func (ns *NamedStopwatch) ElapsedMilliSeconds(name string) float64 {
	if ns == nil {
		return float64(0)
	}
	ns.RLock()
	defer ns.RUnlock()

	if s, ok := ns.stopwatches[name]; ok {
		return s.ElapsedMilliSeconds()
	}
	return float64(0)
}

// AddElapsedSince adds the duration since the reference time to the given named stopwatch.
func (ns *NamedStopwatch) AddElapsedSince(name string, t time.Time) {
	if ns == nil {
		return
	}
	ns.Lock()
	defer ns.Unlock()

	if s, ok := ns.stopwatches[name]; ok {
		s.AddElapsedSince(t)
	}
}

// SetLogger is used to set the logger interface
func (ns *NamedStopwatch) SetLogger(logger logger) {
	ns.logger = logger
}
