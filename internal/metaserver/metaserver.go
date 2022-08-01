package metaserver

import "github.com/rfeverything/rfs/internal/etcd"

type MetaServer struct {
	EtcdClient *etcd.EtcdClient

	UniqueId int32
}

func NewMetaServer() *MetaServer {
	return &MetaServer{}
}
