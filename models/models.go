package models

type OrderBook struct {
	LatencyTrackingID string
	EventType         string          `json:"e"`
	Symbol            string          `json:"s"`
	EventTime         int64           `json:"E"`
	Bids              [][]interface{} `json:"b"`
	Asks              [][]interface{} `json:"a"`
}

type Order struct {
	ID        int     `json:"id"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Status    string  `json:"status"`
	OrderType string  `json:"orderType"`
}

type Signal struct {
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	ShortSMA float64 `json:"short_sma"`
	LongSMA  float64 `json:"long_sma"`
	Reason   string  `json:"reason"`
}
