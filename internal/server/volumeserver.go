package server

import (
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"github.com/rfeverything/rfs/internal/volume"
)

type VolumeServer struct {
	*vpb.UnimplementedVolumeServerServer
	vs *volume.VolumeServer
}

func NewVolumeServer() (*VolumeServer, error) {
	return &VolumeServer{
		vs: volume.NewVolumeServer(),
	}, nil
}
