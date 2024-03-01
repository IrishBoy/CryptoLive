package domain

type NotionTableRow struct {
	ID                string
	Coin              string
	CurrentCointPrice float64
	BoughtAmount      float64
	Gain              float64
	PercentageGain    float64
	SoldCoin          string
	SoldAmount        float64
}

type NotionTable struct {
	DatabaseID string
	Rows       []NotionTableRow
}
