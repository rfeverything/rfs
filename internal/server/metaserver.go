package server

import (
	"github.com/rfeverything/rfs/internal/logger"
	mst "github.com/rfeverything/rfs/internal/metaserver"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

type MetaServer struct {
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
		return nil, err
	}
	return ms, err
}
