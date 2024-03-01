package domain

import "time"

type SearchResponse struct {
	Object         string             `json:"object"`
	Results        []SearchObjectType `json:"results"`
	NextCursor     any                `json:"next_cursor"`
	HasMore        bool               `json:"has_more"`
	Type           string             `json:"type"`
	PageOrDatabase any                `json:"page_or_database"`
	RequestID      string             `json:"request_id"`
}

type SearchObjectType struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	CreatedBy      struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"created_by"`
	LastEditedBy struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"last_edited_by"`
	Cover  any `json:"cover"`
	Icon   any `json:"icon"`
	Parent struct {
		Type        string `json:"type"`
		DatabaseID  string `json:"database_id,omitempty"`
		PageID      string `json:"page_id,omitempty"`
		IsWorkspace bool   `json:"workspace,omitempty"`
	} `json:"parent"`
	Archived   bool `json:"archived"`
	Properties struct {
	} `json:"properties"`
	URL       string `json:"url"`
	PublicURL any    `json:"public_url"`
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

func (nc *NotionClient) UpdateTablePayload(coinPrice, profit float64, profitValue float64) map[string]interface{} {
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"Current Coin Price": map[string]interface{}{
				"type":   "number",
				"number": coinPrice,
			},
			"Profit": map[string]interface{}{
				"type":   "number",
				"number": profit,
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
				"annotations": map[string]interface{}{
					"bold":          false,
					"italic":        false,
					"strikethrough": false,
					"underline":     false,
					"code":          false,
					"color":         "default",
				},
				"plain_text": "Operations",
				"href":       nil,
			},
		},
		"properties": map[string]interface{}{
			"Coin bought": map[string]interface{}{
				"name":   "Coin bought",
				"type":   "select",
				"select": map[string]interface{}{"options": []interface{}{}},
			},
			"Bought Amount": map[string]interface{}{
				"name":   "Bought Amount",
				"type":   "number",
				"number": map[string]interface{}{"format": "number_with_commas"},
			},
			"Coin Sold": map[string]interface{}{
				"name":   "Coin Sold",
				"type":   "select",
				"select": map[string]interface{}{"options": []interface{}{}},
			},
			"Sold Amount": map[string]interface{}{
				"name":   "Sold Amount",
				"type":   "number",
				"number": map[string]interface{}{"format": "number_with_commas"},
			},
			"Profit": map[string]interface{}{
				"name":   "Profit",
				"type":   "number",
				"number": map[string]interface{}{"format": "number_with_commas"},
			},
			"Profit Percentage": map[string]interface{}{
				"name":   "Profit Percentage",
				"type":   "number",
				"number": map[string]interface{}{"format": "percent"},
			},
			"Operation Type": map[string]interface{}{
				"id":    "title",
				"name":  "Operation Type",
				"type":  "title",
				"title": map[string]interface{}{},
			},
		},
	}
}
