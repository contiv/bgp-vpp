package GlobalConfigurator

import (
	"github.com/contiv/bgp-vpp/plugins/bgp/model"
)

// CreatePlugin creates a new plugin in etcd
func (b *GlobalConfHandler) CreatePlugin(val *model.GlobalConf) error {
	err := b.broker.Put(val.GetName(), val)
	if err != nil {
		b.log.Errorf("Could not create plugin")
		return err
	}

	return nil
}

// DeletePlugin deletes a plugin in etcd
func (b *GlobalConfHandler) DeletePlugin(key string) error {
	existed, err := b.broker.Delete(key)
	if err != nil {
		b.log.Errorf("Could not delete plugin")
	}
	b.log.Infof("Plugin existed: %v", existed)

	return nil
}
