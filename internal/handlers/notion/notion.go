package notion

import (
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/handlers/common"
)

type NotionAPI interface {
	GetDatabases() ([]string, error)
	GetDatabase(tableID string) (domain.NotionTable, error)
	UpdateDatabase(pageID string, operationID string, coinPrice float64, profitValue float64) error
}

type BinanceAPI interface {
	GetCoinPrice(coin string) (float64, error)
}
type NotionTables struct {
	notionProvider  NotionAPI
	binanceProvider BinanceAPI
}

func (n *NotionTables) UpdateDatabase() {
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
		for _, row := range table.Rows {
			common.AddString(coins, row.Coin)
		}

	}

	uniqueCoins := common.GetUniqueValues(coins)
	for _, coin := range uniqueCoins {
		fmt.Print(n.binanceProvider.GetCoinPrice(coin))
	}
}

func New(notionProvider NotionAPI, binanceProvider BinanceAPI) *NotionTables {
	return &NotionTables{notionProvider: notionProvider, binanceProvider: binanceProvider}
}
