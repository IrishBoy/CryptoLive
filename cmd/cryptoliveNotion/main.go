// cmd/main.go
package main

import (
	"github.com/IrishBoy/CryptoLive/internal/domain"
	"github.com/IrishBoy/CryptoLive/internal/providers/notion"
)

func main() {
	// Create an instance of the Notion type
	notionInstance := &notion.Notion{
		NotionClient: *domain.NewNotionClient(""),
	}

	// Call the GetTable method on the Notion instance
	notionInstance.GetDatabase("")
}
