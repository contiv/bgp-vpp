package GlobalConfigurator

import (
	"github.com/contiv/bgp-vpp/plugins/bgp/model"
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logrus"
)

// GlobalConfAPI provides methods for CRUD operations on etcd
type GlobalConfAPI interface {
	GlobalConfWrite
}

// GlobalConfWrite provides write methods for ETCD
type GlobalConfWrite interface {
	// CreateGlobalConf adds new plugin to etcd
	CreateGlobalConf(val *model.GlobalConf) error
	// DeleteGlobalConf deletes plugin from etcd
	DeleteGlobalConf(key string) error
}

// GlobalConfHandler is accessor to etcd related GlobalConf methods
type GlobalConfHandler struct {
	log    logging.Logger
	broker keyval.ProtoBroker
}

// NewGlobalConfHandler creates new instance of GlobalConfHandler
func NewGlobalConfHandler(log logging.Logger, broker keyval.ProtoBroker) *GlobalConfHandler {
	if log == nil {
		log = logrus.NewLogger("global-conf-handler")
	}
	return &GlobalConfHandler{
		log:    log,
		broker: broker,
	}
}
