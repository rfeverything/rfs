package metaserver

import (
	"math/rand"

	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

type MetaServer struct {
	Store         *EtcdStore
	UniqueId      int32
	VolumeClients map[string]vpb.VolumeServerClient
}

func NewMetaServer() (*MetaServer, error) {
	UniqueID := rand.Int31()
	Store := NewEtcdStore(UniqueID)
	return &MetaServer{
		Store:    Store,
		UniqueId: UniqueID,
	}, nil
}
