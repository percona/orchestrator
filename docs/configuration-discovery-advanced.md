# Configuration: advanced discovery

The Orchestrator uses an internal queue to manage instances for discovery. When an instance is ready for discovery, it gets added to the queue. Discovery workers process the queue. The `DiscoveryMaxConcurrency` setting in a configuration file controls the number of workers. This setting determines how many discoveries can happen in parallel.

The Orchestrator uses this mechanism to periodically monitor all instances. `InstancePollSeconds` configuration parameter says how often the Orchestrator should refresh the information.

When there is a lot of inaccessible or unhealthy instances, the Orchestrator may lose the proper view of the cluster and be late with needed recovery actions. This is because discoveries of such instances may take a long time and finish with failure anyway, at the same time consuming workers from the discovery workers pool. Healthy instances wait in the queue and they are not checked in a timely manner.

To avoid this, Orchestrator can be configured to maintain a separate discovery queue for unhealthy instances. This queue is processed by a separate pool of workers. Additionally, an exponential time backoff mechanism can be applied for rechecking such instances.

Configuration example:
```json
{
  "DeadInstanceDiscoveryMaxConcurrency": 100,
  "DeadInstancePollSecondsMultiplyFactor": 1.5,
  "DeadInstancePollSecondsMax": 60,
  "DeadInstanceDiscoveryLogsEnabled": true
}
```

`DeadInstanceDiscoveryMaxConcurrency` (default: 0) - Determines the number of discovery workers dedicated to dead instances. If this pool size is grater than 0, the Orchestrator maintains a separate queue for dead instances.

`DeadInstancePollSecondsMultiplyFactor` (default: 1) - Floating point number, allowed values are >= 1. Determines how aggressive the backoff mechanism is. By default, when `DeadInstancePollSecondsMultiplyFactor = 1`, the instance is checked every `InstancePollSeconds` seconds. If the parameter value is greater than 1, every consecutive try `n` is done after the period calculated according to the formula:

dT(n) = InstancePollSeconds * DeadInstancePollSecondsMultiplyFactor ^ (n-1)

Example:

Let's use `D` as `DeadInstancePollSecondsMultiplyFactor`

f(1) = 1\
f(2) = f(1) * D\
f(3) = f(2) * D\
f(4) = f(3) * D

That means:

f(4) = 1 * D * D * D = D^3

or in other words

f(n) = DeadInstancePollSecondsMultiplyFactor ^ (n-1)

so:

dT(n) = InstancePollSeconds * f(n)\
dT(n) = InstancePollSeconds * DeadInstancePollSecondsMultiplyFactor ^ (n-1)

Note that `DeadInstanceDiscoveryMaxConcurrency` controls if the separate pool of discovery workers is created but has no impact on the backoff mechanism controlled by `DeadInstancePollSecondsMultiplyFactor`. It has the following implications:

1. `DeadInstanceDiscoveryMaxConcurrency > 0` and `DeadInstancePollSecondsMultiplyFactor > 1`:\
The separate discovery queue for dead instances is created, and dead instances are checked by a dedicated pool of go workers, and the instance is checked with exponential backoff mechanism time
2. `DeadInstanceDiscoveryMaxConcurrency = 0` and `DeadInstancePollSecondsMultiplyFactor > 1`:\
No separate discovery queue for dead instances is created, and dead instances are checked by the same pool of go workers as healthy instances however, an exponential backoff mechanism is applied for dead instances
3. `DeadInstanceDiscoveryMaxConcurrency > 0` and `DeadInstancePollSecondsMultiplyFactor = 1`:\
A separate discovery queue for dead instances is created, and dead instances are checked by a dedicated pool of go workers. No exponential backoff mechanism is applied for dead instances
4. `DeadInstanceDiscoveryMaxConcurrency = 0` and `DeadInstancePollSecondsMultiplyFactor = 1`:\
There is no separate discovery queue for dead instances, no dedicated go workers, no backoff mechanism. This is the default working mode.

`DeadInstancePollSecondsMax` (default: 300) - Controls the maximum time for backoff mechanism. If the backoff calculation goes beyond this value, it is considered as saturated and stays at `DeadInstancePollSecondsMax`

## Diagnostics
Orchestrator provides `debug/metrics` web endpoint for diagnostics.

`discoveries.dead_instances` - provides the number of instances currently registered as dead.\
`discoveries.dead_instances_queue_length` - provides the current length of the queue dedicate for dead instances. Note this is valid only when `DeadInstanceDiscoveryMaxConcurrency > 0`, so when a separate queue is used. In other cases it is always zero.

Other diagnostics endpoints:

`api/discovery-queue-metrics-raw/:seconds` - provides the raw metrics for a given time for the `DEFAULT` discovery queue.\
`api/discovery-queue-metrics-raw/:queue/:seconds` - provides the raw metrics for a given time for the supplied (`DEFAULT` or `DEADINSTANCES`) discovery queue.\
`discovery-queue-metrics-aggregated/:seconds` - provides aggregated metrics for a given time for the `DEFAULT` discovery queue.\
`discovery-queue-metrics-aggregated/:queue/:seconds` - provides aggregated metrics for a given time for the supplied (`DEFAULT` or `DEADINSTANCES`) discovery queue.


Note that `DEADINSTANCES` queue is available only if `DeadInstanceDiscoveryMaxConcurrency > 0`

## Logging
Logging of dead instances discovery process is controlled vial `DeadInstanceDiscoveryLogsEnabled` bool parameter. It is disabled by default.