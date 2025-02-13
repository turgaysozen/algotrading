package utils

import (
	"strconv"

	"github.com/turgaysozen/algotrading/monitoring/metrics"
)

func StringToFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		metrics.RecordError("StringToFloat64_error")
		return 0
	}
	return val
}
