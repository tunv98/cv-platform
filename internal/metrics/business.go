package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CVUploadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cv_uploads_total",
			Help: "Total number of CV upload operations",
		},
		[]string{"operation", "status"},
	)

	CVUploadDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cv_upload_duration_seconds",
			Help:    "Duration of CV upload operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	ProfileRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "profile_requests_total",
			Help: "Total number of profile requests",
		},
		[]string{"operation", "status"},
	)
)

// RecordCVUpload records CV upload business metrics
func RecordCVUpload(operation, status string, duration float64) {
	CVUploadsTotal.WithLabelValues(operation, status).Inc()
	if duration > 0 {
		CVUploadDuration.WithLabelValues(operation).Observe(duration)
	}
}

// RecordProfileRequest records profile request business metrics
func RecordProfileRequest(operation, status string) {
	ProfileRequestsTotal.WithLabelValues(operation, status).Inc()
}
