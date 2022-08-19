package volume

import (
	"context"

	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"github.com/tecbot/gorocksdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VolumeServer struct {
	db *gorocksdb.DB
}

func (vs *VolumeServer) PutChunk(ctx context.Context, req *vpb.PutChunkRequest) (*vpb.PutChunkResponse, error) {
	for _, chunk := range req.Chunks {
		key := chunk.Chunkid
		value := chunk.Content
		err := vs.db.Put(gorocksdb.NewDefaultWriteOptions(), []byte{byte(key)}, value)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
func (vs *VolumeServer) GetChunk(ctx context.Context, req *vpb.GetChunkRequest) (*vpb.GetChunkResponse, error) {
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
	key := req.ChunkId
	err := vs.db.Delete(gorocksdb.NewDefaultWriteOptions(), []byte{byte(key)})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (vs *VolumeServer) VolumeStatus(ctx context.Context, req *vpb.VolumeStatusRequest) (*vpb.VolumeStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeStatus not implemented")
}
func (vs *VolumeServer) Ping(ctx context.Context, req *vpb.PingRequest) (*vpb.PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
