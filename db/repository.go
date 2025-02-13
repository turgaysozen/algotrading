package db

import (
	"database/sql"
	"log"

	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/monitoring/metrics"
)

func SaveOrderBook(eventType, symbol string, eventTime int64, bestBid, bestAsk float64) (int64, error) {
	query := `
        INSERT INTO order_books (event_type, symbol, event_time, best_bid, best_ask)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id, event_time)
        DO UPDATE SET 
            event_type = EXCLUDED.event_type,
            symbol = EXCLUDED.symbol,
            best_bid = EXCLUDED.best_bid,
            best_ask = EXCLUDED.best_ask
        RETURNING id
    `

	var orderBookID int64
	err := Database.QueryRow(query, eventType, symbol, eventTime, bestBid, bestAsk).Scan(&orderBookID)
	if err != nil {
		log.Printf("Error saving order book: %v", err)
		metrics.RecordError("db_save_order_book_error")
		return 0, err
	}

	return orderBookID, nil
}

func SaveOrder(order models.Order) error {
	query := `
		INSERT INTO orders (price, quantity, status, order_type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, created_at) DO UPDATE SET
			price = EXCLUDED.price,
			quantity = EXCLUDED.quantity,
			status = EXCLUDED.status,
			order_type = EXCLUDED.order_type
	`
	_, err := Database.Exec(query, order.Price, order.Quantity, order.Status, order.OrderType)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		metrics.RecordError("db_save_order_error")
		return err
	}
	return nil
}

func GetLastOpenOrder() (*models.Order, error) {
	var order models.Order
	query := `
		SELECT id, price, quantity, status, order_type
		FROM orders
		WHERE status = 'open'
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := Database.QueryRow(query).Scan(&order.ID, &order.Price, &order.Quantity, &order.Status, &order.OrderType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error retrieving last open order: %v", err)
		metrics.RecordError("db_get_last_open_order_error")
		return nil, err
	}

	return &order, nil
}

func CloseOrder(orderID int) error {
	query := `
		UPDATE orders
		SET status = 'closed', updated_at = NOW()
		WHERE id = $1
	`
	_, err := Database.Exec(query, orderID)
	if err != nil {
		log.Printf("Error closing order with ID %d: %v", orderID, err)
		metrics.RecordError("db_close_order_error")
		return err
	}

	log.Printf("Order with ID %d has been closed.\n", orderID)
	return nil
}

func SaveSignal(signal models.Signal) error {
	query := `
		INSERT INTO signals (type, price, short_sma, long_sma, reason)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id, timestamp) DO UPDATE SET
			type = EXCLUDED.type,
			price = EXCLUDED.price,
			short_sma = EXCLUDED.short_sma,
			long_sma = EXCLUDED.long_sma,
			reason = EXCLUDED.reason
	`
	_, err := Database.Exec(query, signal.Type, signal.Price, signal.ShortSMA, signal.LongSMA, signal.Reason)
	if err != nil {
		log.Printf("Error saving signal: %v", err)
		metrics.RecordError("db_save_signal_error")
		return err
	}
	return nil
}
