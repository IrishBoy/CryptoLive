package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/common"
)

type Coinbase struct {
	CoinbaseClient domain.CoinbaseClient
}

func CreateURLCoinPrice(baseURL string, currency string) string {
	return fmt.Sprintf("%s/prices/%s-USDT/buy", baseURL, currency)
}

func (b *Coinbase) makeRequest(method string, url string, payloadBytes []byte) (*http.Response, error) {
	client := common.CreateHTTPClient()

	req, err := common.CreateRequest(method, url, payloadBytes)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (b *Coinbase) GetCoinPrice(ctx context.Context, coin string) (float64, error) {
	url := CreateURLCoinPrice(b.CoinbaseClient.BaseURL, coin)
	resp, err := b.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return -1, fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("non-OK status code - %s", resp.Status)
	}

	var result domain.CoinbaseResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return -1, fmt.Errorf("cannot unmarshal JSON")
	}
	amountStr := result.Data.Amount
	price, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return -1, fmt.Errorf("cannot get price")
	}

	return price, nil
}
