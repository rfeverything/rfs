package server

import (
	mst "github.com/rfeverything/rfs/internal/metastore"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	"go.uber.org/zap"
)

type MetaServer struct {
	Logger    *zap.Logger
	MetaStore *mst.MetaStore

	*mpb.UnimplementedMetaServerServer

	UniqueId int32
}

func NewMetaServer() *MetaServer {
	return &MetaServer{}
}
