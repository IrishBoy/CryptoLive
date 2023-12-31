package domain

import (
	"time"
)

type NotionClient struct {
	BaseURL string
	APIKey  string
}
type GetTablesResponse struct {
	Object     string `json:"object"`
	Results    []any  `json:"results"`
	NextCursor any    `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	RequestID  string `json:"request_id"`
}
type GetTableRequest struct {
}

type GetTableResponse struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	Title          []any     `json:"title"`
	Description    []any     `json:"description"`
	Properties     struct {
	} `json:"properties"`
	Archived  bool `json:"archived"`
	IsInline  bool `json:"is_inline"`
	PublicURL any  `json:"public_url"`
}

func NewNotionClient(apiKey string) *NotionClient {
	return &NotionClient{
		BaseURL: "https://api.notion.com/v1",
		APIKey:  apiKey,
	}
}

func (nc *NotionClient) CreateRequestHeaders() map[string]string {
	return map[string]string{
		"Notion-Version": "2021-05-13",
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
