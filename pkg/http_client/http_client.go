package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &retryRoundTripper{
				transport:  http.DefaultTransport,
				maxRetries: 2,
				retryDelay: 500 * time.Millisecond,
			},
		},
	}
}

func (c *Client) Post(url string, bodyObject any, target any) error {
	bodyBytes, err := json.Marshal(bodyObject)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return c.do(request, target)
}

func (c *Client) Get(url string, target any) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return c.do(request, target)
}

func (c *Client) do(request *http.Request, target any) error {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code: %d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch v := target.(type) {
	case nil:
		// do nothing
	case *string:
		*v = string(responseBytes)
	case *[]byte:
		*v = responseBytes
	default:
		if err := json.Unmarshal(responseBytes, target); err != nil {
			return err
		}
	}

	return nil
}
