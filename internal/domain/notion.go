package domain

import (
	"time"
)

type NotionClient struct {
	BaseURL string
	APIKey  string
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

func (nc *NotionClient) UpdateTablePayload(coinPrice, profitValue float64, operationID string) map[string]interface{} {
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
			"ID": map[string]interface{}{
				"id":   "title",
				"type": "title",
				"title": []map[string]interface{}{
					{
						"type": "text",
						"text": map[string]interface{}{
							"content": operationID,
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
						"plain_text": operationID,
						"href":       false,
					},
				},
			},
		},
	}
}
