package common

import (
	"fmt"
	"net/http"
)

func CreateRequest(method string, URL string) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	return req, nil
}

func CreateHTTPClient() *http.Client {
	return &http.Client{}
}

func AddHeaders(headers map[string]string, request *http.Request) *http.Request {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}

func FormatURL(baseURL string, group string, databaseID string) string {
	return fmt.Sprintf("%s/%s/%s/query", baseURL, group, databaseID)
}
