package domain

type CoinbaseClient struct {
	BaseURL string
}
type CoinbaseResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

func NewCoinbaseClient() *CoinbaseClient {
	return &CoinbaseClient{
		BaseURL: "https://api.coinbase.com/v2",
	}
}
