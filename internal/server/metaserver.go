package server

import (
	"github.com/rfeverything/rfs/internal/logger"
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

func NewMetaServer() (*MetaServer, error) {
	ms := &MetaServer{}
	var err error
	logger.Global().Info("NewMetaServer")
	ms.metaserver, err = mst.NewMetaServer()
	if err != nil {
		ms.logger.Fatal("NewMetaServer", zap.Error(err))
	}
	return ms, err
}
