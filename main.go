package main

import (
	"fmt"
	"time"

	"github.com/NikhilSharmaWe/crypto-exchange/client"
	"github.com/NikhilSharmaWe/crypto-exchange/server"
)

func main() {
	// var mutex sync.Mutex

	go server.StartServer()
	time.Sleep(1 * time.Second)

	c := client.NewClient()

	bidParams := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  10_000,
		Size:   1000,
		Market: server.MarketETH,
	}

	go func() {
		for {
			// mutex.Lock()
			resp, err := c.PlaceLimitOrder(bidParams)
			if err != nil {
				panic(err)
			}
			fmt.Println("order id =>", resp.OrderID)

			if err := c.CancelOrder(resp.OrderID); err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	askParams := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  8_000,
		Size:   1000,
		Market: server.MarketETH,
	}

	for {
		// mutex.Lock()
		resp, err := c.PlaceLimitOrder(askParams)
		if err != nil {
			panic(err)
		}
		fmt.Println("order id =>", resp.OrderID)
		// mutex.Unlock()

		if err := c.CancelOrder(resp.OrderID); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	select {}
}
