package metaserver

import (
	"time"

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
		for cnt := 0; cnt < 3; cnt++ {
			logger.Global().Debug("watchVolumeState", zap.String("volumeId", st.VolumeId), zap.Int64("free", int64(st.Free)))
			vc, err := NewVolumeClient(st.Address)
			if err != nil {
				logger.Global().Error("watchVolumeState", zap.Error(err), zap.String("volumeId", st.VolumeId), zap.Int64("free", int64(st.Free)), zap.String("address", st.Address), zap.Int("cnt", cnt))
				time.Sleep(time.Second)
				continue
			}
			ms.VolumeClients[st.VolumeId] = vc
			break
		}
	}

	for event := range ms.Store.GetVolumesStatusChan() {
		logger.Global().Info("watchVolumeState", zap.Any("event", event))
		if event.VolumeId == "" {
			continue
		}
		switch event.Type {
		case VolumeEventTypePut:
			for cnt := 0; cnt < 3; cnt++ {
				logger.Global().Info("watchVolumeState", zap.String("volumeId", event.VolumeId))
				vc, err := NewVolumeClient(event.Status.Address)
				if err != nil {
					logger.Global().Error("watchVolumeState", zap.Error(err), zap.String("volumeId", event.VolumeId), zap.String("address", event.Status.Address), zap.Int("cnt", cnt))
					time.Sleep(time.Second)
					continue
				}
				ms.VolumeClients[event.VolumeId] = vc
				break
			}
		case VolumeEventTypeDelete:
			delete(ms.VolumeClients, event.VolumeId)
		}
	}
}
