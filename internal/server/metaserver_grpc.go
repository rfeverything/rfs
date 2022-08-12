package server

import (
	"context"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

func (m *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	return nil, nil
}

func (m *MetaServer) AppendFile(ctx context.Context, req *mpb.AppendFileRequest) (*mpb.AppendFileResponse, error) {
	return nil, nil
}

func (m *MetaServer) GetFileStat(ctx context.Context, req *mpb.GetFileStatRequest) (*mpb.GetFileStatResponse, error) {
	return nil, nil
}

func (m *MetaServer) DeleteFile(ctx context.Context, req *mpb.DeleteFileRequest) (*mpb.DeleteFileResponse, error) {
	return nil, nil
}

func (m *MetaServer) RenameFile(ctx context.Context, req *mpb.RenameFileRequest) (*mpb.RenameFileResponse, error) {
	return nil, nil
}

func (m *MetaServer) Mkdir(ctx context.Context, req *mpb.MkdirRequest) (*mpb.MkdirResponse, error) {
	return nil, nil
}

func (m *MetaServer) DeleteDir(ctx context.Context, req *mpb.DeleteDirRequest) (*mpb.DeleteDirResponse, error) {
	return nil, nil
}

func (m *MetaServer) List(ctx context.Context, req *mpb.ListRequest) (*mpb.ListResponse, error) {
	return nil, nil
}

func (m *MetaServer) RenameDir(ctx context.Context, req *mpb.RenameDirRequest) (*mpb.RenameDirResponse, error) {
	return nil, nil
}

func (m *MetaServer) MoveDir(ctx context.Context, req *mpb.MoveDirRequest) (*mpb.MoveDirResponse, error) {
	return nil, nil
}
