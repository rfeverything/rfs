package server

import (
	"context"

	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (vs *VolumeServer) PutChunk(context.Context, *vpb.PutChunkRequest) (*vpb.PutChunkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutChunk not implemented")
}
func (vs *VolumeServer) GetChunk(context.Context, *vpb.GetChunkRequest) (*vpb.GetChunkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChunk not implemented")
}
func (vs *VolumeServer) DeleteChunk(context.Context, *vpb.DeleteChunkRequest) (*vpb.DeleteChunkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChunk not implemented")
}
func (vs *VolumeServer) VolumeStatus(context.Context, *vpb.VolumeStatusRequest) (*vpb.VolumeStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeStatus not implemented")
}
func (vs *VolumeServer) Ping(context.Context, *vpb.PingRequest) (*vpb.PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
