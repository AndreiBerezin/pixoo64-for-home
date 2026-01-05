package pixoo64

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	addr       string
	httpClient *http.Client
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Post(data any) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/post", c.addr), bytes.NewBuffer(dataBytes))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
