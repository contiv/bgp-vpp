package model

import "github.com/ligato/vpp-agent/pkg/models"

const ModuleName = "bgp"

var (
	ModelBgpGlobal = models.Register(&GlobalConf{}, models.Spec{
		Module:  ModuleName,
		Version: "v1",
		Type:    "global",
	}, models.WithNameTemplate("{{.Name}}"))
	ModelBgpPeer = models.Register(&PeerConf{}, models.Spec{
		Module:  ModuleName,
		Version: "v1",
		Type:    "peers",
	})
)

//given the peer name, this function will return the key
func PeerKey(name string) string {
	return models.Key(&PeerConf{
		Name: name,
	})
}

func GlobalKey(name string) string {
	return models.Key(&GlobalConf{
		Name: name,
	})
}
