package metrics

import (
	"main/pkg/constants"
	"main/pkg/types"
	"main/pkg/utils"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Manager struct {
	logger zerolog.Logger
	config types.MetricsConfig

	registry *prometheus.Registry

	reporterEnabledGauge   *prometheus.GaugeVec
	reporterQueriesCounter *prometheus.CounterVec

	successQueriesCounter *prometheus.CounterVec
	failedQueriesCounter  *prometheus.CounterVec

	appVersionGauge *prometheus.GaugeVec
	startTimeGauge  *prometheus.GaugeVec
}

func NewManager(logger *zerolog.Logger, config types.MetricsConfig) *Manager {
	registry := prometheus.NewRegistry()

	reporterEnabledGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: constants.PrometheusMetricsPrefix + "reporter_enabled",
		Help: "Whether the reporter is enabled (1 if yes, 0 if no)",
	}, []string{"name"})
	reporterQueriesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: constants.PrometheusMetricsPrefix + "reporter_queries",
		Help: "Reporters' queries count ",
	}, []string{"name", "query"})

	successQueriesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: constants.PrometheusMetricsPrefix + "queries_successful",
		Help: "Counter of successful queries towards the external services.",
	}, []string{"chain", "query"})

	failedQueriesCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: constants.PrometheusMetricsPrefix + "queries_failed",
		Help: "Counter of failed queries towards the external services.",
	}, []string{"chain", "query"})

	appVersionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: constants.PrometheusMetricsPrefix + "version",
		Help: "App version",
	}, []string{"version"})
	startTimeGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: constants.PrometheusMetricsPrefix + "start_time",
		Help: "Unix timestamp on when the app was started. Useful for annotations.",
	}, []string{})

	registry.MustRegister(reporterEnabledGauge)
	registry.MustRegister(reporterQueriesCounter)
	registry.MustRegister(successQueriesCounter)
	registry.MustRegister(failedQueriesCounter)
	registry.MustRegister(appVersionGauge)
	registry.MustRegister(startTimeGauge)

	startTimeGauge.
		With(prometheus.Labels{}).
		Set(float64(time.Now().Unix()))

	return &Manager{
		logger:                 logger.With().Str("component", "metrics").Logger(),
		config:                 config,
		registry:               registry,
		reporterEnabledGauge:   reporterEnabledGauge,
		reporterQueriesCounter: reporterQueriesCounter,
		successQueriesCounter:  successQueriesCounter,
		failedQueriesCounter:   failedQueriesCounter,
		appVersionGauge:        appVersionGauge,
		startTimeGauge:         startTimeGauge,
	}
}

func (m *Manager) Start() {
	if !m.config.Enabled.Bool {
		m.logger.Info().Msg("Metrics not enabled")
		return
	}

	m.logger.Info().
		Str("addr", m.config.ListenAddr).
		Msg("Metrics handler listening")

	http.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry}))
	if err := http.ListenAndServe(m.config.ListenAddr, nil); err != nil {
		m.logger.Panic().
			Err(err).
			Str("addr", m.config.ListenAddr).
			Msg("Cannot start metrics handler")
	}
}

func (m *Manager) LogReporterQuery(reporter string, query string) {
	m.reporterQueriesCounter.
		With(prometheus.Labels{
			"name":  reporter,
			"query": query,
		}).
		Inc()
}

func (m *Manager) LogReporterEnabled(name string, enabled bool) {
	m.reporterEnabledGauge.
		With(prometheus.Labels{"name": name}).
		Set(utils.BoolToFloat64(enabled))
}

func (m *Manager) LogAppVersion(version string) {
	m.appVersionGauge.
		With(prometheus.Labels{"version": version}).
		Set(1)
}

func (m *Manager) LogQueryInfo(queryInfo types.QueryInfo) {
	if queryInfo.Success {
		m.successQueriesCounter.
			With(prometheus.Labels{"chain": queryInfo.Chain, "query": queryInfo.Query}).
			Inc()
	} else {
		m.failedQueriesCounter.
			With(prometheus.Labels{"chain": queryInfo.Chain, "query": queryInfo.Query}).
			Inc()
	}
}
