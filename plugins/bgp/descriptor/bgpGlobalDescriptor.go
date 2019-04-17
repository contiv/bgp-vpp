package descriptor

import (
	"context"
	"net"

	"github.com/contiv/bgp-vpp/plugins/bgp/descriptor/adapter"
	"github.com/contiv/bgp-vpp/plugins/bgp/model"
	"github.com/ligato/cn-infra/logging"
	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	bgpapi "github.com/osrg/gobgp/api"
	gobgp "github.com/osrg/gobgp/pkg/server"
)

const (
	globalDescriptorName = "global-conf"
)

//our descriptor
type GlobalDescriptor struct {
	log    logging.Logger
	server *gobgp.BgpServer
}

// NewGlobalConfDescriptor creates a new instance of the descriptor.
func NewGlobalConfDescriptor(log logging.PluginLogger, server *gobgp.BgpServer) *kvs.KVDescriptor {
	d := &GlobalDescriptor{log: log, server: server}

	// Set plugin descriptor init values
	gcd := &adapter.GlobalConfDescriptor{
		Name:          globalDescriptorName,
		NBKeyPrefix:   model.ModelBgpGlobal.KeyPrefix(),
		ValueTypeName: model.ModelBgpGlobal.ProtoName(),
		KeySelector:   model.ModelBgpGlobal.IsKeyValid,
		KeyLabel:      model.ModelBgpGlobal.StripKeyPrefix,
		Create:        d.Create,
		Delete:        d.Delete,
		UpdateWithRecreate: func(key string, oldValue, newValue *model.GlobalConf, metadata interface{}) bool {
			// Modify always performed via re-creation
			return true
		},
	}
	return adapter.NewGlobalConfDescriptor(gcd)
}

// Create creates new value.
func (d *GlobalDescriptor) Create(key string, value *model.GlobalConf) (metadata interface{}, err error) {
	//Checks if RouterId is a valid IP Address
	if net.ParseIP(value.RouterId) == nil {
		d.log.Errorf("Invalid IP Address for RouterId = %s", value.RouterId)
		return nil, err
	}
	//Checks if AS is above 64512 and below 65536
	if value.As < 64512 || value.As > 65536 {
		d.log.Errorf("Invalid AS Number = %d. AS Number should be above 64512 and below 65536", value.As)
		return nil, err
	}
	//Checks if ListenPort is -1
	if value.ListenPort != -1 {
		d.log.Errorf("Invalid ListenPort = %d. ListenPort should be -1", value.ListenPort)
		return nil, err
	}
	d.log.Infof("Creating GlobalConf As = %d, RouterId = %s, ListenPort = %d",
		value.As, value.RouterId, value.ListenPort)
	err = d.server.StartBgp(context.Background(), &bgpapi.StartBgpRequest{
		Global: &bgpapi.Global{
			As:         value.As,
			RouterId:   value.RouterId,
			ListenPort: value.ListenPort,
		},
	})

	if err != nil {
		d.log.Errorf("Error creating GlobalConf = %s", err)
		return nil, err
	}

	return nil, nil
}

// Delete removes an existing value.
func (d *GlobalDescriptor) Delete(key string, value *model.GlobalConf, metadata interface{}) error {
	d.log.Infof("Deleting GlobalConf As = %d, RouterId = %s, ListenPort = %d",
		value.As, value.RouterId, value.ListenPort)
	err := d.server.StopBgp(context.Background(), &bgpapi.StopBgpRequest{})
	if err != nil {
		d.log.Errorf("Error deleting GlobalConf = %s", err)
		return err
	}
	return nil
}
