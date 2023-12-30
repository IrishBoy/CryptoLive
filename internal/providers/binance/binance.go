package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/common"
)

type Binance struct {
	BinanceClient domain.BinanceClient
}

func CreateURLCoinPrice(baseURL string, currency string) string {
	return fmt.Sprintf("%s/avgPrice?symbol=%sUSDT", baseURL, currency)
}

func (b *Binance) makeRequest(method string, url string, payloadBytes []byte) (*http.Response, error) {
	client := common.CreateHTTPClient()

	req, err := common.CreateRequest(method, url, payloadBytes)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (b *Binance) GetCoinPrice(coin string) (float64, error) {
	url := CreateURLCoinPrice(b.BinanceClient.BaseURL, coin)
	resp, err := b.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return math.NaN(), fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return math.NaN(), fmt.Errorf("non-OK status code - %s", resp.Status)
	}

	var result domain.CoinPriceResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return math.NaN(), fmt.Errorf("cannot unmarshal JSON")
	}
	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return math.NaN(), fmt.Errorf("cannot get price")
	}

	return price, nil
}
