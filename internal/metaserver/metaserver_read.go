package metaserver

import (
	"context"
	"errors"
	"sync"

	"github.com/rfeverything/rfs/internal/logger"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

func (ms *MetaServer) GetFile(ctx context.Context, req *mpb.GetFileRequest) (*mpb.GetFileResponse, error) {
	pth := req.GetPath()
	e, err := ms.Store.GetEntry(ctx, pth)
	if err != nil {
		return nil, err
	}
	wait := sync.WaitGroup{}
	gerr := make(chan error, len(e.Chunks))
	for _, chunk := range e.Chunks {
		wait.Add(1)
		req := &vpb.GetChunkRequest{
			ChunkId: chunk.Chunkid,
		}
		go func(chunk *rfspb.FileChunk, req *vpb.GetChunkRequest) {
			for _, volumeserver := range chunk.VolumeIds {
				if _, exist := ms.VolumeClients[volumeserver]; !exist {
					logger.Global().Sugar().Errorf("volumeserver %s not exist", volumeserver)
					continue
				}
				resp, err := ms.VolumeClients[volumeserver].GetChunk(ctx, req)
				if err != nil {
					logger.Global().Sugar().Errorf("get chunk error: %v", err)
					continue
				}
				if resp.Error != "" {
					logger.Global().Sugar().Errorf("get chunk error: %v", resp.Error)
					continue
				}
				chunk.Content = resp.GetChunk().Content
				wait.Done()
				return
			}
			gerr <- errors.New("no volume server available")
			wait.Done()
		}(chunk, req)
	}
	wait.Wait()
	close(gerr)
	for err := range gerr {
		if err != nil {
			return nil, err
		}
	}
	if err := e.CombineChunksGetContent(); err != nil {
		return nil, err
	}

	re := &rfspb.Entry{}
	if err := e.ToExistingProtoEntry(re, true); err != nil {
		return nil, err
	}

	return &mpb.GetFileResponse{
		Error: "",
		Entry: re,
	}, nil
}

func (ms *MetaServer) Stat(ctx context.Context, req *mpb.StatRequest) (*mpb.StatResponse, error) {
	logger.Global().Sugar().Debugf("stat %s", req.GetPath())
	pth := req.GetPath()
	e, err := ms.Store.GetEntry(ctx, pth)
	if err != nil {
		logger.Global().Sugar().Errorf("get entry error: %v", err)
		return nil, err
	}
	re := &rfspb.Entry{}

	if err := e.ToExistingProtoEntry(re, false); err != nil {
		logger.Global().Sugar().Errorf("to proto entry error: %v", err)
		return nil, err
	}

	return &mpb.StatResponse{
		Entry: re,
	}, nil
}

func (ms *MetaServer) List(ctx context.Context, req *mpb.ListRequest) (*mpb.ListResponse, error) {
	logger.Global().Sugar().Debugf("list %s", req.GetDir())
	pth := req.GetDir()
	// e, err := ms.Store.GetEntry(ctx, pth)
	// if err != nil {
	// 	return nil, err
	// }
	// if !e.IsDir() {
	// 	return nil, errors.New("not a directory")
	// }
	entries, err := ms.Store.ListEntries(ctx, pth)
	if err != nil {
		logger.Global().Sugar().Errorf("list entries error: %v", err)
		return nil, err
	}
	re := make([]*rfspb.Entry, 0)
	for _, entry := range entries {
		res := &rfspb.Entry{}
		if err := entry.ToExistingProtoEntry(res, false); err != nil {
			logger.Global().Sugar().Errorf("to proto entry error: %v", err)
			return nil, err
		}
		re = append(re, res)
	}
	return &mpb.ListResponse{
		Entries: re,
	}, nil
}
