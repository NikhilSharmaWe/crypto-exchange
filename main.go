package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/NikhilSharmaWe/crypto-exchange/client"
	"github.com/NikhilSharmaWe/crypto-exchange/server"
)

var tick = 1 * time.Second

func makeMarketSimple(c *client.Client) {
	for {
		ticker := time.NewTicker(tick)
		<-ticker.C

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

		// here we are tighting the spreads

		limitOrderBid := client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Price:  bestBid.Price + 100,
			Size:   1000,
			Market: server.MarketETH,
		}

		limitOrderBidResp, err := c.PlaceLimitOrder(&limitOrderBid)
		if err != nil {
			log.Println(limitOrderBidResp.OrderID)
		}

		limitOrderAsk := client.PlaceOrderParams{
			UserID: 7,
			Bid:    false,
			Price:  bestBid.Price - 100,
			Size:   1000,
			Market: server.MarketETH,
		}

		limitOrderAskResp, err := c.PlaceLimitOrder(&limitOrderAsk)
		if err != nil {
			log.Println(limitOrderAskResp.OrderID)
		}

		fmt.Println("best ask price =>", bestAsk.Price)
		fmt.Println("best bid price =>", bestBid.Price)
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

	makeMarketSimple(c)

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
