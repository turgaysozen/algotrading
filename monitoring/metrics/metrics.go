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

	orderExecutionAvgLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_execution_avg_processing_latency_seconds",
			Help: "Average latency of processing order execution data",
		},
	)

	tradingSignalAvgLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "trading_signal_avg_processing_latency_seconds",
			Help: "Average latency of processing trading signal data",
		},
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

	activeTimers   = make(map[string]time.Time)
	activeTimersMu sync.Mutex
	latencySums    = make(map[string]float64)
	latencyCounts  = make(map[string]int)
	latencyMutex   = sync.Mutex{}
)

func init() {
	prometheus.MustRegister(
		orderbookAvgLatency,
		orderExecutionAvgLatency,
		tradingSignalAvgLatency,
		cpuUsage,
		memoryUsage,
		errors,
		dataLoss,
	)
}

func SetStartTime(metricType string) {
	activeTimersMu.Lock()
	activeTimers[metricType] = time.Now()
	activeTimersMu.Unlock()
}

func RecordLatency(metricType string) {
	activeTimersMu.Lock()
	startTime, exists := activeTimers[metricType]
	if exists {
		delete(activeTimers, metricType)
	}
	activeTimersMu.Unlock()

	if !exists {
		return
	}

	latency := time.Since(startTime).Seconds()

	switch metricType {
	case "orderbook_avg":
		orderbookAvgLatency.Set(updateAverageLatency("avg", latency))
	case "signal_avg":
		tradingSignalAvgLatency.Set(updateAverageLatency("avg", latency))
	case "order_avg":
		orderExecutionAvgLatency.Set(updateAverageLatency("avg", latency))
	}
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
