package server

import (
	"context"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

func (m *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	return m.metaserver.CreateFile(ctx, req)
}

func (m *MetaServer) AppendFile(ctx context.Context, req *mpb.AppendFileRequest) (*mpb.AppendFileResponse, error) {
	return m.metaserver.AppendFile(ctx, req)
}

func (m *MetaServer) GetFile(ctx context.Context, req *mpb.GetFileRequest) (*mpb.GetFileResponse, error) {
	return m.metaserver.GetFile(ctx, req)
}

func (m *MetaServer) Stat(ctx context.Context, req *mpb.StatRequest) (*mpb.StatResponse, error) {
	return m.metaserver.Stat(ctx, req)
}

func (m *MetaServer) Remove(ctx context.Context, req *mpb.RemoveRequest) (*mpb.RemoveResponse, error) {
	return m.metaserver.Remove(ctx, req)
}

func (m *MetaServer) Mkdir(ctx context.Context, req *mpb.MkdirRequest) (*mpb.MkdirResponse, error) {
	return m.metaserver.Mkdir(ctx, req)
}

func (m *MetaServer) List(ctx context.Context, req *mpb.ListRequest) (*mpb.ListResponse, error) {
	return m.metaserver.List(ctx, req)
}

func (m *MetaServer) Move(ctx context.Context, req *mpb.MoveRequest) (*mpb.MoveResponse, error) {
	return m.metaserver.Move(ctx, req)
}
