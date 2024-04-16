package main

import (
	"context"
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	commonHandler "github.com/IrishBoy/CryptoLive/internal/handlers/common"
	notionHandler "github.com/IrishBoy/CryptoLive/internal/handlers/notion"
	"github.com/IrishBoy/CryptoLive/internal/providers/binance"
	"github.com/IrishBoy/CryptoLive/internal/providers/coinbase"
	"github.com/IrishBoy/CryptoLive/internal/providers/notion"
	"github.com/IrishBoy/CryptoLive/internal/providers/storage"
)

func main() {
	// Retrieve environment variables
	API_NOTION_KEY := commonHandler.GoDotEnvVariable("NOTION_API_KEY")
	REDIS_USERNAME := commonHandler.GoDotEnvVariable("REDIS_USERNAME")
	REDIS_PASSWORD := commonHandler.GoDotEnvVariable("REDIS_PASSWORD")
	REDIS_ADDRESS := commonHandler.GoDotEnvVariable("REDIS_ADDRESS")

	// Initialize Notion client
	notionClient := domain.NewNotionClient(API_NOTION_KEY)
	notionInstance := &notion.Notion{
		NotionClient: *notionClient,
	}

	// Initialize Redis client
	redisClient, err := domain.InitializeRedis(REDIS_ADDRESS, REDIS_PASSWORD, REDIS_USERNAME, 0)
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}

	// Initialize storage provider with Redis client
	redisInstance := storage.NewRedisStorage(redisClient)

	// Initialize Binance and Coinbase clients
	binanceClient := domain.NewBinanceClient()
	coinbaseClient := domain.NewCoinbaseClient()

	// Create instances of Binance and Coinbase providers
	binanceInstance := &binance.Binance{
		BinanceClient: *binanceClient,
	}

	coinbaseInstance := &coinbase.Coinbase{
		CoinbaseClient: *coinbaseClient,
	}

	// Create an instance of NotionTables
	notionTables := notionHandler.New(notionInstance, binanceInstance, coinbaseInstance, redisInstance)
	ctx := context.Background()
	// Use the instance to update databases
	err = notionTables.UpdateDatabases(ctx)
	fmt.Println(err)
}
