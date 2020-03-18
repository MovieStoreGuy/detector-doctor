package client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

const (
	// UserAgent is set to ensure that application is mistaken as a bot
	UserAgent = "detector-doctor"

	// ContentType is set on each out going request that is not a GET method
	ContentType = "application/json"

	// DefaultTimeout used when making outbound requests
	DefaultTimeout = 10 * time.Second
)

// NewConfiguredClient returns a configured client and allows the caller to
// update any client settings they desire
func NewConfiguredClient(opts ...func(*http.Client)) *http.Client {
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 2 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			MaxIdleConns:    6,
			IdleConnTimeout: DefaultTimeout,
		},
		Timeout: DefaultTimeout,
	}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

// NewConfiguredRequestFunc caches the token and applies them to each newly created request.
func NewConfiguredRequestFunc(token string) func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	// No-Op field, used just to cache token within the request generation
	return func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-SF-TOKEN", token)
		switch method {
		case http.MethodGet:
			// Do nothing
		default:
			req.Header.Add("Content-Type", ContentType)
		}
		// Override the default language user agent to be clearly defined
		req.Header.Add("User-Agent", UserAgent)
		return req, nil
	}
}
