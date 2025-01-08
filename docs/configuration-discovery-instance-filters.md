# Configuration: instance filters

At times, you may want to exclude particular hosts from discovery. For example, when utility nodes like gh-ost, Debzium, Tungsten, or other nodes that don't require monitoring are part of a large, complex environment.

To skip a node during the Discovery process, specify the host name using a regex pattern. If the host name matches the pattern, Orchestrator skips the node and continues the process:

```json
{
  "DiscoveryIgnoreHostnameFilters": [
    "a_host_i_want_to_ignore[.]example[.]com:5000",
    ".*[.]example2.com:5000",
    "192.168.0.1:6000"
  ]
}
```

You can also use the regex patterns to prevent auto-discovering new replicas. These replicas may be on unreachable servers due to firewall rules:

```json
{
  "DiscoveryIgnoreReplicaHostnameFilters": [
    "a_host_i_want_to_ignore[.]example[.]com",
    ".*[.]ignore_all_hosts_from_this_domain[.]example[.]com",
    "a_host_with_extra_port_i_want_to_ignore[.]example[.]com:3307"
  ]
}
```

Use the regexp patterns to prevent auto-discovering a master:

```json
{
  "DiscoveryIgnoreMasterHostnameFilters": [
    "a_host_i_want_to_ignore[.]example[.]com:5000",
    ".*[.]example2.com:5000",
    "192.168.0.1:6000"
  ]
}
```

It is also possible to specify regexp filters used to check the replication user name used by the replica. When Orchestrator detects that a given instance uses the replication user matching one of the patterns, such instance is skipped during the discovery:

```json
{
    "DiscoveryIgnoreReplicationUsernameFilters": [
            "replication_user",
            "rpl_[0-9].*"
    ]
}
```
For the Orchestrator to discover the replication topology it is enough to point it to any node. The Orchestrator will detect possible replication sources of the given node and its replicas. Then it will continue with the discovery using the detected source/replicas and build a picture of the whole replication graph.

Let's assume that we have the following replication channel:

node_1 -> node_2(repl_user_2) -> node_3(repl_user_3)

node_1 : source\
node_2 : intermediate source that uses `repl_user_2` to replicate from node_1\
node_3 : replica that uses `repl_user_3` to replicate from node_2

and
```json
{
    "DiscoveryIgnoreReplicationUsernameFilters": [
            "repl_user_2"
    ]
}
```

Let's consider following discovery cases:

## Discovery process through node_1
Orchestrator learns that node_2 is the replica during the examination of node_1. If the replication user used by node_2 is known at this point, node_2 will be filtered out immediately and never queried.\
When Orchestrator examines node_1 and finds node_2 as a replica, two scenarios can occur:
1. Using Show Slave Hosts (`DiscoverByShowSlaveHosts=true`):\
Orchestrator can skip node_2 immediately if:
    1. node_1 runs with `--show-replica-auth-info=1`
    2. node_2 sets `report_user=<replication_user>`

   Without these settings, Orchestrator must query node_2 directly. To avoid this overhead, configure node_1 and node_2 properly.
2. Using Process List (`DiscoverByShowSlaveHosts=false`):\
Orchestrator checks `information_schema.processlist`. This table shows the replica's replication user. Node_2 gets skipped without direct queries.

As a result, only node_1 is discovered.

## Discovery process through node_2
When starting discovery from node_2, Orchestrator checks the replication user name. If the name matches the filter pattern, discovery stops. No further nodes get discovered.

## Discovery process through node_3
When starting with node_3, Orchestrator finds node_2 as the source. The replication user for node_2 is unknown initially. Orchestrator schedules node_2 for discovery.\
Upon examining node_2:
* Finds matching replication user and skips further discovery
* Node_1 remains undiscovered

Only node_3 stays in the discovered set.

## Logs
Orchestrator logs such discovery skips in its error log. Logging is enabled by default and controlled by the following setting:

```json
{
    "EnableDiscoveryFiltersLogs": false
}
```
As for now, logging control applies only to `DiscoveryIgnoreReplicationUsernameFilters`.