package swagger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func BuildSummaryVec(metricName string, metricHelp string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        metricName,
			Help:        metricHelp,
			ConstLabels: prometheus.Labels{"service": "secret_api"},
		},
		[]string{"handler", "code"},
	)
	prometheus.Register(summaryVec)
	return summaryVec
}

// WithMonitoring optionally adds a middleware that stores request duration and response size into the supplied
// summaryVec
func WithMonitoring(next http.Handler, route Route, summary *prometheus.SummaryVec) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := NewMonitoringResponseWriter(rw)
		next.ServeHTTP(lrw, req)
		statusCode := lrw.statusCode
		duration := time.Since(start)

		// Store duration of request
		summary.WithLabelValues(route.Name, strconv.FormatInt(int64(statusCode), 10)).Observe(duration.Seconds() * 1000)
	})
}

type monitoringResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMonitoringResponseWriter(w http.ResponseWriter) *monitoringResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &monitoringResponseWriter{w, http.StatusOK}
}

func (lrw *monitoringResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
