package main

import (
	"log"
	"os"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	notionHandler "github.com/IrishBoy/CryptoLive/internal/handlers/notion"
	"github.com/IrishBoy/CryptoLive/internal/providers/binance"
	"github.com/IrishBoy/CryptoLive/internal/providers/notion"
	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	API_NOTION_KEY := goDotEnvVariable("NOTION_API_KEY")
	notionClient := domain.NewNotionClient(API_NOTION_KEY)
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
