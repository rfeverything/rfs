package metastore

import (
	"context"
	"math/rand"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

type MetaStore struct {
	Store    *EtcdStore
	UniqueId int32
}

func NewMetaStore() (*MetaStore, error) {
	UniqueID := rand.Int31()
	Store := NewEtcdStore(UniqueID)
	return &MetaStore{
		Store:    Store,
		UniqueId: UniqueID,
	}, nil
}

func (ms *MetaStore) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	return nil, nil
}
