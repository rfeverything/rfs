package server

import (
	mst "github.com/rfeverything/rfs/internal/metaserver"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	"go.uber.org/zap"
)

type MetaServer struct {
	Logger    *zap.Logger
	MetaStore *mst.MetaServer

	*mpb.UnimplementedMetaServerServer

	UniqueId int32
}

func NewMetaServer() *MetaServer {
	return &MetaServer{}
}
