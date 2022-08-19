package server

import (
	mst "github.com/rfeverything/rfs/internal/metaserver"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	"go.uber.org/zap"
)

type MetaServer struct {
	logger     *zap.Logger
	metaserver *mst.MetaServer

	*mpb.UnimplementedMetaServerServer

	UniqueId int32
}

func NewMetaServer() *MetaServer {
	ms := &MetaServer{}
	var err error
	ms.metaserver, err = mst.NewMetaServer()
	if err != nil {
		ms.logger.Fatal("NewMetaServer", zap.Error(err))
	}
	return ms
}
