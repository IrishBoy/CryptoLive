package domain

type GetPagesResponse struct {
	Object         string `json:"object"`
	Results        []any  `json:"results"`
	NextCursor     any    `json:"next_cursor"`
	HasMore        bool   `json:"has_more"`
	Type           string `json:"type"`
	PageOrDatabase struct {
	} `json:"page_or_database"`
	RequestID string `json:"request_id"`
}

type NotionClient struct {
	BaseURL string
	APIKey  string
}

func NewNotionClient(apiKey string) *NotionClient {
	return &NotionClient{
		BaseURL: "https://api.notion.com/v1",
		APIKey:  apiKey,
	}
}

func (nc *NotionClient) CreateRequestHeaders(headersType string) map[string]string {

	var version string
	switch headersType {
	case "old":
		version = "2021-05-13"
	default:
		version = "2022-06-28"
	}

	return map[string]string{
		"Notion-Version": version,
		"Authorization":  "Bearer " + nc.APIKey,
		"Content-Type":   "application/json",
		"Accept":         "application/json",
	}

}

func (nc *NotionClient) UpdateTablePayload(coinPrice, profitValue float64) map[string]interface{} {
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"Current Coin Price": map[string]interface{}{
				"type":   "number",
				"number": coinPrice,
			},
			"ProfitValue": map[string]interface{}{
				"type":   "number",
				"number": profitValue,
			},
		},
	}
}

func (nc *NotionClient) UpdatePagePayload(pageID string) map[string]interface{} {
	return nil
}

func (nc *NotionClient) CreateDatabasePayload(pageID string) map[string]interface{} {
	return map[string]interface{}{
		"parent": map[string]interface{}{
			"type":    "page_id",
			"page_id": pageID,
		},
		"title": []map[string]interface{}{
			{
				"type": "text",
				"text": map[string]interface{}{
					"content": "Operations",
					"link":    nil,
				},
			},
		},
		"properties": map[string]interface{}{
			"Coin bought": map[string]interface{}{
				"select": map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"name": "BTC",
						},
						{
							"name": "ETH",
						},
					},
				},
			},
			"Bought amount": map[string]interface{}{
				"number": map[string]interface{}{
					"format": "number_with_commas",
				},
			},
			"Sold amount": map[string]interface{}{
				"number": map[string]interface{}{
					"format": "number_with_commas",
				},
			},
		},
	}
}
