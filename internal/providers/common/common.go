// common/common.go

package common

import (
	"bytes"
	"net/http"
)

func CreateRequest(method string, URL string, payloadBytes []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(payloadBytes))

	if err != nil {
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

func CalculateAverage(list []float64) float64 {
	total := 0.0
	for _, num := range list {
		total += num
	}
	return total / float64(len(list))
}
