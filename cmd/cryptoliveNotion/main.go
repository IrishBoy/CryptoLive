package main

import (
	"github.com/IrishBoy/CryptoLive/internal/domain"
	notionHandler "github.com/IrishBoy/CryptoLive/internal/handlers/notion"
	"github.com/IrishBoy/CryptoLive/internal/providers/binance"
	"github.com/IrishBoy/CryptoLive/internal/providers/notion"
)

func main() {
	notionClient := domain.NewNotionClient("")
	notionInstance := &notion.Notion{
		NotionClient: *notionClient,
	}

	binanceClient := domain.NewBinanceClient()

	// Create a Binance instance and assign the BinanceClient instance
	binanceInstance := &binance.Binance{
		BinanceClient: *binanceClient,
	}

	// Create an instance of NotionTables
	notionTables := notionHandler.New(notionInstance, binanceInstance)

	// Use the instance to update databases
	notionTables.UpdateDatabases()
}
