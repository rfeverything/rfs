package metaserver

import (
	"errors"
	"sort"

	"github.com/rfeverything/rfs/internal/constant"
	"github.com/rfeverything/rfs/internal/logger"
	"github.com/rfeverything/rfs/internal/proto/rfs"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

func (ms *MetaServer) assignVolumeForChunks(volumes []*vpb.VolumeStatus, chunks []*rfs.FileChunk) error {
	logger.Global().Debug("assignVolumeForChunks", zap.Any("volumes", volumes), zap.Any("chunks", chunks))
	for _, chunk := range chunks {
		sort.Slice(volumes, func(i, j int) bool {
			return volumes[i].Free < volumes[j].Free
		})
		if chunk.VolumeIds == nil {
			chunk.VolumeIds = make([]string, 0)
		}
		isassigned := false
		replicationCount := constant.Replication
		for _, volume := range volumes {
			if slices.Contains(chunk.VolumeIds, volume.VolumeId) {
				continue
			}
			if volume.Free >= chunk.Size {
				chunk.VolumeIds = append(chunk.VolumeIds, volume.VolumeId)
				volume.Free -= chunk.Size
				replicationCount--
				if replicationCount == 0 {
					//TODO: volumeservers less than replication count
					isassigned = true
				}
				break
			}
		}
		if !isassigned {
			return errors.New("no enough space")
		}
	}
	return nil
}
