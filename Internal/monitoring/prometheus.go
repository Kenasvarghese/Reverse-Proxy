package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Observer struct {
	ReqCount                prometheus.Counter
	ReqDuration             prometheus.Histogram
	TimeForRatelimiterCheck *prometheus.HistogramVec
	InflightReq             prometheus.Gauge
	IpRateLimitedCount      prometheus.Counter
	ActiveIpCount           prometheus.Gauge
	GlobalRateLimitedCount  prometheus.Counter
}

var Buckets = []float64{
	10,
	20,
	50,
	100,
	500,
	1000,
	5000,
	10000,
	50000,
	100000,
}

func NewObserver() *Observer {
	reqCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_request_count",
		Help: "Total number of requests",
	})
	reqDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "proxy_request_duration_seconds",
		Help:    "Total number of requests",
		Buckets: Buckets,
	})
	timeForRatelimiterCheck := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "proxy_ratelimiter_check_time_seconds",
		Help:    "Time taken for ratelimiter check",
		Buckets: Buckets,
	}, []string{"status"})
	inflightReq := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_request_inflight_total",
		Help: "Total number of requests",
	})
	ipRateLimitedCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_ip_rate_limited_count",
		Help: "Total number of requests that were rate limited by IP",
	})
	activeIpCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_active_ip_count",
		Help: "Total number of active IPs",
	})
	globalRateLimitedCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_global_rate_limited_count",
		Help: "Total number of requests that were rate limited by global",
	})
	prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)

	prometheus.MustRegister(reqCount, reqDuration, timeForRatelimiterCheck, inflightReq, ipRateLimitedCount, activeIpCount, globalRateLimitedCount)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatal("Failed to start metrics server: ", err)
		}
	}()
	return &Observer{
		ReqCount:                reqCount,
		ReqDuration:             reqDuration,
		TimeForRatelimiterCheck: timeForRatelimiterCheck,
		InflightReq:             inflightReq,
		IpRateLimitedCount:      ipRateLimitedCount,
		ActiveIpCount:           activeIpCount,
		GlobalRateLimitedCount:  globalRateLimitedCount,
	}
}
