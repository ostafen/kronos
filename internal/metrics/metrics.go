package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var webhookRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "webhook_requests_total",
		Help: "Total number of webhook requests",
	},
	[]string{"url", "code"},
)

var scheduleFailures = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "webhook_notifications",
		Help: "Total number of webhook notifications",
	},
	[]string{"id"},
)

// schedule_failures_total, schedule_success_total

func init() {
	prometheus.MustRegister(webhookRequestsTotal, scheduleFailures)
}

func IncWebhookRequests(url string, code int) {
	webhookRequestsTotal.WithLabelValues(url, strconv.Itoa(code)).Inc()
}

func IncScheduleFailures(id string) {
	scheduleFailures.WithLabelValues(id).Inc()
}

func ResetScheduleFailures(id string) {
	scheduleFailures.WithLabelValues(id).Set(0)
}
