package http_client

import (
	"net/http"
	"time"
)

type retryRoundTripper struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(r.retryDelay)
		}

		resp, err = r.transport.RoundTrip(req)
		if err == nil && resp.StatusCode < 400 {
			return resp, nil
		}

		if resp != nil {
			_ = resp.Body.Close()
		}
	}

	return resp, err
}
