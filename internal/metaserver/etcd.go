package metaserver

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rfeverything/rfs/internal/logger"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type EtcdStore struct {
	uniqueID int32
	client   *clientv3.Client
	isleader bool
	closed   bool

	volumeUpdateChan chan *VolumeUpdateEvent
}

func genKey(dir, FileName string) (key string) {
	dir = strings.TrimPrefix(dir, "/")
	return strings.Join([]string{"file", dir, FileName}, "/")
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

	klresp, kaerr := client.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		logger.Global().Fatal(err.Error())
	}
	//clean lease keepalive response queue
	go func() {
		for {
			<-klresp
		}
	}()
	es := &EtcdStore{
		client:   client,
		uniqueID: UniqueID,
	}
	es.volumeUpdateChan = make(chan *VolumeUpdateEvent, 1)
	go es.election()
	go es.watchVolume()
	logger.Global().Info("NewEtcdStore Done", zap.Int32("UniqueID", UniqueID))

	return es
}

func (es *EtcdStore) Close() {
	es.closed = false
	es.client.Close()
}

func (es *EtcdStore) InsertEntry(ctx context.Context, dir string, entry *Entry) error {
	key := genKey(dir, entry.Name)

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
	key := genKey(dir, "")
	resp, err := es.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("EtcdStore.Get: %v", err)
	}
	var entries []*Entry
	for _, kv := range resp.Kvs {
		var entry Entry
		if err := entry.DecodeAttributesAndChunks(kv.Value); err != nil {
			continue
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
		if err := json.Unmarshal(kv.Value, &volume); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %v", err)
		}
		volumes = append(volumes, &volume)
	}
	return volumes, nil
}

func (es *EtcdStore) GetVolumesStatusChan() chan *VolumeUpdateEvent {
	return es.volumeUpdateChan
}

func (es *EtcdStore) watchVolume() {
	ctx := context.TODO()
	rch := es.client.Watch(ctx, "/rfs/volumes/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				var volume vpb.VolumeStatus
				if err := json.Unmarshal(ev.Kv.Value, &volume); err != nil {
					logger.Global().Error("json.Unmarshal", zap.Error(err))
					continue
				}
				logger.Global().Info("watchVolume: PUT", zap.String("key", string(ev.Kv.Key)), zap.Any("volume", &volume))
				es.volumeUpdateChan <- &VolumeUpdateEvent{
					Type:     VolumeEventTypePut,
					Status:   &volume,
					VolumeId: string(ev.Kv.Key[len("/rfs/volumes/") : len(ev.Kv.Key)-1]),
				}
			case mvccpb.DELETE:
				logger.Global().Info("watchVolume: DELETE", zap.String("key", string(ev.Kv.Key)))
				es.volumeUpdateChan <- &VolumeUpdateEvent{
					Type:     VolumeEventTypeDelete,
					VolumeId: string(ev.Kv.Key[len("/rfs/volumes/") : len(ev.Kv.Key)-1]),
				}
			}
		}
	}
}
