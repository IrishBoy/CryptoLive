package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	RedisClient *redis.Client
}

// NewRedisStorage creates a new Storage instance with Redis client
func NewRedisStorage(client *redis.Client) *Storage {
	return &Storage{
		RedisClient: client,
	}
}

// SetCoinPriceBySource sets the coin price in Redis storage by source
func (s *Storage) SetCoinPriceBySource(ctx context.Context, source string, price float64, coin string) error {
	key := coin + "_" + source
	err := s.RedisClient.Set(ctx, key, price, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set coin price by source: %v", err)
	}
	return nil
}

// GetAverageCoinPrice calculates the average coin price from Redis storage
func (s *Storage) GetAverageCoinPrice(ctx context.Context, token string) (float64, error) {
	var coinPrices []float64
	// fmt.Printf("coin value: %s\n", token)
	pattern := token + "_*"
	iter := s.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		// fmt.Printf("Iterator value: %s\n", iter.Val())
		price, err := s.GetCoinPriceByKey(ctx, iter.Val())
		if err != nil {
			fmt.Printf("Error getting coin price for key %s: %v\n", iter.Val(), err)
			continue
		}
		coinPrices = append(coinPrices, price)
	}
	if len(coinPrices) == 0 {
		return 0, nil // No prices found
	}
	avgPrice := calculateAverage(coinPrices)
	return avgPrice, nil
}

// GetCoinPriceByKey retrieves the coin price from Redis storage by key
func (s *Storage) GetCoinPriceByKey(ctx context.Context, key string) (float64, error) {
	priceString, err := s.RedisClient.Get(ctx, key).Float64()
	if err != nil {
		return -1, fmt.Errorf("failed to get coin price by key: %v", err)
	}
	return priceString, nil
}

// GetCoinPriceBySource retrieves the coin price from Redis storage by source
func (s *Storage) GetCoinPriceBySource(ctx context.Context, coin, source string) (float64, error) {
	key := coin + "_" + source
	price, err := s.GetCoinPriceByKey(ctx, key)
	if err != nil {
		return -1, fmt.Errorf("failed to get coin price by source: %v", err)
	}
	return price, nil
}

// calculateAverage calculates the average of a slice of float64 numbers
func calculateAverage(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
