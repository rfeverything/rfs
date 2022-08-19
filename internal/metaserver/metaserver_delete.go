package metaserver

import (
	"context"
	"errors"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

func (ms *MetaServer) Remove(ctx context.Context, req *mpb.RemoveRequest) (*mpb.RemoveResponse, error) {
	pth := req.GetPath()
	e, err := ms.Store.GetEntry(ctx, pth)
	if err != nil {
		return nil, err
	}
	if e.Chunks != nil {
		for _, chunk := range e.Chunks {
			for _, volumeserver := range chunk.VolumeIds {
				resp, err := ms.VolumeClients[volumeserver].DeleteChunk(ctx, &vpb.DeleteChunkRequest{
					ChunkId: chunk.Chunkid,
				})
				if err != nil {
					return nil, err
				}
				if resp.Error != "" {
					return nil, errors.New(resp.Error)
				}
				break
			}
		}
	}
	if err := ms.Store.DeleteEntry(ctx, pth); err != nil {
		return nil, err
	}
	return nil, nil
}
