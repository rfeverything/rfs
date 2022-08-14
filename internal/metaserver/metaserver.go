package metaserver

import (
	"context"
	"math/rand"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

type MetaServer struct {
	Store    *EtcdStore
	UniqueId int32
}

func NewMetaStore() (*MetaServer, error) {
	UniqueID := rand.Int31()
	Store := NewEtcdStore(UniqueID)
	return &MetaServer{
		Store:    Store,
		UniqueId: UniqueID,
	}, nil
}

func (ms *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	return nil, nil
}
