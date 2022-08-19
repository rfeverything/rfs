package server

import vpb "github.com/rfeverything/rfs/internal/proto/volume_server"

type VolumeServer struct {
	*vpb.UnimplementedVolumeServerServer
}
