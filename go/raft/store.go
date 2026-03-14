package orcraft

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/openark/golib/log"

	"github.com/hashicorp/raft"
)

type Store struct {
	raftDir       string
	raftBind      string
	raftAdvertise string

	raft *raft.Raft // The consensus mechanism

	applier                CommandApplier
	snapshotCreatorApplier SnapshotCreatorApplier
}

type storeCommand struct {
	Op    string `json:"op,omitempty"`
	Value []byte `json:"value,omitempty"`
}

// NewStore inits and returns a new store
func NewStore(raftDir string, raftBind string, raftAdvertise string, applier CommandApplier, snapshotCreatorApplier SnapshotCreatorApplier) *Store {
	return &Store{
		raftDir:                raftDir,
		raftBind:               raftBind,
		raftAdvertise:          raftAdvertise,
		applier:                applier,
		snapshotCreatorApplier: snapshotCreatorApplier,
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
func (store *Store) Open(peerNodes []string) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.SnapshotThreshold = 1
	config.SnapshotInterval = snapshotInterval
	config.ShutdownOnRemove = false
	config.LocalID = raft.ServerID(store.raftAdvertise)

	// Setup Raft communication.
	advertise, err := net.ResolveTCPAddr("tcp", store.raftAdvertise)
	if err != nil {
		return err
	}
	log.Debugf("raft: advertise=%+v", advertise)

	transport, err := raft.NewTCPTransport(store.raftBind, advertise, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	log.Debugf("raft: transport=%+v", transport)

	if _, err := os.Stat(store.raftDir); err != nil {
		if os.IsNotExist(err) {
			// path does not exist
			log.Debugf("raft: creating data dir %s", store.raftDir)
			if err := os.MkdirAll(store.raftDir, os.ModePerm); err != nil {
				return log.Errorf("RaftDataDir (%s) does not exist and cannot be created: %+v", store.raftDir, err)
			}
		} else {
			// Other error
			return log.Errorf("RaftDataDir (%s) error: %+v", store.raftDir, err)
		}
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := NewFileSnapshotStore(store.raftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return log.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	logStore := NewRelationalStore(store.raftDir)
	log.Debugf("raft: logStore=%+v", logStore)

	// Instantiate the Raft systems.
	if store.raft, err = raft.NewRaft(config, (*fsm)(store), logStore, logStore, snapshots, transport); err != nil {
		return fmt.Errorf("error creating new raft: %s", err)
	}

	// Bootstrap if needed: no existing state and either no peers or single-node mode.
	hasState, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		return fmt.Errorf("error checking existing state: %s", err)
	}
	if !hasState {
		var servers []raft.Server
		if len(peerNodes) == 0 {
			// Single-node mode: bootstrap with just this node
			log.Infof("enabling single-node mode")
			servers = []raft.Server{
				{
					ID:      raft.ServerID(store.raftAdvertise),
					Address: raft.ServerAddress(store.raftAdvertise),
				},
			}
		} else {
			// Multi-node: include all peers and this node
			seen := make(map[string]bool)
			for _, peerNode := range peerNodes {
				peerNode = strings.TrimSpace(peerNode)
				if peerNode != "" && !seen[peerNode] {
					seen[peerNode] = true
					servers = append(servers, raft.Server{
						ID:      raft.ServerID(peerNode),
						Address: raft.ServerAddress(peerNode),
					})
				}
			}
			// Ensure this node is included
			if !seen[store.raftAdvertise] {
				servers = append(servers, raft.Server{
					ID:      raft.ServerID(store.raftAdvertise),
					Address: raft.ServerAddress(store.raftAdvertise),
				})
			}
		}

		configuration := raft.Configuration{Servers: servers}
		if err := raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, configuration); err != nil {
			return fmt.Errorf("error bootstrapping cluster: %s", err)
		}
		log.Infof("raft: cluster bootstrapped with %d servers", len(servers))
	}

	log.Infof("new raft created")
	return nil
}

// AddPeer adds a node, located at addr, to this store. The node must be ready to
// respond to Raft communications at that address.
func (store *Store) AddPeer(addr string) error {
	log.Infof("received join request for remote node %s", addr)

	f := store.raft.AddVoter(raft.ServerID(addr), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Infof("node at %s joined successfully", addr)
	return nil
}

// RemovePeer removes a node from this raft setup
func (store *Store) RemovePeer(addr string) error {
	log.Infof("received remove request for remote node %s", addr)

	f := store.raft.RemoveServer(raft.ServerID(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Infof("node at %s removed successfully", addr)
	return nil
}

// genericCommand requests consensus for applying a single command.
// This is an internal orchestrator implementation
func (store *Store) genericCommand(op string, bytes []byte) (response interface{}, err error) {
	if store.raft.State() != raft.Leader {
		return nil, fmt.Errorf("not leader")
	}

	b, err := json.Marshal(&storeCommand{Op: op, Value: bytes})
	if err != nil {
		return nil, err
	}

	f := store.raft.Apply(b, raftTimeout)
	if err = f.Error(); err != nil {
		return nil, err
	}
	r := f.Response()
	if err, ok := r.(error); ok && err != nil {
		// This code checks whether the response itself was an error object. If so, it should
		// indicate failure of the operation.
		return r, err
	}
	return r, nil
}
