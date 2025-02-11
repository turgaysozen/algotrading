package models

import (
	"database/sql"
	"log"
	"time"
)

func SaveOrderBook(db *sql.DB, orderBook OrderBook) error {
	_, err := db.Exec(
		"INSERT INTO order_books (symbol, timestamp, bids, asks) VALUES ($1, $2, $3, $4)",
		orderBook.Symbol, orderBook.EventTime, orderBook.Bids, orderBook.Asks,
	)
	if err != nil {
		log.Printf("Error saving order book: %v", err)
		return err
	}
	return nil
}

func SaveOrder(db *sql.DB, order Order) error {
	_, err := db.Exec(
		"INSERT INTO orders (id, price, quantity, status, order_type) VALUES ($1, $2, $3, $4, $5)",
		order.ID, order.Price, order.Quantity, order.Status, order.OrderType,
	)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		return err
	}
	return nil
}

func SaveSignal(db *sql.DB, signal Signal) error {
	_, err := db.Exec(
		"INSERT INTO signals (timestamp, type, price, short_sma, long_sma, reason) VALUES ($1, $2, $3, $4, $5, $6)",
		signal.Timestamp, signal.Type, signal.Price, signal.ShortSMA, signal.LongSMA, signal.Reason,
	)
	if err != nil {
		log.Printf("Error saving signal: %v", err)
		return err
	}
	return nil
}

func GetOrder(db *sql.DB, orderID int) (*Order, error) {
	var order Order
	err := db.QueryRow(
		"SELECT id, price, quantity, status, order_type FROM orders WHERE id = $1", orderID,
	).Scan(&order.ID, &order.Price, &order.Quantity, &order.Status, &order.OrderType)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return nil, err
	}
	return &order, nil
}

func GetSignal(db *sql.DB, timestamp time.Time) (*Signal, error) {
	var signal Signal
	err := db.QueryRow(
		"SELECT timestamp, type, price, short_sma, long_sma, reason FROM signals WHERE timestamp = $1", timestamp,
	).Scan(&signal.Timestamp, &signal.Type, &signal.Price, &signal.ShortSMA, &signal.LongSMA, &signal.Reason)
	if err != nil {
		log.Printf("Error fetching signal: %v", err)
		return nil, err
	}
	return &signal, nil
}
