package services

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/turgaysozen/algotrading/config"
	"github.com/turgaysozen/algotrading/db"
	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/utils"
)

var priceDataMap sync.Map
var lastSignalMap sync.Map

func ProcessOrderBook(orderBook models.OrderBook) {
	if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
		log.Println("No bids or asks data received.")
		return
	}

	bidPrice := GetBestBidPrice(orderBook.Bids)
	askPrice := GetBestAskPrice(orderBook.Asks)
	midPrice := (bidPrice + askPrice) / 2

	orderBookID, err := db.SaveOrderBook(orderBook.EventType, orderBook.Symbol, orderBook.EventTime, bidPrice, askPrice)
	if err != nil {
		log.Printf("Error saving order book: %v", err)
		return
	}

	log.Printf("ID: %d | Symbol: %s | EventTime: %d | Bid: %.2f | Ask: %.2f | Mid Price: %.2f\n",
		orderBookID, orderBook.Symbol, orderBook.EventTime, bidPrice, askPrice, midPrice)

	value, _ := priceDataMap.LoadOrStore(orderBook.Symbol, &[]float64{})
	priceData := value.(*[]float64)

	appendPriceData(priceData, midPrice)

	if len(*priceData) >= config.MaxPriceCount {
		shortSMA := CalculateSMA(*priceData, config.ShortSMACount)
		longSMA := CalculateSMA(*priceData, config.LongSMACount)

		value, _ := lastSignalMap.LoadOrStore(orderBook.Symbol, "")
		lastSignal := value.(string)

		newSignal, reason := CheckSignal(shortSMA, longSMA, lastSignal)

		if newSignal != lastSignal {
			lastSignalMap.Store(orderBook.Symbol, newSignal)
			saveSignal(newSignal, midPrice, shortSMA, longSMA, reason, orderBook.Symbol)
		}
	}
}

func saveSignal(newSignal string, midPrice, shortSMA, longSMA float64, reason, symbol string) {
	signal := models.Signal{
		Type:     newSignal,
		Price:    midPrice,
		ShortSMA: shortSMA,
		LongSMA:  longSMA,
		Reason:   reason,
	}

	err := db.SaveSignal(signal)
	if err != nil {
		log.Printf("Error saving signal: %v", err)
		return
	}

	signalJSON, _ := json.MarshalIndent(signal, "", "  ")
	log.Println("Signal saved successfully:", string(signalJSON))

	saveOrder(newSignal, midPrice, symbol)
}

func saveOrder(newSignal string, midPrice float64, symbol string) {
	lastOrder, err := db.GetLastOpenOrder()
	if err != nil {
		log.Printf("Error retrieving last open order: %v", err)
		return
	}

	if lastOrder != nil {
		err := db.CloseOrder(lastOrder.ID)
		if err != nil {
			log.Printf("Error closing last open order: %v", err)
			return
		}
		log.Printf("Closing last order with ID: %d\n", lastOrder.ID)
	}

	orderType := "sell"
	if newSignal == "BUY Signal!" {
		orderType = "buy"
	}

	order := models.Order{
		Price:     midPrice,
		Quantity:  1.0,
		Status:    "open",
		OrderType: orderType,
	}

	err = db.SaveOrder(order)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		return
	}

	log.Printf("Order saved successfully: Type= %s, Price= %.2f, Symbol= %s, Timestamp= %s",
		orderType, midPrice, symbol, time.Now())
}

func GetBestBidPrice(bids [][]interface{}) float64 {
	if len(bids) == 0 {
		return 0
	}

	bestBid := bids[0]
	for _, bid := range bids {
		bidPrice := utils.StringToFloat64(bid[0].(string))
		if bidPrice > utils.StringToFloat64(bestBid[0].(string)) {
			bestBid = bid
		}
	}
	return utils.StringToFloat64(bestBid[0].(string))
}

func GetBestAskPrice(asks [][]interface{}) float64 {
	if len(asks) == 0 {
		return 0
	}

	bestAsk := asks[0]
	for _, ask := range asks {
		askPrice := utils.StringToFloat64(ask[0].(string))
		if askPrice < utils.StringToFloat64(bestAsk[0].(string)) {
			bestAsk = ask
		}
	}
	return utils.StringToFloat64(bestAsk[0].(string))
}

func appendPriceData(priceData *[]float64, midPrice float64) {
	*priceData = append(*priceData, midPrice)

	// Keep last 200 records
	if len(*priceData) > config.MaxPriceCount {
		*priceData = (*priceData)[1:]
	}
}
