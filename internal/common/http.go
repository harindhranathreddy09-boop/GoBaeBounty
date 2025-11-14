package common

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/ratelimit"
)

// HTTPClient wraps http.Client with rate limiting and concurrency control
type HTTPClient struct {
	client  *http.Client
	limiter *ratelimit.Limiter
	config  *Config
}

// NewHTTPClient creates an HTTP client with reasonable defaults
func NewHTTPClient(config *Config) *HTTPClient {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	return &HTTPClient{
		client:  client,
		limiter: config.Limiter,
		config:  config,
	}
}

// Do performs an HTTP request with rate limiting and concurrency control
func (h *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := h.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	defer h.limiter.Release()

	h.addDefaultHeaders(req)

	req = req.WithContext(ctx)
	return h.client.Do(req)
}

// Get performs HTTP GET request
func (h *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(ctx, req)
}

// Post performs HTTP POST request with content type and body
func (h *HTTPClient) Post(ctx context.Context, url, contentType string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return h.Do(ctx, req)
}

// Head performs HTTP HEAD request
func (h *HTTPClient) Head(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(ctx, req)
}

// addDefaultHeaders adds custom and default headers
func (h *HTTPClient) addDefaultHeaders(req *http.Request) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "GoBaeBounty/1.0 (https://github.com/harindhranathreddy09-boop/GoBaeBounty)")
	}
	if h.config.CustomHeaders != nil {
		for k, v := range h.config.CustomHeaders {
			req.Header.Set(k, v)
		}
	}
	if h.config.Cookies != nil {
		for k, v := range h.config.Cookies {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}
}
v
