// cmd/main.go
package main

import (
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/binance"
)

func main() {
	// Create an instance of the Notion type
	// notionInstance := &notion.Notion{
	// 	NotionClient: domain.NotionClient{
	// 		APIKey: "",
	// 	}, // Initialize with an appropriate implementation
	// }

	binanceClient := domain.NewBinanceClient()

	// Create a Binance instance and assign the BinanceClient instance
	binanceInstance := &binance.Binance{
		BinanceClient: *binanceClient,
	}
	price, err := binanceInstance.GetCoinPrice("BTC")
	fmt.Print(price, err)
	// notionInstance.GetDatabase("")
	// notionInstance.UpdateDatabase("123123", "321123", 123.33, 123.33)
}
