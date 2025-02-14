package services

func CheckSignal(shortSMA, longSMA float64, lastSignal string) (string, string) {
	if shortSMA > longSMA && lastSignal != "BUY" {
		return "BUY Signal!", "Short SMA crossed above Long SMA"
	} else if shortSMA < longSMA && lastSignal != "SELL" {
		return "SELL Signal!", "Short SMA crossed below Long SMA"
	}
	return "NO Signal", "No significant crossover"
}

// TODO: optimize SMA calculation from O(n) complexity to O(1) by using sliding window

// func CalculateSMA(data []float64, period int) float64 {
// 	if len(data) < period {
// 		return 0
// 	}

// 	var sum float64
// 	for i := len(data) - period; i < len(data); i++ {
// 		sum += data[i]
// 	}

// 	return sum / float64(period)
// }

// Optimize SMA calculation to O(1) complexity by using sliding window
type SMA struct {
	window []float64
	period int
	sum    float64
}

func NewSMA(period int) *SMA {
	return &SMA{
		window: make([]float64, 0, period),
		period: period,
		sum:    0,
	}
}

func (s *SMA) AddPrice(price float64) float64 {
	s.window = append(s.window, price)
	s.sum += price

	if len(s.window) > s.period {
		s.sum -= s.window[0]
		s.window = s.window[1:]
	}

	if len(s.window) < s.period {
		return 0
	}

	return s.sum / float64(s.period)
}
