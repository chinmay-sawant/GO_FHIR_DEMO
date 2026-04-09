package utils

import (
	"context"
	"io"
	"net/http"
	"time"
)

// SharedClient is a pre-configured HTTP client shared across the package.
var SharedClient = &http.Client{
	Timeout: 10 * time.Second,
}

// DrainAndClose closes the response body after draining it to ensure the connection can be reused.
func DrainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}

// PrepareRequestWithTimeout creates an HTTP request with a timeout-aware context.
func PrepareRequestWithTimeout(ctx context.Context, method, url string, body io.Reader, timeout time.Duration) (*http.Request, context.CancelFunc, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	req, err := http.NewRequestWithContext(timeoutCtx, method, url, body)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return req, cancel, nil
}
