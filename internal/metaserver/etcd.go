package metaserver

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rfeverything/rfs/internal/log"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

type EtcdStore struct {
	uniqueID int32
	client   *clientv3.Client
	isleader bool
	closed   bool
}

func genKey(dir, FileName string) (key []byte) {
	key = []byte(dir)
	key = append(key, []byte(FileName)...)
	return key
}

func NewEtcdStore(UniqueID int32) *EtcdStore {
	var client *clientv3.Client
	var err error
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		log.Global().Fatal(err.Error())
	}

	resp, err := client.Grant(context.TODO(), 5)
	if err != nil {
		log.Global().Fatal(err.Error())
	}

	_, err = client.Put(context.TODO(), strings.Join([]string{"metaserver-", string(UniqueID)}, ""), "", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Global().Fatal(err.Error())
	}

	// to renew the lease only once
	_, kaerr := client.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Global().Fatal(err.Error())
	}

	return &EtcdStore{
		client:   client,
		uniqueID: UniqueID,
	}
}

func (es *EtcdStore) Close() {
	es.closed = false
	es.client.Close()
}

func (es *EtcdStore) InsertEntry(ctx context.Context, path string, entry *Entry) error {
	key := genKey(filepath.Split(path))

	meta, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}

	if _, err := es.client.Put(ctx, string(key), string(meta)); err != nil {
		return fmt.Errorf("EtcdStore.Put: %v", err)
	}
	return nil
}

func (es *EtcdStore) FindEntry(ctx context.Context, path string) (*Entry, error) {
	key := genKey(filepath.Split(path))
	resp, err := es.client.Get(ctx, string(key))
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var entry Entry
	if err := json.Unmarshal(resp.Kvs[0].Value, &entry); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return &entry, nil
}

func (es *EtcdStore) DeleteEntry(ctx context.Context, path string) error {
	key := genKey(filepath.Split(path))
	if _, err := es.client.Delete(ctx, string(key)); err != nil {
		return fmt.Errorf("EtcdStore.Delete: %v", err)
	}
	return nil
}

func (es *EtcdStore) ListEntries(ctx context.Context, dir string) ([]*Entry, error) {
	key := genKey(filepath.Split(dir))
	resp, err := es.client.Get(ctx, string(key), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	var entries []*Entry
	for _, kv := range resp.Kvs {
		var entry Entry
		if err := json.Unmarshal(kv.Value, &entry); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %v", err)
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (es *EtcdStore) KvGet(ctx context.Context, key string) (string, error) {
	resp, err := es.client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("EtcdStore.Get: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (es *EtcdStore) KvPut(ctx context.Context, key, value string) error {
	if _, err := es.client.Put(ctx, key, value); err != nil {
		return fmt.Errorf("EtcdStore.Put: %v", err)
	}
	return nil
}

func (es *EtcdStore) IsLeader() bool {
	return es.isleader
}

func (es *EtcdStore) election() {
	for !es.closed {
		s, err := concurrency.NewSession(es.client, concurrency.WithTTL(5))
		if err != nil {
			fmt.Println(err)
			continue
		}
		e := concurrency.NewElection(s, "metaserver-leader")
		ctx := context.TODO()

		if err = e.Campaign(ctx, string(es.uniqueID)); err != nil {
			fmt.Println(err)
			continue
		}

		log.Global().Info("election: success", zap.Int32("uid", es.uniqueID))
		es.isleader = true

		select {
		case <-s.Done():
			es.isleader = false
			log.Global().Debug("election: expired", zap.Int32("uid", es.uniqueID))
		}
	}
}