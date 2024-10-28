package stats

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type StatsCollection struct {
	serviceName      string
	env              string
	host             string
	okRequestGauge   *prometheus.GaugeVec
	badRequestGauge  *prometheus.GaugeVec
	internalErrGauge *prometheus.GaugeVec
}

func newStats(serviceName, env, host string) *StatsCollection {
	okRequestGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "http",
			Name:      "success",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	badRequestGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "http",
			Name:      "bad_request",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	internalErrGauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "http",
			Name:      "internal_error",
		},
		[]string{"service", "env", "host", "endpoint"},
	)

	return &StatsCollection{
		serviceName:      serviceName,
		env:              env,
		host:             host,
		okRequestGauge:   okRequestGauge,
		badRequestGauge:  badRequestGauge,
		internalErrGauge: internalErrGauge,
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
