package metaserver

import (
	"context"
	"errors"
	"os"

	"github.com/rfeverything/rfs/internal/constant"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

func (ms *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	pth := req.GetPath()
	e := &Entry{}
	FromProtoEntry(req.GetEntry(), e)
	if err := e.SplitToChunks(constant.ChunkSize); err != nil {
		return nil, err
	}
	volumes, err := ms.Store.GetVolumesStatus()
	if err != nil {
		return nil, err
	}
	if err := ms.assignVolumeForChunks(volumes, e.Chunks); err != nil {
		return nil, err
	}
	//TODO : exist check
	//TODO : recovery from crash
	if err := ms.Store.InsertEntry(ctx, pth, e); err != nil {
		return nil, err
	}
	//TODO: 2pc
	//TODO: Time out
	for _, chunk := range e.Chunks {
		for _, volumeserver := range chunk.VolumeIds {
			resp, err := ms.VolumeClients[volumeserver].PutChunk(ctx, &vpb.PutChunkRequest{
				Chunks: []*rfspb.FileChunk{chunk},
			})
			if err != nil {
				return nil, err
			}
			if resp.Error != "" {
				return nil, errors.New(resp.Error)
			}
		}
	}
	return &mpb.CreateFileResponse{
		Error: "",
	}, nil
}

func (ms *MetaServer) AppendFile(ctx context.Context, req *mpb.AppendFileRequest) (*mpb.AppendFileResponse, error) {
	return nil, nil
}

func (ms *MetaServer) Mkdir(ctx context.Context, req *mpb.MkdirRequest) (*mpb.MkdirResponse, error) {
	pth := req.GetDirectory()
	e := &Entry{
		Path: pth,
	}
	e.Mode = os.ModeDir | 0755
	if err := ms.Store.InsertEntry(ctx, pth, e); err != nil {
		return nil, err
	}
	return nil, nil
}
