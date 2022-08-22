package volume

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	gorocksdb "github.com/linxGnu/grocksdb"
	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/internal/logger"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	clientv3 "go.etcd.io/etcd/client/v3"

	"google.golang.org/protobuf/proto"
)

type VolumeServer struct {
	ID              string
	TotalChunkCount uint64
	Size            uint64
	Used            uint64
	ChunksSet       map[int64]uint64 // chunkid -> size

	Host string
	Port int

	sdb     *gorocksdb.DB
	db      *gorocksdb.DB
	etcd    *clientv3.Client
	leaseID clientv3.LeaseID
}

func NewVolumeServer() *VolumeServer {
	vs := &VolumeServer{
		Host: config.Global().GetString("volume.host"),
		Port: config.Global().GetInt("volume.port"),
	}

	id := config.Global().GetString("volume.uid")
	if id == "" {
		id = uuid.New().String()
		config.Global().Set("volume.uid", id)
		config.Global().WriteConfig()
	}
	vs.ID = id

	size := config.Global().GetString("volume.size")
	vs.Size, _ = strconv.ParseUint(size, 10, 64)

	opt := gorocksdb.NewDefaultOptions()
	opt.SetCreateIfMissing(true)
	sdb, err := gorocksdb.OpenDb(opt, "./volume_state")
	if err != nil {
		logger.Global().Sugar().Fatalf("open volume state db failed: %v", err)
	}
	vs.sdb = sdb
	db, err := gorocksdb.OpenDb(opt, "./volume_db")
	if err != nil {
		logger.Global().Sugar().Fatalf("open volume db failed: %v", err)
	}
	vs.db = db
	vs.recoverFromPersistence()
	logger.Global().Sugar().Infof("volume server %s started", vs.ID)
	vs.registerToEtcd()
	vs.persist()
	return vs
}

func (vs *VolumeServer) registerToEtcd() {
	cfg := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		logger.Global().Sugar().Fatalf("connect to etcd failed: %v", err)
	}
	resp, err := c.Grant(context.TODO(), 5)
	if err != nil {
		logger.Global().Sugar().Fatalf("grant failed: %v", err)
	}
	klresp, kaerr := c.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		logger.Global().Sugar().Fatalf("keep alive failed: %v", kaerr)
	}
	//clean lease keepalive response queue
	go func() {
		for {
			<-klresp
		}
	}()
	vs.etcd = c
	vs.leaseID = resp.ID
}

func (vs *VolumeServer) persist() {
	s := vs.getStatus()
	logger.Global().Sugar().Infof("persist volume status: %v", s)
	data, err := proto.Marshal(s)
	if err != nil {
		logger.Global().Sugar().Errorf("marshal volume status failed: %v", err)
	}
	err = vs.sdb.Put(gorocksdb.NewDefaultWriteOptions(), []byte(vs.ID), data)
	if err != nil {
		logger.Global().Sugar().Errorf("put volume status failed: %v", err)
	}
	jdata, err := json.Marshal(s)
	if err != nil {
		logger.Global().Sugar().Errorf("marshal volume status failed: %v", err)
	}
	_, err = vs.etcd.Put(context.Background(), "/rfs/volumes/"+vs.ID, string(jdata), clientv3.WithLease(vs.leaseID))
	if err != nil {
		logger.Global().Sugar().Errorf("put volume status to etcd failed: %v", err)
	}
	logger.Global().Sugar().Infof("volume server %s persisted", vs.ID)
}

func (vs *VolumeServer) recoverFromPersistence() {
	data, err := vs.sdb.Get(gorocksdb.NewDefaultReadOptions(), []byte(vs.ID))
	if err != nil {
		logger.Global().Sugar().Errorf("get volume status from db failed: %v", err)
	}
	s := &vpb.VolumeStatus{}
	err = proto.Unmarshal(data.Data(), s)
	if err != nil {
		logger.Global().Sugar().Errorf("unmarshal volume status failed: %v", err)
	}
	vs.TotalChunkCount = s.ChunkCount
	vs.Used = s.Used
	vs.ChunksSet = s.ChunkIds
}

func (vs *VolumeServer) PutChunk(ctx context.Context, req *vpb.PutChunkRequest) (*vpb.PutChunkResponse, error) {
	for _, chunk := range req.Chunks {
		logger.Global().Sugar().Infof("put chunk %d", chunk.Chunkid)
		key := chunk.Chunkid
		value := chunk.Content
		err := vs.db.Put(gorocksdb.NewDefaultWriteOptions(), []byte{byte(key)}, value)
		if err != nil {
			return nil, err
		}
		vs.ChunksSet[key] = uint64(len(value))
		vs.TotalChunkCount++
		vs.Used += uint64(len(value))
		vs.persist()
	}
	return nil, nil
}
func (vs *VolumeServer) GetChunk(ctx context.Context, req *vpb.GetChunkRequest) (*vpb.GetChunkResponse, error) {
	logger.Global().Sugar().Infof("get chunk %d", req.ChunkId)
	key := req.ChunkId
	value, err := vs.db.Get(gorocksdb.NewDefaultReadOptions(), []byte{byte(key)})
	if err != nil {
		return nil, err
	}

	return &vpb.GetChunkResponse{
		Chunk: &rfspb.FileChunk{
			Chunkid: key,
			Content: value.Data(),
		},
	}, nil

}
func (vs *VolumeServer) DeleteChunk(ctx context.Context, req *vpb.DeleteChunkRequest) (*vpb.DeleteChunkResponse, error) {
	logger.Global().Sugar().Infof("delete chunk %d", req.ChunkId)
	key := req.ChunkId
	err := vs.db.Delete(gorocksdb.NewDefaultWriteOptions(), []byte{byte(key)})
	if err != nil {
		return nil, err
	}
	delete(vs.ChunksSet, key)
	vs.TotalChunkCount--
	vs.Used -= uint64(vs.ChunksSet[key])
	vs.persist()
	return nil, nil
}

func (vs *VolumeServer) VolumeStatus(ctx context.Context, req *vpb.VolumeStatusRequest) (*vpb.VolumeStatusResponse, error) {
	s := vs.getStatus()
	s.VolumeId = vs.ID
	return &vpb.VolumeStatusResponse{
		VolumeStatus: s,
	}, nil
}

func (vs *VolumeServer) getStatus() *vpb.VolumeStatus {
	return &vpb.VolumeStatus{
		ChunkCount: vs.TotalChunkCount,
		Size:       vs.Size,
		Used:       vs.Used,
		Free:       vs.Size - vs.Used,
		Address:    vs.Host + ":" + fmt.Sprint(vs.Port),
		ChunkIds:   vs.ChunksSet,
	}
}

func (vs *VolumeServer) Ping(ctx context.Context, req *vpb.PingRequest) (*vpb.PingResponse, error) {
	return &vpb.PingResponse{}, nil
}
