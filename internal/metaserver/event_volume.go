package metaserver

import vpb "github.com/rfeverything/rfs/internal/proto/volume_server"

type VolumeUpdateEvent struct {
	Type     int
	Status   *vpb.VolumeStatus
	VolumeId string
}

const (
	VolumeEventTypePut = iota
	VolumeEventTypeDelete
)
