package models

import "time"

type OrderBook struct {
	Bids [][]interface{} `json:"b"`
	Asks [][]interface{} `json:"a"`
}

type Order struct {
	ID        int     `json:"id"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Status    string  `json:"status"`
	OrderType string  `json:"orderType"`
}

type Signal struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	ShortSMA  float64   `json:"short_sma"`
	LongSMA   float64   `json:"long_sma"`
	Reason    string    `json:"reason"`
}
