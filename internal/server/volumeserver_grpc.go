package server

import (
	"context"
	"time"

	"github.com/rfeverything/rfs/internal/metrics"

	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
)

func (vs *VolumeServer) PutChunk(ctx context.Context, req *vpb.PutChunkRequest) (*vpb.PutChunkResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("PutChunk").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("PutChunk").Observe(time.Since(start).Seconds())
	}()
	return vs.vs.PutChunk(ctx, req)
}
func (vs *VolumeServer) GetChunk(ctx context.Context, req *vpb.GetChunkRequest) (*vpb.GetChunkResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("GetChunk").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("GetChunk").Observe(time.Since(start).Seconds())
	}()
	return vs.vs.GetChunk(ctx, req)
}
func (vs *VolumeServer) DeleteChunk(ctx context.Context, req *vpb.DeleteChunkRequest) (*vpb.DeleteChunkResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("DeleteChunk").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("DeleteChunk").Observe(time.Since(start).Seconds())
	}()
	return vs.vs.DeleteChunk(ctx, req)
}
func (vs *VolumeServer) VolumeStatus(ctx context.Context, req *vpb.VolumeStatusRequest) (*vpb.VolumeStatusResponse, error) {
	metrics.VolumeServerRequestCounter.WithLabelValues("VolumeStatus").Inc()
	start := time.Now()
	defer func() {
		metrics.VolumeServerRequestDuration.WithLabelValues("VolumeStatus").Observe(time.Since(start).Seconds())
	}()
	return vs.vs.VolumeStatus(ctx, req)
}
func (vs *VolumeServer) Ping(ctx context.Context, req *vpb.PingRequest) (*vpb.PingResponse, error) {
	return vs.vs.Ping(ctx, req)
}
