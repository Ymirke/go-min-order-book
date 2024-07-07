package main

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sort"
	"time"
)

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)

type Order struct {
	ID        string
	Type      OrderType
	Price     float64
	Quantity  float64
	Timestamp time.Time
}

type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

func NewOrder(id string, orderType OrderType, price float64, quantity float64) Order {
	return Order{
		ID:        id,
		Type:      orderType,
		Price:     price,
		Quantity:  quantity,
		Timestamp: time.Now(),
	}
}

func (ob *OrderBook) AddOrder(order Order) {
	fmt.Printf("Added order: %s %f @ %f \n", order.Type, order.Quantity, order.Price)

	if order.Type == Buy {
		ob.BuyOrders = append(ob.BuyOrders, order)
		sort.Slice(ob.BuyOrders, func(i, j int) bool {
			if ob.BuyOrders[i].Price == ob.BuyOrders[j].Price {
				return ob.BuyOrders[i].Timestamp.Before(ob.BuyOrders[j].Timestamp)
			}
			return ob.BuyOrders[i].Price > ob.BuyOrders[j].Price
		})
	} else {
		ob.SellOrders = append(ob.SellOrders, order)
		sort.Slice(ob.SellOrders, func(i, j int) bool {
			if ob.SellOrders[i].Price == ob.SellOrders[j].Price {
				return ob.SellOrders[i].Timestamp.Before(ob.SellOrders[j].Timestamp)
			}
			return ob.SellOrders[i].Price < ob.SellOrders[j].Price
		})
	}
}

func (ob *OrderBook) MatchOrders() {
	for len(ob.BuyOrders) > 0 && len(ob.SellOrders) > 0 {
		buyOrder := ob.BuyOrders[0]
		sellOrder := ob.SellOrders[0]

		if buyOrder.Price >= sellOrder.Price {
			tradedQuantity := min(buyOrder.Quantity, sellOrder.Quantity)
			fmt.Printf("Matched: Buy %f @ %f with Sell %f @ %f\n", tradedQuantity, buyOrder.Price, tradedQuantity, sellOrder.Price)

			buyOrder.Quantity -= tradedQuantity
			sellOrder.Quantity -= tradedQuantity

			if buyOrder.Quantity == 0 {
				ob.BuyOrders = ob.BuyOrders[1:]
			} else {
				ob.BuyOrders[0] = buyOrder
			}

			if sellOrder.Quantity == 0 {
				ob.SellOrders = ob.SellOrders[1:]
			} else {
				ob.SellOrders[0] = sellOrder
			}
		} else {
			break
		}
	}
}

func min(a float64, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func randomFloat64(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func randomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

func randomOrderType() OrderType {
	orderTypes := []OrderType{Buy, Sell}
	rand.Seed(time.Now().UnixNano())
	return orderTypes[rand.Intn(len(orderTypes))]
}

func randomOrder() Order {
	id := uuid.New()
	orderType := randomOrderType()
	price := randomFloat64(27, 35)
	amount := randomFloat64(0.1, 25)

	order := NewOrder(id.String(), orderType, price, amount)
	return order
}

func main() {
	orderBook := &OrderBook{}

	// infinite loop
	for true {
		numOrders := randomInt(0, 10)

		for i := 0; i < numOrders; i++ {
			orderBook.AddOrder(randomOrder())
		}

		orderBook.MatchOrders()

		// Log buy and sell orders
		fmt.Println("Remaining Buy Orders:")
		for _, order := range orderBook.BuyOrders {
			fmt.Printf("%+v\n", order)
		}
		fmt.Println("Remaining Sell Orders:")
		for _, order := range orderBook.SellOrders {
			fmt.Printf("%+v\n", order)
		}
		fmt.Println("\n")

		tenSeconds := 10000
		time.Sleep(time.Duration(tenSeconds) * time.Millisecond)
	}
}
