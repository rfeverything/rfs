package metaserver

import (
	"github.com/rfeverything/rfs/internal/logger"
	"go.uber.org/zap"
)

func (ms *MetaServer) watchVolumeState() {
	sts, err := ms.Store.GetVolumesStatus()
	if err != nil {
		logger.Global().Error("failed to get volumes status", zap.Error(err))
		return
	}
	for _, st := range sts {
		ms.VolumeClients[st.VolumeId] = NewVolumeClient(st.Address)
	}

	for event := range ms.Store.GetVolumesStatusChan() {
		logger.Global().Info("watchVolumeState", zap.Any("event", event))
		if event.VolumeId == "" {
			continue
		}
		switch event.Type {
		case VolumeEventTypePut:
			ms.VolumeClients[event.VolumeId] = NewVolumeClient(event.Status.Address)
		case VolumeEventTypeDelete:
			delete(ms.VolumeClients, event.VolumeId)
		}
	}
}
