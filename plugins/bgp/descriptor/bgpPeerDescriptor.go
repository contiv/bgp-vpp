package descriptor

import (
	"context"

	"github.com/contiv/bgp-vpp/plugins/bgp/descriptor/adapter"
	"github.com/contiv/bgp-vpp/plugins/bgp/model"
	"github.com/ligato/cn-infra/logging"
	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	bgp_api "github.com/osrg/gobgp/api"
	gobgp "github.com/osrg/gobgp/pkg/server"
)

const (
	peerDescriptorName = "peer-conf"
)

//our descriptor
type PeerDescriptor struct {
	log logging.Logger
	//scheduler kvs.KVScheduler
	server *gobgp.BgpServer
}

// NewPeerConfDescriptor creates a new instance of the descriptor.
func NewPeerConfDescriptor(log logging.PluginLogger) *PeerDescriptor {
	// Set plugin descriptor init values
	return &PeerDescriptor{
		log: log.NewLogger("peer-conf-descriptor"),
	}
}

// GetDescriptor returns descriptor suitable for registration (via adapter) with the KVScheduler.
func (d *PeerDescriptor) GetDescriptor() *adapter.PeerConfDescriptor {
	return &adapter.PeerConfDescriptor{
		Name:          peerDescriptorName,
		NBKeyPrefix:   model.ModelBgpPeer.KeyPrefix(),
		ValueTypeName: model.ModelBgpPeer.ProtoName(),
		KeySelector:   model.ModelBgpPeer.IsKeyValid,
		KeyLabel:      model.ModelBgpPeer.StripKeyPrefix,
		Create:        d.Create,
		Delete:        d.Delete,
		//UpdateWithRecreate: d.UpdateWithRecreate,
		//Retrieve:             d.Retrieve,
		Dependencies:         d.Dependencies,
		RetrieveDependencies: []string{},
	}
}

// Create creates new value.
func (d *PeerDescriptor) Create(key string, value *model.PeerConf) (metadata interface{}, err error) {
	n := &bgp_api.Peer{
		Conf: &bgp_api.PeerConf{
			NeighborAddress: value.NeighborAddress,
			PeerAs:          value.PeerAs,
		},
	}
	err = d.server.AddPeer(context.Background(), &bgp_api.AddPeerRequest{
		Peer: n,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete removes an existing value.
func (d *PeerDescriptor) Delete(key string, value *model.PeerConf, metadata interface{}) error {
	err := d.server.DeletePeer(context.Background(), &bgp_api.DeletePeerRequest{})
	if err != nil {
		return err
	}
	return nil
}

// UpdateWithRecreate returns true if value update requires full re-creation.
/*func (d *GlobalDescriptor) UpdateWithRecreate(key string, old, new *model.GlobalConf, metadata interface{}) bool {
	return true
}*/

// Retrieve retrieves values from SB.
/*func (d *GlobalDescriptor) Retrieve(correlate []adapter.PluginKVWithMetadata) (retrieved []adapter.PluginKVWithMetadata, err error) {
	return retrieved, nil
}*/

// Dependencies lists dependencies of the given value.
func (d *PeerDescriptor) Dependencies(key string, value *model.PeerConf) (deps []kvs.Dependency) {
	return deps
}
