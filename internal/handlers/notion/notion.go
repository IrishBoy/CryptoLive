package notion

import (
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/handlers/common"
)

type NotionAPI interface {
	GetDatabases() ([]string, error)
	GetDatabase(tableID string) (domain.NotionTable, error)
	UpdateDatabase(pageID string, coinPrice float64, profitValue float64) error
}

type BinanceAPI interface {
	GetCoinPrice(coin string) (float64, error)
}
type NotionTables struct {
	notionProvider  NotionAPI
	binanceProvider BinanceAPI
}

func (n *NotionTables) UpdateDatabases() {
	tables := []domain.NotionTable{}
	coins := make(map[string]bool)
	ids, err := n.notionProvider.GetDatabases()
	if err != nil {
		fmt.Println("Error getting databases")
	}

	for _, databaseID := range ids {
		table, err := n.notionProvider.GetDatabase(databaseID)
		if err != nil {
			fmt.Println("Error getting database")
		}
		tables = append(tables, table)
		for _, row := range table.Rows {
			common.AddString(coins, row.Coin)
		}

	}

	uniqueCoins := common.GetUniqueValues(coins)
	coinsPrices := n.GetCoinsPrices(uniqueCoins)

	for databaseID, database := range tables {
		for rowID, row := range database.Rows {
			updatedRow, err := UpdateCoinPrice(row, coinsPrices)
			if err != nil {
				fmt.Println("Error updating coin price")
			}
			updatedRow, err = CointGains(updatedRow)
			if err != nil {
				fmt.Println("Error updating gain")
			}
			tables[databaseID].Rows[rowID] = updatedRow
			err = n.notionProvider.UpdateDatabase(updatedRow.ID, updatedRow.CurrentCointPrice, updatedRow.PercentageGain)
			if err != nil {
				fmt.Println("Error updating database:", databaseID)
			}
		}

	}
}

func (n *NotionTables) GetCoinsPrices(coins []string) map[string]float64 {
	coinsPrices := make(map[string]float64)
	for _, coin := range coins {
		price, err := n.binanceProvider.GetCoinPrice(coin)
		if err != nil {
			fmt.Println("Error getting coin price")
		}
		coinsPrices[coin] = price
	}
	return coinsPrices
}

// write verification on no coin price
func UpdateCoinPrice(row domain.NotionTableRow, coins map[string]float64) (domain.NotionTableRow, error) {
	row.CurrentCointPrice = float64(coins[row.Coin])
	return row, nil
}

// write verification on row.SoldAmoint == 0
func CointGains(row domain.NotionTableRow) (domain.NotionTableRow, error) {
	row.PercentageGain = float64((row.CurrentCointPrice*row.CoinAmount - row.SoldAmount) / row.SoldAmount)
	return row, nil

}

func New(notionProvider NotionAPI, binanceProvider BinanceAPI) *NotionTables {
	return &NotionTables{notionProvider: notionProvider, binanceProvider: binanceProvider}
}
