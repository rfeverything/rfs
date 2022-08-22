package metaserver

import (
	"context"
	"errors"

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
	for _, chunk := range e.Chunks {
		for _, volumeserver := range chunk.VolumeIds {
			resp, err := ms.VolumeClients[volumeserver].GetChunk(ctx, &vpb.GetChunkRequest{
				ChunkId: chunk.Chunkid,
			})
			if err != nil {
				return nil, err
			}
			if resp.Error != "" {
				return nil, errors.New(resp.Error)
			}
			chunk.Content = resp.GetChunk().Content
			break
		}
	}
	if err := e.CombineChunksGetContent(); err != nil {
		return nil, err
	}

	re := &rfspb.Entry{}
	if err := e.ToExistingProtoEntry(re); err != nil {
		return nil, err
	}

	return &mpb.GetFileResponse{
		Error: "",
		Entry: re,
	}, nil
}

func (ms *MetaServer) Stat(ctx context.Context, req *mpb.StatRequest) (*mpb.StatResponse, error) {
	pth := req.GetPath()
	e, err := ms.Store.GetEntry(ctx, pth)
	if err != nil {
		return nil, err
	}
	re := &rfspb.Entry{}

	if err := e.ToExistingProtoEntry(re); err != nil {
		return nil, err
	}

	return &mpb.StatResponse{
		Entry: re,
	}, nil
}

func (ms *MetaServer) List(ctx context.Context, req *mpb.ListRequest) (*mpb.ListResponse, error) {
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
		return nil, err
	}
	re := make([]*rfspb.Entry, 0)
	for _, entry := range entries {
		res := &rfspb.Entry{}
		if err := entry.ToExistingProtoEntry(res); err != nil {
			return nil, err
		}
		re = append(re, res)
	}
	return &mpb.ListResponse{
		Entries: re,
	}, nil
}
