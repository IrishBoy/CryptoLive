package notion

import (
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/handlers/common"
)

type NotionAPI interface {
	GetDatabases() ([]string, error)
	GetDatabase(tableID string) (domain.NotionTable, error)
	UpdateDatabase(pageID string, coinPrice float64, profit float64, profitValue float64) error
	Search() (domain.SearchResponse, error)
	UpdatePage(pageID string) (err error)
	CreateDatabase(pageID string) (err error)
}

type BinanceAPI interface {
	GetCoinPrice(coin string) (float64, error)
}
type CoinbaseAPI interface {
	GetCoinPrice(coin string) (float64, error)
}

type NotionTables struct {
	notionProvider   NotionAPI
	binanceProvider  BinanceAPI
	coinbaseProvider CoinbaseAPI
}

func (n *NotionTables) UpdateDatabases() {
	tables := []domain.NotionTable{}
	coins := make(map[string]bool)
	ids, err := n.notionProvider.GetDatabases()
	if err != nil {
		fmt.Println("Error getting databases")
	}
	// Rewrtite so it will be done in parallel
	for _, databaseID := range ids {
		table, err := n.notionProvider.GetDatabase(databaseID)
		if err != nil {
			fmt.Println("Error getting database:")
		}
		tables = append(tables, table)
		for _, row := range table.Rows {
			if row.Coin == "USDT" {
				common.AddString(coins, row.SoldCoin)
			} else {
				common.AddString(coins, row.Coin)
			}

		}

	}

	uniqueCoins := common.GetUniqueValues(coins)
	coinsPrices := n.GetCoinsPrices(uniqueCoins)

	for databaseID, database := range tables {
		// Rewrtite so it will be done in parallel

		for rowID, row := range database.Rows {
			updatedRow, err := UpdateCoinPrice(row, coinsPrices)
			if err != nil {
				fmt.Println("Error updating coin price")
			}
			updatedRow, err = CoinGain(updatedRow)
			if err != nil {
				fmt.Println("Error updating gain")
			}
			updatedRow, err = CoinPercentageGain(updatedRow)
			if err != nil {
				fmt.Println("Error updating percetage gain")
			}
			tables[databaseID].Rows[rowID] = updatedRow
			err = n.notionProvider.UpdateDatabase(updatedRow.ID, updatedRow.CurrentCointPrice, updatedRow.Gain, updatedRow.PercentageGain)
			if err != nil {
				fmt.Println("Error updating database:", databaseID)
			}
		}

	}
}

func (n *NotionTables) CreateSpaces() {
	// Get pages without child databases
	// Create blocks with parameters
	// Create databases
	searchNotion, _ := n.notionProvider.Search()
	parentPages := n.FilterParentPages(searchNotion)
	for _, pageId := range parentPages {
		err := n.notionProvider.UpdatePage(pageId)
		if err != nil {
			fmt.Println("Error updating page:", err)
		}
		err = n.notionProvider.CreateDatabase(pageId)
		if err != nil {
			fmt.Println("Error creaing database:", err)
		}
	}

}

// Return only pages that
// 1) Are not pages of a database already  -> Check parent type
// 2) don't have databse as a child -> ?????
func (n *NotionTables) FilterParentPages(allPages domain.SearchResponse) []string {
	parentPages := []string{}
	databases_parents := []string{}
	for _, page := range allPages.Results {
		// fmt.Printf("%+v\n", page.Parent)
		// fmt.Printf("%#v\n", page)
		if page.Parent.Type == "page_id" && page.Object != "database" {
			//add only pages that are not ones of a database already
			parentPages = append(parentPages, page.ID)
		}
		if page.Parent.Type == "workspace" {
			parentPages = append(parentPages, page.ID)
		}
		if page.Object == "database" {
			databases_parents = append(databases_parents, page.Parent.PageID)
		}
		// TODO: Check all the tree
	}
	new_pages := common.FilterArray(parentPages, databases_parents)

	return new_pages

}

func (n *NotionTables) GetCoinsPrices(coins []string) map[string]float64 {
	coinsPrices := make(map[string]float64)
	// Rewrtite so it will be done in parallel
	for _, coin := range coins {
		price, err := n.binanceProvider.GetCoinPrice(coin)
		if err != nil {
			fmt.Printf("Error getting coin price from binance: %s\n", coin)
			price, err = n.coinbaseProvider.GetCoinPrice(coin)
			if err != nil {
				fmt.Printf("Error getting coin price from coinbase: %s\n", coin)
			}
		}

		coinsPrices[coin] = price
	}
	return coinsPrices
}

// write verification on no coin price
func UpdateCoinPrice(row domain.NotionTableRow, coins map[string]float64) (domain.NotionTableRow, error) {
	if row.Coin == "USDT" {
		row.CurrentCointPrice = 1 / float64(coins[row.SoldCoin])
	} else {
		row.CurrentCointPrice = float64(coins[row.Coin])
	}

	return row, nil
}

// write verification on row.SoldAmoint == 0
func CoinPercentageGain(row domain.NotionTableRow) (domain.NotionTableRow, error) {
	if row.Coin == "USDT" {
		row.PercentageGain = float64(row.Gain / row.BoughtAmount)
	} else {
		row.PercentageGain = float64(row.Gain / row.SoldAmount)
	}
	return row, nil

}

func CoinGain(row domain.NotionTableRow) (domain.NotionTableRow, error) {
	if row.Coin == "USDT" {
		row.Gain = -float64((1/row.CurrentCointPrice)*row.SoldAmount - row.BoughtAmount)
	} else {
		row.Gain = float64(row.CurrentCointPrice*row.BoughtAmount - row.SoldAmount)
	}
	return row, nil
}

func New(notionProvider NotionAPI, binanceProvider BinanceAPI, coinbaseProvider CoinbaseAPI) *NotionTables {
	return &NotionTables{notionProvider: notionProvider, binanceProvider: binanceProvider, coinbaseProvider: coinbaseProvider}
}
