package metrics

import (
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var (
	orderbookProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "orderbook_processing_latency_seconds",
			Help:    "Latency of processing order book data",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"latencyTrackingID"},
	)

	tradeSignalLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "trade_signal_latency_seconds",
			Help:    "Latency of generating trade signals",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"latencyTrackingID"},
	)

	orderExecutionLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_execution_latency_seconds",
			Help:    "Latency of executing trade orders",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"latencyTrackingID"},
	)

	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage percentage",
		},
		[]string{"cpu"},
	)

	memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	activeTimers = make(map[string]time.Time)
)

func init() {
	prometheus.MustRegister(
		orderbookProcessingLatency,
		tradeSignalLatency,
		orderExecutionLatency,
		cpuUsage,
		memoryUsage,
	)
}

func SetStartTime(metricType, latencyTrackingID string) {
	activeTimers[metricType+latencyTrackingID] = time.Now()
}

func RecordLatency(metricType, latencyTrackingID string) {
	startTime, exists := activeTimers[metricType+latencyTrackingID]
	if !exists {
		return
	}

	latency := time.Since(startTime).Seconds()

	switch metricType {
	case "orderbook":
		orderbookProcessingLatency.WithLabelValues(latencyTrackingID).Observe(latency)
	case "signal":
		tradeSignalLatency.WithLabelValues(latencyTrackingID).Observe(latency)
	case "order":
		orderExecutionLatency.WithLabelValues(latencyTrackingID).Observe(latency)
	}

	delete(activeTimers, metricType+latencyTrackingID)
}

func CollectSystemMetrics() {
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		log.Println("Error collecting CPU usage:", err)
	}
	for i, percent := range cpuPercent {
		cpuUsage.WithLabelValues(strconv.Itoa(i)).Set(percent)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsage.Set(float64(m.Sys))
}
