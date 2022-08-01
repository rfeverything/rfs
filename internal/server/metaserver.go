package metaserver

import (
	"github.com/rfeverything/rfs/internal/etcd"
	"go.uber.org/zap"
)

type MetaServer struct {
	EtcdClient *etcd.EtcdClient
	Logger     *zap.Logger

	UniqueId int32
}

func NewMetaServer() *MetaServer {
	return &MetaServer{}
}
