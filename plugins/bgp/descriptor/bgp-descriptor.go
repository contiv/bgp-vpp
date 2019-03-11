package descriptor

/*
const (
	sswanPoolDescriptorName = "sswan-pool"
)

// NewSSwanPoolDescriptor creates an instance of a SSwanPoolDescriptor

func NewSSwanPoolDescriptor(configurator configurator.SSwanConfigurator) *kvs.KVDescriptor {

	poolDescriptor := &adapter.SSwanPoolDescriptor{

		Name: sswanPoolDescriptorName,

		// TODO: replace these four fields with a Model

		NBKeyPrefix: model.ModelSSwanPool.KeyPrefix(),

		ValueTypeName: model.ModelSSwanPool.ProtoName(),

		KeySelector: model.ModelSSwanPool.IsKeyValid,

		KeyLabel: model.ModelSSwanPool.StripKeyPrefix,

		// Add = C for CRUD (will be renamed to Create very soon)

		Create: func(key string, value *model.SSwanPool) (metadata interface{}, err error) {

			err = configurator.AddRecord(value)

			return

		},

		Delete: func(key string, value *model.SSwanPool, metadata interface{}) (err error) {

			err = configurator.DeleteRecord(value)

			return

		},

		// Modify (will be renamed to Update) is not needed, change is always performed with a full re-creation

		// Dump (Read) is not yet supported, leave undefined

		UpdateWithRecreate: func(key string, oldValue, newValue *model.SSwanPool, metadata interface{}) bool {

			// Modify always performed via re-creation

			return true

		},

		// No Dependencies

	}

	return adapter.NewSSwanPoolDescriptor(poolDescriptor)

}
*/