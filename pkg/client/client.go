package rfs

import (
	"context"
	"io"
	"net/url"

	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

type RfsClient struct {
	MetaServerURL *url.URL
	mpbclient     mpb.MetaServerClient
}

func (rc *RfsClient) GetFile(ctx context.Context, dir string, callback func(io.Reader) error) (fileName string, err error) {
	return nil, nil
}

func (rc *RfsClient) PutFile(ctx context.Context, dir string, fileName string, file io.Reader) (err error) {
	return nil
}

func (rc *RfsClient) Mkdir(ctx context.Context, dir string) (err error) {
	return nil
}

func (rc *RfsClient) Move(ctx context.Context, srcDir string, dstDir string) (err error) {
	return nil
}

func (rc *RfsClient) Remove(ctx context.Context, dir string) (err error) {
	return nil
}

func (rc *RfsClient) List(ctx context.Context, dir string) (fileNames []string, err error) {
	return nil, nil
}

func (rc *RfsClient) Stat(ctx context.Context, dir string) (fileInfo *mpb.FileInfo, err error) {
	return nil, nil
}
