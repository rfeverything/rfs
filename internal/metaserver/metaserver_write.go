package metaserver

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/rfeverything/rfs/internal/constant"
	"github.com/rfeverything/rfs/internal/logger"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"go.uber.org/zap"
)

func (ms *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	logger.Global().Info("CreateFile", zap.String("path", req.GetDir()))
	dir := req.GetDir()
	e := &Entry{}
	FromProtoEntry(req.GetEntry(), e)
	if err := e.SplitToChunks(constant.ChunkSize); err != nil {
		logger.Global().Error("CreateFile", zap.Error(err))
		return nil, err
	}
	volumes, err := ms.Store.GetVolumesStatus()
	if err != nil {
		logger.Global().Error("CreateFile", zap.Error(err))
		return nil, err
	}
	if err := ms.assignVolumeForChunks(volumes, e.Chunks); err != nil {
		logger.Global().Error("CreateFile", zap.Error(err))
		return nil, err
	}
	//TODO : exist check
	//TODO : recovery from crash
	if err := ms.Store.InsertEntry(ctx, dir, e); err != nil {
		logger.Global().Error("CreateFile", zap.Error(err))
		return nil, err
	}
	//TODO: 2pc
	//TODO: Time out
	wait := sync.WaitGroup{}
	gerr := make(chan error, len(e.Chunks))
	for _, chunk := range e.Chunks {
		for _, volumeserver := range chunk.VolumeIds {
			logger.Global().Info("CreateFile", zap.String("volumeserver", volumeserver))
			if _, e := ms.VolumeClients[volumeserver]; !e {
				logger.Global().Error("CreateFile", zap.String("volumeserver", volumeserver))
				return nil, errors.New("volume server not found")
			}
			req := &vpb.PutChunkRequest{
				Chunks: []*rfspb.FileChunk{chunk},
			}
			wait.Add(1)
			go func(volumeserver string, req *vpb.PutChunkRequest) {
				defer wait.Done()
				resp, err := ms.VolumeClients[volumeserver].PutChunk(ctx, req)
				if err != nil {
					logger.Global().Error("CreateFile", zap.Error(err))
					gerr <- err
					return
				}
				if resp.Error != "" {
					logger.Global().Error("CreateFile", zap.Error(errors.New(resp.Error)))
					gerr <- errors.New(resp.Error)
					return
				}
			}(volumeserver, req)
		}
	}
	wait.Wait()
	if len(gerr) > 0 {
		return nil, <-gerr
	}
	return &mpb.CreateFileResponse{
		Error: "",
	}, nil
}

// TODO
func (ms *MetaServer) AppendFile(ctx context.Context, req *mpb.AppendFileRequest) (*mpb.AppendFileResponse, error) {
	return nil, nil
}

func (ms *MetaServer) Mkdir(ctx context.Context, req *mpb.MkdirRequest) (*mpb.MkdirResponse, error) {
	pth := req.GetDirectory()
	e := &Entry{
		Name: pth,
	}
	e.Mode = os.ModeDir | 0755
	if err := ms.Store.InsertEntry(ctx, pth, e); err != nil {
		return nil, err
	}
	return nil, nil
}
