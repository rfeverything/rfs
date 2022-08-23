package metaserver

import (
	"math/rand"

	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/internal/logger"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	"go.uber.org/zap"
)

type MetaServer struct {
	Store         *EtcdStore
	UniqueId      int32
	VolumeClients map[string]vpb.VolumeServerClient // key: volumeId
}

func NewMetaServer() (*MetaServer, error) {
	UniqueID := config.Global().GetInt32("metaserver.unique_id")
	if UniqueID == 0 {
		UniqueID = rand.Int31()
		config.Global().Set("metaserver.unique_id", UniqueID)
		if err := config.Global().WriteConfig(); err != nil {
			logger.Global().Fatal(err.Error())
		}
		logger.Global().Info("NewMetaServer generate new uid", zap.Int32("uid", UniqueID))
	} else {
		logger.Global().Info("NewMetaServer use store uid", zap.Int32("uid", UniqueID))
	}
	Store := NewEtcdStore(UniqueID)
	logger.Global().Info("NewMetaServer", zap.Int32("UniqueID", UniqueID))
	ms := &MetaServer{
		Store:    Store,
		UniqueId: UniqueID,
	}
	ms.VolumeClients = make(map[string]vpb.VolumeServerClient)
	go ms.watchVolumeState()
	return ms, nil
}
