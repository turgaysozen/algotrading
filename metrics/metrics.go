package metrics

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	// Collect CPU usage metrics
	cpuPercent, err := cpu.Percent(0, true) // 'true' to get per-core usage
	if err != nil {
		log.Println("Error collecting CPU usage:", err)
	}
	for i, percent := range cpuPercent {
		cpuUsage.WithLabelValues(string(i)).Set(percent)
	}

	// Collect memory usage metrics
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error collecting memory usage:", err)
	}
	memoryUsage.Set(float64(vmStat.Used))

	// Optionally, you can call this function periodically
}
