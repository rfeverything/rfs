package metrics

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MetaServerRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_server_request_total",
			Help: "Total number of requests received by the meta server",
		},
		[]string{"type"},
	)
	MetaServerRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "meta_server_request_duration_seconds",
			Help:    "Duration of requests received by the meta server",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 24),
		},
		[]string{"type", "status"},
	)
	MetaServerStoreCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_server_store_total",
			Help: "Total number of stores received by the meta server",
		},
		[]string{"store", "type"},
	)
	MetaServerStoreDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "meta_server_store_duration_seconds",
			Help:    "Duration of stores received by the meta server",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 24),
		},
		[]string{"store", "type", "status"},
	)
	VolumeServerRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volume_server_request_total",
			Help: "Total number of requests received by the volume server",
		},
		[]string{"type"},
	)
	VolumeServerRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "volume_server_request_duration_seconds",
			Help:    "Duration of requests received by the volume server",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 24),
		},
		[]string{"type", "status"},
	)
	VolumeServerStoreCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volume_server_store_total",
			Help: "Total number of stores received by the volume server",
		},
		[]string{"store", "type"},
	)
	VolumeServerStoreDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "volume_server_store_duration_seconds",
			Help:    "Duration of stores received by the volume server",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 24),
		},
		[]string{"store", "type", "status"},
	)
	VolumeServerDiskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "volume_server_disk_usage_bytes",
			Help: "Disk usage of the volume server",
		},
		[]string{"store"},
	)
)

func init() {
	prometheus.MustRegister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	prometheus.MustRegister(MetaServerRequestCounter)
	prometheus.MustRegister(MetaServerRequestDuration)
	prometheus.MustRegister(MetaServerStoreCounter)
	prometheus.MustRegister(MetaServerStoreDuration)
	prometheus.MustRegister(VolumeServerRequestCounter)
	prometheus.MustRegister(VolumeServerRequestDuration)
	prometheus.MustRegister(VolumeServerStoreCounter)
	prometheus.MustRegister(VolumeServerStoreDuration)
	prometheus.MustRegister(VolumeServerDiskUsage)
}

func StartMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(strings.Join([]string{":", port}, ""), nil)
}