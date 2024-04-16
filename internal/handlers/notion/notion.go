package notion

import (
	"context"
	"fmt"
	"sync"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/handlers/common"
)

type NotionAPI interface {
	GetDatabases(ctx context.Context) ([]string, error)
	GetDatabase(ctx context.Context, tableID string) (domain.NotionTable, error)
	UpdateDatabaseBuyPage(ctx context.Context, pageID string, coinPrice, profit, profitValue float64) error
	Search(ctx context.Context) (domain.SearchResponse, error)
	UpdatePage(ctx context.Context, pageID string) error
	CreateDatabase(ctx context.Context, pageID string) error
}

type BinanceAPI interface {
	GetCoinPrice(ctx context.Context, coin string) (float64, error)
}

type CoinbaseAPI interface {
	GetCoinPrice(ctx context.Context, coin string) (float64, error)
}

type StorageAPI interface {
	SetCoinPriceBySource(ctx context.Context, source string, price float64, coin string) error
	GetAverageCoinPrice(ctx context.Context, token string) (float64, error)
}

type NotionTables struct {
	notionProvider   NotionAPI
	binanceProvider  BinanceAPI
	coinbaseProvider CoinbaseAPI
	storageProvider  StorageAPI
}

func (n *NotionTables) UpdateDatabases(ctx context.Context) error {
	ids, err := n.notionProvider.GetDatabases(ctx)
	if err != nil {
		return fmt.Errorf("error getting databases: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(ids))

	for _, databaseID := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			table, err := n.notionProvider.GetDatabase(ctx, id)
			if err != nil {
				errCh <- fmt.Errorf("error getting database %s: %v", id, err)
				return
			}

			for _, row := range table.Rows {
				if row.Coin != "USDT" {
					coinPrice, err := n.storageProvider.GetAverageCoinPrice(ctx, row.Coin)
					if err != nil {
						errCh <- fmt.Errorf("error getting average coin price for %s: %v", row.Coin, err)
						return
					}

					if coinPrice == 0 {
						price, err := n.binanceProvider.GetCoinPrice(ctx, row.Coin)
						if err != nil {
							errCh <- fmt.Errorf("error getting coin price from Binance for %s: %v", row.Coin, err)
							return
						}

						if err := n.storageProvider.SetCoinPriceBySource(ctx, "binance", price, row.Coin); err != nil {
							errCh <- fmt.Errorf("error setting price from Binance for %s: %v", row.Coin, err)
							return
						}

						price, err = n.coinbaseProvider.GetCoinPrice(ctx, row.Coin)
						if err != nil {
							errCh <- fmt.Errorf("error getting coin price from Coinbase for %s: %v", row.Coin, err)
							return
						}

						if err := n.storageProvider.SetCoinPriceBySource(ctx, "coinBase", price, row.Coin); err != nil {
							errCh <- fmt.Errorf("error setting price from Coinbase for %s: %v", row.Coin, err)
							return
						}

						coinPrice, err = n.storageProvider.GetAverageCoinPrice(ctx, row.Coin)
						if err != nil {
							errCh <- fmt.Errorf("error getting average coin price for %s: %v", row.Coin, err)
							return
						}
					}

					if err := n.UpdateDatabaseRow(ctx, row, coinPrice); err != nil {
						errCh <- fmt.Errorf("error updating database row: %v", err)
						return
					}
				}
			}
		}(databaseID)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errors)
	}

	return nil
}

func (n *NotionTables) UpdateDatabaseRow(ctx context.Context, row domain.NotionTableRow, currentCoinPrice float64) error {
	switch row.OpeartionType {
	case "Buy":
		updatedRow, err := UpdateCoinPrice(row, currentCoinPrice)
		if err != nil {
			return fmt.Errorf("error updating coin price: %v", err)
		}
		updatedRow, err = CoinGain(updatedRow)
		if err != nil {
			return fmt.Errorf("error updating gain: %v", err)
		}
		updatedRow, err = CoinPercentageGain(updatedRow)
		if err != nil {
			return fmt.Errorf("error updating percentage gain: %v", err)
		}
		fmt.Println(updatedRow)
		if err := n.notionProvider.UpdateDatabaseBuyPage(ctx, updatedRow.ID, updatedRow.CurrentCointPrice, updatedRow.Gain, updatedRow.PercentageGain); err != nil {
			return fmt.Errorf("error updating database buy page: %v", err)
		}
	default:
		// Handle other cases
	}
	return nil
}

func (n *NotionTables) CreateSpaces(ctx context.Context) error {
	searchNotion, err := n.notionProvider.Search(ctx)
	if err != nil {
		return fmt.Errorf("error searching Notion: %v", err)
	}

	parentPages := n.FilterParentPages(searchNotion)
	for _, pageID := range parentPages {
		if err := n.notionProvider.UpdatePage(ctx, pageID); err != nil {
			return fmt.Errorf("error updating page %s: %v", pageID, err)
		}
		if err := n.notionProvider.CreateDatabase(ctx, pageID); err != nil {
			return fmt.Errorf("error creating database for page %s: %v", pageID, err)
		}
	}

	return nil
}

func (n *NotionTables) FilterParentPages(allPages domain.SearchResponse) []string {
	parentPages := []string{}
	databasesParents := []string{}
	for _, page := range allPages.Results {
		if page.Parent.Type == "page_id" && page.Object != "database" {
			parentPages = append(parentPages, page.ID)
		}
		if page.Parent.Type == "workspace" {
			parentPages = append(parentPages, page.ID)
		}
		if page.Object == "database" {
			databasesParents = append(databasesParents, page.Parent.PageID)
		}
	}
	newPages := common.FilterArray(parentPages, databasesParents)
	return newPages
}

func (n *NotionTables) GetCoinsPrices(ctx context.Context, coins []string) (map[string]float64, error) {
	coinsPrices := make(map[string]float64)
	for _, coin := range coins {
		price, err := n.binanceProvider.GetCoinPrice(ctx, coin)
		if err != nil {
			fmt.Printf("Error getting coin price from Binance for %s: %v\n", coin, err)
			price, err = n.coinbaseProvider.GetCoinPrice(ctx, coin)
			if err != nil {
				fmt.Printf("Error getting coin price from Coinbase for %s: %v\n", coin, err)
				continue
			}
		}
		coinsPrices[coin] = price
	}
	return coinsPrices, nil
}

func UpdateCoinPrice(row domain.NotionTableRow, coinPrice float64) (domain.NotionTableRow, error) {
	row.CurrentCointPrice = coinPrice
	return row, nil
}

func CoinPercentageGain(row domain.NotionTableRow) (domain.NotionTableRow, error) {
	row.PercentageGain = float64(row.Gain / row.SoldAmount)
	return row, nil
}

func CoinGain(row domain.NotionTableRow) (domain.NotionTableRow, error) {
	if row.CurrentCointPrice >= 0 {
		row.Gain = float64(row.CurrentCointPrice*row.BoughtAmount - row.SoldAmount)
	}
	return row, nil
}

func New(notionProvider NotionAPI, binanceProvider BinanceAPI, coinbaseProvider CoinbaseAPI, storageProvider StorageAPI) *NotionTables {
	return &NotionTables{notionProvider: notionProvider, binanceProvider: binanceProvider, coinbaseProvider: coinbaseProvider, storageProvider: storageProvider}
}
