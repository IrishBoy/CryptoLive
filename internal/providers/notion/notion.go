package notion

import (
	"fmt"

	"github.com/IrishBoy/CryptoLive/internal/domain"
	common "github.com/IrishBoy/CryptoLive/internal/providers"
)

type Notion struct {
	NotionClient domain.NotionClient
}

func (c *Notion) GetDatabase(tableID string) {
	URL := common.FormatURL(c.NotionClient.BaseURL, "databases", tableID)
	client := common.CreateHTTPClient()

	req, _ := common.CreateRequest("POST", URL)

	headers := c.NotionClient.CreateRequestHeaders()

	req = common.AddHeaders(headers, req)
	fmt.Println(URL)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	fmt.Println(resp)
	defer resp.Body.Close()
}
