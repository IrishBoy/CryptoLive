package domain

import (
	"time"
)

type NotionClient struct {
	BaseURL string
	// HTTPClient *http.Client
	APIKey string
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

func (nc *NotionClient) CreateRequestHeaders() map[string]string {
	return map[string]string{
		"Notion-Version": "2021-05-13",
		"Authorization":  "Bearer " + nc.APIKey,
		"Content-Type":   "application/json",
		"Accept":         "application/json",
	}
}
func NewNotionClient(apiKey string) *NotionClient {
	return &NotionClient{
		BaseURL: "https://api.notion.com/v1",
		APIKey:  apiKey,
	}
}
