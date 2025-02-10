package services

func CheckSignal(shortSMA, longSMA float64, lastSignal string) (string, string) {
	if shortSMA > longSMA && lastSignal != "BUY" {
		return "BUY Signal!", "Short SMA crossed above Long SMA"
	} else if shortSMA < longSMA && lastSignal != "SELL" {
		return "SELL Signal!", "Short SMA crossed below Long SMA"
	}
	return "NO Signal", "No significant crossover"
}

func CalculateSMA(data []float64, period int) float64 {
	if len(data) < period {
		return 0
	}

	var sum float64
	for i := len(data) - period; i < len(data); i++ {
		sum += data[i]
	}

	return sum / float64(period)
}
