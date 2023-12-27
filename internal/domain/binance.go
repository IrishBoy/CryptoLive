package domain

type BinanceClient struct {
	BaseURL string
}
type CoinPriceResponse struct {
	Mins      int    `json:"mins"`
	Price     string `json:"price"`
	CloseTime int64  `json:"closeTime"`
}

func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		BaseURL: "https://api.binance.com/api/v3",
	}
}
