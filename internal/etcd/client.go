package etcd

import (
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
}

func NewEtcdClient() *EtcdClient {
	var client *clientv3.Client
	var err error
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return nil
	}

	return &EtcdClient{client: client}
}
