// notion/notion.go

package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/common"
)

// Notion is a struct representing the Notion API client.
type Notion struct {
	NotionClient domain.NotionClient
}

func CreateURLSearch(baseURL string) string {
	return fmt.Sprintf("%s/search", baseURL)
}
func CreateURLDatabase(baseURL string, databaseID string) string {
	return fmt.Sprintf("%s/databases/%s/query", baseURL, databaseID)
}

func CreateURLDatabases(baseURL string) string {
	return fmt.Sprintf("%s/databases", baseURL)
}

func CreateURLPages(baseURL string, pageID string) string {
	return fmt.Sprintf("%s/pages/%s", baseURL, pageID)
}

func CreateURLAppendBlock(baseURL string, pageID string) string {
	return fmt.Sprintf("%s/blocks/%s/children", baseURL, pageID)
}

func CreateURLCreateDatabase(baseURL string) string {
	return fmt.Sprintf("%s/databases", baseURL)
}

// makeRequest is a common function for making HTTP requests.
// SHoul be moved to common package
func (n *Notion) makeRequest(method string, url string, payloadBytes []byte, headersType string) (*http.Response, error) {
	client := common.CreateHTTPClient()

	req, err := common.CreateRequest(method, url, payloadBytes)
	if err != nil {
		return nil, err
	}

	headers := n.NotionClient.CreateRequestHeaders(headersType)
	req = common.AddHeaders(headers, req)
	return client.Do(req)
}

// GetDatabase retrieves data from a Notion database using the specified tableID.
func (n *Notion) GetDatabase(tableID string) (domain.NotionTable, error) {
	url := CreateURLDatabase(n.NotionClient.BaseURL, tableID)

	response, err := n.makeRequest(http.MethodPost, url, nil, "new")
	if err != nil {
		fmt.Println("Error making request:", err)
		return domain.NotionTable{}, err
	}
	defer response.Body.Close()
	var resp map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		fmt.Println("Error decoding response:", err)
		return domain.NotionTable{}, err
	}

	var notionTable domain.NotionTable // Declare notionTable outside the loop
	var rows []domain.NotionTableRow
	// Rewrtite so it will be done in parallel
	for _, v := range resp["results"].([]interface{}) {
		id, ok := v.(map[string]interface{})["id"].(string)
		if !ok {
			fmt.Println("Error: 'id' not found or not a map")
			continue
		}
		properties, ok := v.(map[string]interface{})["properties"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: 'properties' not found or not a map")
			continue
		}

		coinSelect, ok := properties["Coin bought"].(map[string]interface{})["select"]
		if !ok || coinSelect == nil {
			fmt.Println("Error: 'Coin bought' or 'select' is nil or not found")
			continue
		}

		coin, ok := coinSelect.(map[string]interface{})["name"].(string)
		if !ok {
			fmt.Println("Error: 'name' not found or not a string")
			continue
		}

		soldCoinSelect, ok := properties["Coin Sold"].(map[string]interface{})["select"]
		if !ok || soldCoinSelect == nil {
			fmt.Println("Error: 'Coin bought' or 'select' is nil or not found")
			continue
		}

		coinSold, ok := soldCoinSelect.(map[string]interface{})["name"].(string)
		if !ok {
			fmt.Println("Error: 'name' not found or not a string")
			continue
		}

		row := domain.NotionTableRow{
			ID:                id,
			Coin:              coin,
			CurrentCointPrice: 0,
			BoughtAmount:      properties["Bought Amount"].(map[string]interface{})["number"].(float64),
			Gain:              0,
			PercentageGain:    0,
			SoldCoin:          coinSold,
			SoldAmount:        properties["Sold Amount"].(map[string]interface{})["number"].(float64),
		}

		rows = append(rows, row)

	}
	notionTable = domain.NotionTable{
		DatabaseID: tableID,
		Rows:       rows,
	}

	return notionTable, nil
}

// GetDatabases retrieves a list of databases from Notion.
func (n *Notion) GetDatabases() ([]string, error) {

	url := CreateURLDatabases(n.NotionClient.BaseURL)
	response, err := n.makeRequest(http.MethodGet, url, nil, "old")

	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer response.Body.Close()

	var resp map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		fmt.Println("Error decoding response:", err)
	}

	var ids []string
	for _, v := range resp["results"].([]interface{}) {
		ids = append(ids, v.(map[string]interface{})["id"].(string))
	}

	return ids, nil
}

func (n *Notion) UpdateDatabase(pageID string, coinPrice float64, profit float64, profitValue float64) error {
	url := CreateURLPages(n.NotionClient.BaseURL, pageID)

	payload := n.NotionClient.UpdateTablePayload(coinPrice, profit, profitValue)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := n.makeRequest(http.MethodPatch, url, payloadBytes, "new")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (n *Notion) Search() (domain.SearchResponse, error) {
	url := CreateURLSearch(n.NotionClient.BaseURL)
	response, err := n.makeRequest(http.MethodPost, url, nil, "new")

	if err != nil {
		return domain.SearchResponse{}, err
	}
	defer response.Body.Close()
	var result domain.SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return domain.SearchResponse{}, err
	}

	return result, nil
}

// We may want to set out own column names, so we need to give users functionality to set it in the parent page
// So, we need to retrieve those names
// If page is connected to our extention and there is
// 1) Some pattern that we will configure
// 2) Database as a child page
// -> We need to get this parameters
func (n *Notion) GetColumns(pageID string) ([]string, error) {
	return []string{}, nil
}

// If a user gives us an access to some page we can
// Create a pattern for the page so he can cofigure column names
// Create databse as a child page -> So a user will not need to do this
func (n *Notion) UpdatePage(pageID string) (err error) {
	url := CreateURLAppendBlock(n.NotionClient.BaseURL, pageID)
	pageBlock, err := os.ReadFile("internal/domain/page_template.json")

	if err != nil {
		fmt.Println("Error reading page block file:", err)
		return err
	}

	_, err = n.makeRequest(http.MethodPatch, url, pageBlock, "new")

	if err != nil {
		fmt.Println("Error making request to append blocks to page:", err)

	}

	return nil
}

func (n *Notion) CreateDatabase(pageID string) (err error) {
	url := CreateURLCreateDatabase(n.NotionClient.BaseURL)
	databaseBlock := n.NotionClient.CreateDatabasePayload(pageID)
	databaseBytes, err := json.Marshal(databaseBlock)

	if err != nil {
		fmt.Println("Error marshaling database block to json:", err)
		return err
	}

	_, err = n.makeRequest(http.MethodPost, url, databaseBytes, "new")
	if err != nil {
		fmt.Println("Error making request to create database:", err)
		return err

	}

	return nil
}
