package metaserver

import (
	"errors"

	"github.com/rfeverything/rfs/internal/proto/rfs"
)

func (e *Entry) SplitToChunks(chunkSize uint64) error {
	if e.Chunks != nil {
		return errors.New("entry already has chunks")
	}
	if e.Content == nil || len(e.Content) == 0 {
		return errors.New("entry has no content")
	}
	e.Chunks = make([]*rfs.FileChunk, 0)
	for i := uint64(0); i < uint64(len(e.Content)); i += chunkSize {
		chunk := &rfs.FileChunk{
			Offset: i,
			Size:   chunkSize,
		}
		if i+chunkSize > uint64(len(e.Content)) {
			chunk.Size = uint64(len(e.Content)) - i
		}
		chunk.Content = e.Content[i : i+chunk.Size]
		e.Chunks = append(e.Chunks, chunk)
	}
	return nil
}

func (e *Entry) CombineChunksGetContent() error {
	if e.Chunks == nil || len(e.Chunks) == 0 {
		return errors.New("entry has no chunks")
	}
	if e.Content != nil {
		return errors.New("entry already has content")
	}
	e.Content = make([]byte, 0)
	for _, chunk := range e.Chunks {
		e.Content = append(e.Content, chunk.Content...)
	}
	return nil
}
