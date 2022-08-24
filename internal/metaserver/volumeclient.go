package metaserver

import (
	"github.com/rfeverything/rfs/internal/logger"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewVolumeClient(addr string) (vpb.VolumeServerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Global().Fatal(err.Error())
		return nil, err
	}
	logger.Global().Debug("new volume client", zap.String("addr", addr))
	return vpb.NewVolumeServerClient(conn), nil
}
