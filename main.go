package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/NikhilSharmaWe/crypto-exchange/client"
	"github.com/NikhilSharmaWe/crypto-exchange/server"
)

const (
	maxOrders = 3
)

var (
	tick = 1 * time.Second

	myAsks = make(map[float64]int64)
	myBids = make(map[float64]int64)
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		marketSellOrder := client.PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Size:   1000,
			Market: server.MarketETH,
		}

		orderResp, err := c.PlaceMarketOrder(&marketSellOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := client.PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Size:   1000,
			Market: server.MarketETH,
		}

		orderResp, err = c.PlaceMarketOrder(&marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}

		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}
		spread := math.Abs(bestBid.Price - bestAsk.Price)
		fmt.Println("exchange spread", spread)

		// place the bid
		if len(myBids) < 3 {
			bidLimit := client.PlaceOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid.Price + 100,
				Size:   1000,
				Market: server.MarketETH,
			}

			bidOrderResp, err := c.PlaceLimitOrder(&bidLimit)
			if err != nil {
				log.Println(bidOrderResp.OrderID)
			}

			myBids[bidLimit.Price] = bidOrderResp.OrderID
		}

		// place the ask
		if len(myAsks) < 3 {
			askLimit := client.PlaceOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk.Price - 100,
				Size:   1000,
				Market: server.MarketETH,
			}

			askOrderResp, err := c.PlaceLimitOrder(&askLimit)
			if err != nil {
				log.Println(askOrderResp.OrderID)
			}

			myAsks[askLimit.Price] = askOrderResp.OrderID
		}

		fmt.Println("best ask price =>", bestAsk.Price)
		fmt.Println("best bid price =>", bestBid.Price)

		<-ticker.C
	}
}

func seedMarket(c *client.Client) error {
	ask := client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  10_000,
		Size:   5_000_000,
		Market: server.MarketETH,
	}

	bid := client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   1_000_000,
		Market: server.MarketETH,
	}

	_, err := c.PlaceLimitOrder(&ask)
	if err != nil {
		return err
	}

	_, err = c.PlaceLimitOrder(&bid)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	go server.StartServer()
	time.Sleep(1 * time.Second)

	c := client.NewClient()

	if err := seedMarket(c); err != nil {
		panic(err)
	}

	go makeMarketSimple(c)
	time.Sleep(time.Second)

	marketOrderPlacer(c)

	// limitOrderParams := &client.PlaceOrderParams{
	// 	UserID: 8,
	// 	Bid:    false,
	// 	Price:  10_000,
	// 	Size:   5_000_000,
	// 	Market: server.MarketETH,
	// }

	// _, err := c.PlaceLimitOrder(limitOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// otherLimitOrderParams := &client.PlaceOrderParams{
	// 	UserID: 666,
	// 	Bid:    false,
	// 	Price:  9_000,
	// 	Size:   500_000,
	// 	Market: server.MarketETH,
	// }

	// _, err = c.PlaceLimitOrder(otherLimitOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// buyLimitOrderParams := &client.PlaceOrderParams{
	// 	UserID: 666,
	// 	Bid:    true,
	// 	Price:  11_000,
	// 	Size:   500_000,
	// 	Market: server.MarketETH,
	// }

	// _, err = c.PlaceLimitOrder(buyLimitOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Println("placed limit order from the client =>", resp.OrderID)

	// marketParams := &client.PlaceOrderParams{
	// 	UserID: 7,
	// 	Bid:    true,
	// 	Size:   1_000_000,
	// 	Market: server.MarketETH,
	// }

	// _, err = c.PlaceMarketOrder(marketParams)
	// if err != nil {
	// 	panic(err)
	// }

	// bestBid, err := c.GetBestBid()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("best bid price =>", bestBid.Price)

	// bestAsk, err := c.GetBestAsk()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("best ask price =>", bestAsk.Price)
	// // fmt.Println("placed market order from the client =>", resp.OrderID)

	select {}
}
