// notion/notion.go

package notion

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/common"
)

// Notion is a struct representing the Notion API client.
type Notion struct {
	NotionClient domain.NotionClient
}

func CreateURLDatabase(baseURL string, group string, databaseID string) string {
	return fmt.Sprintf("%s/databases/%s/query", baseURL, databaseID)
}

func CreateURLDatabases(baseURL string, group string) string {
	return fmt.Sprintf("%s/databases", baseURL)
}

func CreateURLPages(baseURL string, pageID string) string {
	return fmt.Sprintf("%s/pages/%s", baseURL, pageID)
}

// makeRequest is a common function for making HTTP requests.
// SHoul be moved to common package
func (n *Notion) makeRequest(method string, url string, payloadBytes []byte) (*http.Response, error) {
	client := common.CreateHTTPClient()

	req, err := common.CreateRequest(method, url, payloadBytes)
	if err != nil {
		return nil, err
	}

	headers := n.NotionClient.CreateRequestHeaders()
	req = common.AddHeaders(headers, req)

	return client.Do(req)
}

// GetDatabase retrieves data from a Notion database using the specified tableID.
func (n *Notion) GetDatabase(tableID string) {
	url := CreateURLDatabase(n.NotionClient.BaseURL, "databases", tableID)
	resp, err := n.makeRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}

// GetDatabasesList retrieves a list of databases from Notion.
func (n *Notion) GetDatabasesList() {
	url := CreateURLDatabases(n.NotionClient.BaseURL, "databases")
	resp, err := n.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}

func (n *Notion) UpdateDatabase(pageID string, operationID string, coinPrice float64, profitValue float64) {
	url := CreateURLPages(n.NotionClient.BaseURL, pageID)

	payload := n.NotionClient.UpdateTablePayload(coinPrice, profitValue, operationID)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Print(fmt.Errorf("error encoding JSON payload: %v", err))
	}
	resp, err := n.makeRequest(http.MethodGet, url, payloadBytes)
	if err != nil {
		fmt.Print(fmt.Errorf("error making request: %v", err))
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

}
