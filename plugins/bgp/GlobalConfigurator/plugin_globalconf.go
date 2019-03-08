package GlobalConfigurator

import (
	"github.com/contiv/bgp-vpp/plugins/bgp/model"
)

// CreateGlobalConf creates a new plugin in etcd
func (b *GlobalConfHandler) CreateGlobalConf(val *model.GlobalConf) error {
	err := b.broker.Put(val.GetName(), val)
	if err != nil {
		b.log.Errorf("Could not create configuration")
		return err
	}

	return nil
}

// DeleteGlobalConf deletes a plugin in etcd
func (b *GlobalConfHandler) DeleteGlobalConf(key string) error {
	existed, err := b.broker.Delete(key)
	if err != nil {
		b.log.Errorf("Could not delete configuration")
	}
	b.log.Infof("Configuration existed: %v", existed)

	return nil
}
