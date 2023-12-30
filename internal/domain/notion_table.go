package domain

type NotionTableRow struct {
	ID                string
	Coin              string
	OldCoinPrice      float64
	CurrentCointPrice float64
	CoinAmount        float64
	Gain              float64
	PercentageGain    float64
	SoldCoin          string
	SoldAmount        float64
}

type NotionTable struct {
	DatabaseID string
	Rows       []NotionTableRow
}
