# Test the error handling when the master fails to take over gracefully and replication is stopped

This error occurred in production when two takeovers were executed in a row.
The first take over was successful, but the second one failed.

This error is reproducable by stopping replication on the replica that
is supposed to take over. Another way to reproduce this error is to
run two takeovers one immediately after the other:

```sh
orchestrator-client -c graceful-master-takeover -i 127.0.0.1:10111 -d 127.0.0.1:10112
orchestrator-client -c graceful-master-takeover -i 127.0.0.1:10112 -d 127.0.0.1:10111
```

In the end the topology will be in a partially failed state, with
replication stopped for replica `127.0.0.1:10112`, and the other two
replicas placed behind it. Though, `master` will still be writable,
and `PostUnsuccessfulGracefulTakeoverProcesses` hooks will be executed
to help the cluster recover.
