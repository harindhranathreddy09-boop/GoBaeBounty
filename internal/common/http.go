package common

import (
    "bytes"
    "context"
    "crypto/tls"
    "io"
    "net"
    "net/http"
    "time"
)

type Limiter interface {
    Wait(context.Context) error
    Release()
}

// HTTPClient wraps http.Client with rate limiting and custom configuration
type HTTPClient struct {
    client  http.Client
    limiter Limiter
    config  Config
}

func NewHTTPClient(config Config) HTTPClient {
    tr := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   10 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSClientConfig:     &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12},
        TLSHandshakeTimeout: 10 * time.Second,
    }
    client := http.Client{
        Timeout:   30 * time.Second,
        Transport: tr,
    }
    return HTTPClient{client: client, limiter: config.Limiter, config: config}
}

func (h HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
    if err := h.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    defer h.limiter.Release()
    req = req.WithContext(ctx)
    h.addHeaders(req)
    return h.client.Do(req)
}

func (h HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    return h.Do(ctx, req)
}

func (h HTTPClient) Post(ctx context.Context, url, contentType string, body []byte) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", contentType)
    return h.Do(ctx, req)
}

// Add headers, cookies, etc. to request as needed
func (h HTTPClient) addHeaders(req *http.Request) {
    if req.Header.Get("User-Agent") == "" {
        req.Header.Set("User-Agent", "GoBaeBounty1.0 https://github.com/harindhranathreddy09-boop/GoBaeBounty")
    }
    // Add custom headers and cookies as needed from config...
}

// Helpers for reading response bodies
func ReadResponseBody(resp *http.Response) ([]byte, error) {
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func ReadResponseBodyLimited(resp *http.Response, maxSize int64) ([]byte, error) {
    defer resp.Body.Close()
    return io.ReadAll(io.LimitReader(resp.Body, maxSize))
}
func (h HTTPClient) Head(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
    if err != nil {
        return nil, err
    }
    return h.Do(ctx, req)
}
