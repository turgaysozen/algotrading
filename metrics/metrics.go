package metrics

import (
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var (
	orderbookAvgLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "orderbook_avg_processing_latency_seconds",
			Help: "Average latency of processing order book data",
		},
	)

	orderbookSingleProcessingLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "orderbook_single_processing_latency_seconds",
			Help: "Latency of processing single order book data",
		},
		[]string{"latencyTrackingID"},
	)

	tradeSignalLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "trade_signal_latency_seconds",
			Help: "Latency of generating trade signals",
		},
		[]string{"latencyTrackingID"},
	)

	orderExecutionLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "order_execution_latency_seconds",
			Help: "Latency of executing trade orders",
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

	errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_count",
			Help: "Total number of errors in the system",
		},
		[]string{"error_type"},
	)

	dataLoss = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dataloss_error_count",
			Help: "Total number of data loss events",
		},
		[]string{"data_loss_type"},
	)

	activeTimers  = make(map[string]time.Time)
	latencySums   = make(map[string]float64)
	latencyCounts = make(map[string]int)
	latencyMutex  = sync.Mutex{}
)

func init() {
	prometheus.MustRegister(
		orderbookAvgLatency,
		orderbookSingleProcessingLatency,
		tradeSignalLatency,
		orderExecutionLatency,
		cpuUsage,
		memoryUsage,
		errors,
		dataLoss,
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
	case "orderbook_avg":
		orderbookAvgLatency.Set(updateAverageLatency("avg", latency))
	case "orderbook_single":
		orderbookSingleProcessingLatency.WithLabelValues(latencyTrackingID).Set(latency)
	case "signal":
		tradeSignalLatency.WithLabelValues(latencyTrackingID).Set(latency)
	case "order":
		orderExecutionLatency.WithLabelValues(latencyTrackingID).Set(latency)
	}

	delete(activeTimers, metricType+latencyTrackingID)
}

func updateAverageLatency(metricType string, latency float64) float64 {
	latencyMutex.Lock()
	defer latencyMutex.Unlock()

	latencySums[metricType] += latency
	latencyCounts[metricType]++

	avgLatency := latencySums[metricType] / float64(latencyCounts[metricType])
	return avgLatency
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

func RecordError(errorType string) {
	errors.WithLabelValues(errorType).Inc()
}

func RecordDataLoss(dataLossType string) {
	dataLoss.WithLabelValues(dataLossType).Inc()
}
