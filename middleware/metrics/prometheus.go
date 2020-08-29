package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var summaryVector *prometheus.SummaryVec
var counterVector *prometheus.CounterVec
var creatorLock sync.Mutex

func getSummaryVector() *prometheus.SummaryVec {
	creatorLock.Lock()
	defer creatorLock.Unlock()

	if summaryVector != nil {
		return summaryVector
	} else {
		summaryVector = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "pac_backend_api",
				Name:      "api_requests_summary",
				Help:      "Info about the API requests processed total",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service"},
		)
		prometheus.MustRegister(summaryVector)
		return summaryVector
	}
}

func getRequestsCounterVec() *prometheus.CounterVec {
	creatorLock.Lock()
	defer creatorLock.Unlock()

	if counterVector != nil {
		return counterVector
	} else {
		counterVector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "pac_backend_api",
				Name: "api_requests_total",
				Help: "How many API requests processed, partitioned by status",
			},
			[]string{"status"},
		)

		prometheus.MustRegister(counterVector)
		return counterVector
	}
}

func Prometheus(next http.Handler) http.Handler {

	counterVec := getRequestsCounterVec()
	summaryVec := getSummaryVector()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(rw)

		start := time.Now()
		next.ServeHTTP(lrw, r) // the call to next in the chain (with the wrapped ResponseWriter!)
		duration := time.Since(start)

		// Store duration of request
		summaryVec.WithLabelValues("duration").Observe(duration.Seconds())

		// Store size of response, if possible
		size, err := strconv.Atoi(rw.Header().Get("Content-Length"))
		if err == nil {
			summaryVec.WithLabelValues("size").Observe(float64(size))
		}

		// Store the request success/error count
		if lrw.statusCode >= 200 && lrw.statusCode < 300 {
			counterVec.WithLabelValues("success").Inc()
		} else {
			counterVec.WithLabelValues("error").Inc()
		}
	})
}
