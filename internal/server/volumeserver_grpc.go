package server

import (
	"context"

	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

func (vs *VolumeServer) PutChunk(ctx context.Context, req *vpb.PutChunkRequest) (*vpb.PutChunkResponse, error) {
	return vs.vs.PutChunk(ctx, req)
}
func (vs *VolumeServer) GetChunk(ctx context.Context, req *vpb.GetChunkRequest) (*vpb.GetChunkResponse, error) {
	return vs.vs.GetChunk(ctx, req)
}
func (vs *VolumeServer) DeleteChunk(ctx context.Context, req *vpb.DeleteChunkRequest) (*vpb.DeleteChunkResponse, error) {
	return vs.vs.DeleteChunk(ctx, req)
}
func (vs *VolumeServer) VolumeStatus(ctx context.Context, req *vpb.VolumeStatusRequest) (*vpb.VolumeStatusResponse, error) {
	return vs.vs.VolumeStatus(ctx, req)
}
func (vs *VolumeServer) Ping(ctx context.Context, req *vpb.PingRequest) (*vpb.PingResponse, error) {
	return vs.vs.Ping(ctx, req)
}
