package rfs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"net/url"

	"github.com/google/uuid"
	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/internal/log"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	"github.com/rfeverything/rfs/internal/proto/rfs"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RfsClient struct {
	UUID          uuid.UUID
	MetaServerURL *url.URL
	mpbclient     mpb.MetaServerClient
}

func NewRfsClient() *RfsClient {
	url, err := url.Parse(config.Global().GetString("MetaServerURL"))
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(config.Global().GetString("MetaServerURL"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	suid := config.Global().GetString("UUID")
	var uid uuid.UUID
	if suid != "" {
		uid, err = uuid.Parse(suid)
		if err != nil {
			panic(err)
		}
	} else {
		uid = uuid.New()
		config.Global().Set("UUID", uid.String())
		if err := config.Global().WriteConfig(); err != nil {
			panic(err)
		}
	}

	return &RfsClient{
		MetaServerURL: url,
		mpbclient:     mpb.NewMetaServerClient(conn),
		UUID:          uid,
	}
}

func (rc *RfsClient) GetFile(ctx context.Context, dir string, callback func(io.Reader) error) (fileName string, err error) {
	req := &mpb.GetFileRequest{
		Directory: dir,
		ClientId:  rc.UUID.String(),
	}
	log.Global().Debug("get file", zap.String("dir", dir))
	resp, err := rc.mpbclient.GetFile(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}
	callback(bytes.NewReader(resp.Entry.Content))
	log.Global().Debug("get file done", zap.String("dir", dir))

	return resp.Entry.Name, nil
}

func (rc *RfsClient) PutFile(ctx context.Context, dir string, file fs.File) (err error) {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		panic(err)
	}
	req := &mpb.CreateFileRequest{
		Directory: dir,
		ClientId:  rc.UUID.String(),
		Entry: &rfs.Entry{
			Name: stat.Name(),
			Attributes: &rfspb.FuseAttributes{
				FileSize: uint64(stat.Size()),
				FileMode: uint32(stat.Mode()),
				Crtime:   stat.ModTime().Unix(),
				Mtime:    stat.ModTime().Unix(),
			},
			Content: buf.Bytes(),
		},
	}

	log.Global().Debug("put file", zap.String("dir", dir), zap.String("file", stat.Name()))
	resp, err := rc.mpbclient.CreateFile(ctx, req)
	if err != nil {
		panic(err)
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

func (rc *RfsClient) Mkdir(ctx context.Context, dir string) (err error) {
	req := &mpb.CreateFileRequest{
		Directory: dir,
		ClientId:  rc.UUID.String(),
	}
	log.Global().Debug("mkdir", zap.String("dir", dir))
	resp, err := rc.mpbclient.CreateFile(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	log.Global().Debug("mkdir done", zap.String("dir", dir))
	return nil
}

func (rc *RfsClient) Move(ctx context.Context, srcDir string, dstDir string) (err error) {
	req := &mpb.MoveRequest{
		ClientId:     rc.UUID.String(),
		SrcDirectory: srcDir,
		DstDirectory: dstDir,
	}
	log.Global().Debug("move", zap.String("srcDir", srcDir), zap.String("dstDir", dstDir))
	resp, err := rc.mpbclient.Move(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	log.Global().Debug("move done", zap.String("srcDir", srcDir), zap.String("dstDir", dstDir))
	return nil
}

func (rc *RfsClient) Remove(ctx context.Context, dir string) (err error) {
	req := &mpb.RemoveRequest{
		ClientId:  rc.UUID.String(),
		Directory: dir,
	}
	log.Global().Debug("remove", zap.String("dir", dir))
	resp, err := rc.mpbclient.Remove(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	log.Global().Debug("remove done", zap.String("dir", dir))
	return nil
}

func (rc *RfsClient) List(ctx context.Context, dir string) (es []*rfspb.Entry, err error) {
	req := &mpb.ListRequest{
		ClientId:  rc.UUID.String(),
		Directory: dir,
	}
	log.Global().Debug("list", zap.String("dir", dir))
	resp, err := rc.mpbclient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	log.Global().Debug("list done", zap.String("dir", dir))
	return resp.FileEntries, nil
}

func (rc *RfsClient) Stat(ctx context.Context, dir string) (fileInfo *rfspb.Entry, err error) {

	return nil, nil
}
