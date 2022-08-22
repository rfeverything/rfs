package server

import (
	"context"
	"time"

	"github.com/rfeverything/rfs/internal/metrics"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
)

func (m *MetaServer) CreateFile(ctx context.Context, req *mpb.CreateFileRequest) (*mpb.CreateFileResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("CreateFile").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("CreateFile").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.CreateFile(ctx, req)
}

func (m *MetaServer) AppendFile(ctx context.Context, req *mpb.AppendFileRequest) (*mpb.AppendFileResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("AppendFile").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("AppendFile").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.AppendFile(ctx, req)
}

func (m *MetaServer) GetFile(ctx context.Context, req *mpb.GetFileRequest) (*mpb.GetFileResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("GetFile").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("GetFile").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.GetFile(ctx, req)
}

func (m *MetaServer) Stat(ctx context.Context, req *mpb.StatRequest) (*mpb.StatResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("Stat").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("Stat").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.Stat(ctx, req)
}

func (m *MetaServer) Remove(ctx context.Context, req *mpb.RemoveRequest) (*mpb.RemoveResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("Remove").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("Remove").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.Remove(ctx, req)
}

func (m *MetaServer) Mkdir(ctx context.Context, req *mpb.MkdirRequest) (*mpb.MkdirResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("Mkdir").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("Mkdir").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.Mkdir(ctx, req)
}

func (m *MetaServer) List(ctx context.Context, req *mpb.ListRequest) (*mpb.ListResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("List").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("List").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.List(ctx, req)
}

func (m *MetaServer) Move(ctx context.Context, req *mpb.MoveRequest) (*mpb.MoveResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("Move").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("Move").Observe(time.Since(start).Seconds())
	}()
	return m.metaserver.Move(ctx, req)
}
