package metastore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdStore struct {
	client *clientv3.Client
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
		fmt.Println(err)
		return nil
	}

	resp, err := client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Put(context.TODO(), strings.Join([]string{"metaserver-", string(UniqueID)}, ""), "", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// to renew the lease only once
	_, kaerr := client.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}

	return &EtcdStore{client: client}
}

func (es *EtcdStore) Close() {
	es.client.Close()
}

func (es *EtcdStore) InsertEntry(ctx context.Context, path string, entry *MetaEntry) error {
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

func (es *EtcdStore) FindEntry(ctx context.Context, path string) (*MetaEntry, error) {
	key := genKey(filepath.Split(path))
	resp, err := es.client.Get(ctx, string(key))
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var entry MetaEntry
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

func (es *EtcdStore) ListEntries(ctx context.Context, dir string) ([]*MetaEntry, error) {
	key := genKey(filepath.Split(dir))
	resp, err := es.client.Get(ctx, string(key), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	var entries []*MetaEntry
	for _, kv := range resp.Kvs {
		var entry MetaEntry
		if err := json.Unmarshal(kv.Value, &entry); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %v", err)
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}
