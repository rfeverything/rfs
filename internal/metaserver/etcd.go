package metaserver

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/rfeverything/rfs/internal/logger"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
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
	key = append(key, []byte("/")...)
	key = append(key, []byte(FileName)...)
	return key
}

func NewEtcdStore(UniqueID int32) *EtcdStore {
	logger.Global().Info("NewEtcdStore", zap.Int32("uid", UniqueID))
	var client *clientv3.Client
	var err error
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		logger.Global().Fatal(err.Error())
	}

	resp, err := client.Grant(context.TODO(), 5)
	if err != nil {
		logger.Global().Fatal(err.Error())
	}

	_, err = client.Put(context.TODO(), strings.Join([]string{"/rfs/meta/", string(UniqueID)}, ""), "", clientv3.WithLease(resp.ID))
	if err != nil {
		logger.Global().Fatal(err.Error())
	}

	_, kaerr := client.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		logger.Global().Fatal(err.Error())
	}
	es := &EtcdStore{
		client:   client,
		uniqueID: UniqueID,
	}
	go es.election()
	logger.Global().Info("NewEtcdStore Done", zap.Int32("UniqueID", UniqueID))

	return es
}

func (es *EtcdStore) Close() {
	es.closed = false
	es.client.Close()
}

func (es *EtcdStore) InsertEntry(ctx context.Context, path string, entry *Entry) error {
	key := genKey(filepath.Split(path))

	meta, err := entry.EncodeAttributesAndChunks()
	if err != nil {
		return fmt.Errorf("EncodeAttributesAndChunks: %v", err)
	}

	if _, err := es.client.Put(ctx, string(key), string(meta)); err != nil {
		return fmt.Errorf("EtcdStore.Put: %v", err)
	}
	return nil
}

func (es *EtcdStore) GetEntry(ctx context.Context, path string) (*Entry, error) {
	key := genKey(filepath.Split(path))
	resp, err := es.client.Get(ctx, string(key))
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var entry Entry
	if err := entry.DecodeAttributesAndChunks(resp.Kvs[0].Value); err != nil {
		return nil, fmt.Errorf("DecodeAttributesAndChunks: %v", err)
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
	key := dir
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
	logger.Global().Info("election")
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

		logger.Global().Info("election: success", zap.Int32("uid", es.uniqueID))
		es.isleader = true

		select {
		case <-s.Done():
			es.isleader = false
			logger.Global().Debug("election: expired", zap.Int32("uid", es.uniqueID))
		default:
			logger.Global().Debug("election: running", zap.Int32("uid", es.uniqueID))
			time.Sleep(time.Second)
		}
	}
}

func (es *EtcdStore) GetVolumesStatus() ([]*vpb.VolumeStatus, error) {
	resp, err := es.client.Get(context.TODO(), "/rfs/volumes/", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	var volumes []*vpb.VolumeStatus
	for _, kv := range resp.Kvs {
		var volume vpb.VolumeStatus
		if err := proto.Unmarshal(kv.Value, &volume); err != nil {
			return nil, fmt.Errorf("proto.Unmarshal: %v", err)
		}
		volumes = append(volumes, &volume)
	}
	return volumes, nil
}
