package services

import (
	"fmt"
	"log"
	"time"

	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/utils"
)

var priceData []float64
var counter int = 1
var lastSignal string
var orders []models.Order
var signals []models.Signal
var orderID int

func ProcessOrderBook(orderBook models.OrderBook) {
	if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
		log.Println("No bids or asks data received.")
		return
	}

	bidPrice := utils.StringToFloat64(orderBook.Bids[0][0].(string))
	askPrice := utils.StringToFloat64(orderBook.Asks[0][0].(string))
	midPrice := (bidPrice + askPrice) / 2

	log.Printf(
		"Counter: %d, Symbol: %s, EventTime: %d, Bid Price: %.2f, Ask Price: %.2f, Mid Price: %.2f\n",
		counter,
		orderBook.Symbol,
		orderBook.EventTime,
		bidPrice,
		askPrice,
		midPrice,
	)

	priceData = append(priceData, midPrice)

	if len(priceData) >= 200 {
		shortSMA := CalculateSMA(priceData, 50)
		longSMA := CalculateSMA(priceData, 200)

		log.Printf("Short SMA: %.2f, Long SMA: %.2f, Mid Price: %.2f\n", shortSMA, longSMA, midPrice)

		newSignal, reason := CheckSignal(shortSMA, longSMA, lastSignal)

		if newSignal != lastSignal {
			executeTrade(newSignal, midPrice, shortSMA, longSMA, reason)
		}
	}

	counter++
}

func executeTrade(newSignal string, midPrice, shortSMA, longSMA float64, reason string) {
	if len(orders) > 0 && orders[len(orders)-1].Status == "open" {
		orders[len(orders)-1].Status = "closed"
		log.Printf("Closing last order with ID: %d\n", orders[len(orders)-1].ID)
	}

	orderType := "sell"
	if newSignal == "BUY Signal!" {
		orderType = "buy"
	}

	orderID++

	orders = append(orders, models.Order{
		ID:        orderID,
		Price:     midPrice,
		Quantity:  1.0,
		Status:    "open",
		OrderType: orderType,
	})

	fmt.Println("order list:", orders)

	signals = append(signals, models.Signal{
		Timestamp: time.Now(),
		Type:      newSignal,
		Price:     midPrice,
		ShortSMA:  shortSMA,
		LongSMA:   longSMA,
		Reason:    reason,
	})

	fmt.Println("signal list:", signals)

	lastSignal = newSignal
	fmt.Println(newSignal)
}
