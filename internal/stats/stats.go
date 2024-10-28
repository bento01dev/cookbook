package stats

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type StatsCollection struct {
	serviceName           string
	env                   string
	host                  string
	okRequestGauge        *prometheus.GaugeVec
	badRequestGauge       *prometheus.GaugeVec
	internalErrGauge      *prometheus.GaugeVec
	responseTimeHistogram *prometheus.HistogramVec
}

func newStats(serviceName, env, host string) *StatsCollection {
	okRequestGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Help:      "success response gauge",
			Namespace: "http",
			Name:      "success",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	badRequestGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Help:      "bad request response gauge",
			Namespace: "http",
			Name:      "bad_request",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	internalErrGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Help:      "internal server error gauge",
			Namespace: "http",
			Name:      "internal_error",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	responseTimeHistogram := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Help:      "response time for endpoints",
			Namespace: "http",
			Subsystem: "response",
			Name:      "ms",
			Buckets:   []float64{0.0, 50.0, 100.0, 200.0, 250.0, 300.0, 400.0},
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	return &StatsCollection{
		serviceName:           serviceName,
		env:                   env,
		host:                  host,
		okRequestGauge:        okRequestGauge,
		badRequestGauge:       badRequestGauge,
		internalErrGauge:      internalErrGauge,
		responseTimeHistogram: responseTimeHistogram,
	}
}

func (s *StatsCollection) StatusOkInc(endpoint string) {
	s.okRequestGauge.
		With(prometheus.Labels{"service": s.serviceName, "env": s.env, "host": s.host, "endpoint": endpoint}).
		Inc()
}

func (s *StatsCollection) BadRequestInc(endpoint string) {
	s.badRequestGauge.
		With(prometheus.Labels{"service": s.serviceName, "env": s.env, "host": s.host, "endpoint": endpoint}).
		Inc()
}

func (s *StatsCollection) InternalServerErrorInc(endpoint string) {
	s.internalErrGauge.
		With(prometheus.Labels{"service": s.serviceName, "env": s.env, "host": s.host, "endpoint": endpoint}).
		Inc()
}

func (s *StatsCollection) ResponseTime(endpoint string, responseTimeMs int64) {
	s.responseTimeHistogram.
		With(prometheus.Labels{"service": s.serviceName, "env": s.env, "host": s.host, "endpoint": endpoint}).
		Observe(float64(responseTimeMs))
}

var (
	stats     *StatsCollection
	statsOnce sync.Once
)

func Stats(getEnv func(string) string) *StatsCollection {
	statsOnce.Do(func() {
		serviceName := getEnv("SERVICE_NAME")
		env := getEnv("ENV")
		host := getEnv("HOST_IP")
		stats = newStats(serviceName, env, host)
	})
	return stats
}
