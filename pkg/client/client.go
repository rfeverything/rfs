package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"net/url"

	"github.com/google/uuid"
	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/internal/logger"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RfsClient struct {
	UUID          uuid.UUID
	MetaServerURL *url.URL
	mpbclient     mpb.MetaServerClient
}

func NewRfsClient() (*RfsClient, error) {
	url, err := url.Parse(config.Global().GetString("client.server"))
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(config.Global().GetString("client.server")+":"+config.Global().GetString("client.port"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	suid := config.Global().GetString("client.uuid")
	var uid uuid.UUID
	if suid != "" {
		uid, err = uuid.Parse(suid)
		if err != nil {
			return nil, err
		}
	} else {
		uid = uuid.New()
		config.Global().Set("client.uuid", uid.String())
		if err := config.Global().WriteConfig(); err != nil {
			return nil, err
		}
	}

	logger.Global().Debug("client uuid", zap.String("uuid", uid.String()))

	return &RfsClient{
		MetaServerURL: url,
		mpbclient:     mpb.NewMetaServerClient(conn),
		UUID:          uid,
	}, nil
}

func (rc *RfsClient) GetFile(ctx context.Context, path string) (fileName string, file io.Writer, err error) {
	req := &mpb.GetFileRequest{
		Path:     path,
		ClientId: rc.UUID.String(),
	}
	logger.Global().Debug("get file", zap.String("path", path))
	resp, err := rc.mpbclient.GetFile(ctx, req)
	if err != nil {
		return "", nil, err
	}
	if resp.Error != "" {
		return "", nil, errors.New(resp.Error)
	}
	logger.Global().Debug("get file done", zap.String("path", path))

	return resp.Entry.Name, bytes.NewBuffer(resp.Entry.Content), nil
}

func (rc *RfsClient) PutFile(ctx context.Context, path string, file fs.File) (err error) {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		panic(err)
	}
	logger.Global().Debug("put file", zap.String("path", path))
	req := &mpb.CreateFileRequest{
		Path:     path,
		ClientId: rc.UUID.String(),
		Entry: &rfspb.Entry{
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

	logger.Global().Debug("put file", zap.String("path", path), zap.String("file", stat.Name()))
	resp, err := rc.mpbclient.CreateFile(ctx, req)
	if err != nil {
		panic(err)
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

func (rc *RfsClient) Mkdir(ctx context.Context, path string) (err error) {
	req := &mpb.CreateFileRequest{
		Path:     path,
		ClientId: rc.UUID.String(),
	}
	logger.Global().Debug("mkdir", zap.String("path", path))
	resp, err := rc.mpbclient.CreateFile(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	logger.Global().Debug("mkdir done", zap.String("path", path))
	return nil
}

func (rc *RfsClient) Move(ctx context.Context, srcPath string, dstPath string) (err error) {
	req := &mpb.MoveRequest{
		ClientId: rc.UUID.String(),
		SrcPath:  srcPath,
		DstPath:  dstPath,
	}
	logger.Global().Debug("move", zap.String("srcDir", srcPath), zap.String("dstDir", dstPath))
	resp, err := rc.mpbclient.Move(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	logger.Global().Debug("move done", zap.String("srcDir", srcPath), zap.String("dstDir", dstPath))
	return nil
}

func (rc *RfsClient) Remove(ctx context.Context, path string) (err error) {
	req := &mpb.RemoveRequest{
		ClientId: rc.UUID.String(),
		Path:     path,
	}
	logger.Global().Debug("remove", zap.String("path", path))
	resp, err := rc.mpbclient.Remove(ctx, req)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	logger.Global().Debug("remove done", zap.String("path", path))
	return nil
}

func (rc *RfsClient) List(ctx context.Context, dir string) (es []*rfspb.Entry, err error) {
	req := &mpb.ListRequest{
		ClientId:  rc.UUID.String(),
		Directory: dir,
	}
	logger.Global().Debug("list", zap.String("dir", dir))
	resp, err := rc.mpbclient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	logger.Global().Debug("list done", zap.String("dir", dir))
	return resp.Entries, nil
}

func (rc *RfsClient) Stat(ctx context.Context, dir string) (fileInfo *rfspb.Entry, err error) {

	return nil, nil
}
