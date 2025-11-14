package endpoint

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Fingerprint probes endpoints to classify as API, admin, graphql, etc.
func Fingerprint(ctx context.Context, urls []string, config *common.Config) ([]common.ScoredEndpoint, error) {
	var results []common.ScoredEndpoint
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, config.Workers)

	client := common.NewHTTPClient(config)

	for _, url := range urls {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sem <- struct{}{}:
			wg.Add(1)

			go func(u string) {
				defer wg.Done()
				defer func() { <-sem }()

				ep := common.ScoredEndpoint{URL: u}

				reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				resp, err := client.Head(reqCtx, u)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				ep.StatusCode = resp.StatusCode
				ep.Headers = map[string]string{}
				for k, v := range resp.Header {
					if len(v) > 0 {
						ep.Headers[k] = v[0]
					}
				}

				ct := resp.Header.Get("Content-Type")
				if strings.Contains(ct, "json") {
					ep.HasJSON = true
				}
				ep.ContentType = ct

				// Detect CORS headers
				if resp.Header.Get("Access-Control-Allow-Origin") != "" {
					ep.CORS = true
					ep.CORSOrigins = []string{resp.Header.Get("Access-Control-Allow-Origin")}
				}

				// Heuristics scoring
				ep.Score = scoreEndpoint(u, ep)

				ep.Methods = []string{"HEAD", "GET", "OPTIONS"}

				// Mark type by heuristics
				ep.Type = classifyURL(u)

				mu.Lock()
				results = append(results, ep)
				mu.Unlock()
			}(url)
		}
	}

	wg.Wait()
	return results, nil
}

// Helper scoring function for endpoint importance heuristic
func scoreEndpoint(url string, ep common.ScoredEndpoint) int {
	score := 0
	lower := strings.ToLower(url)

	if strings.Contains(lower, "admin") {
		score += 10
	}
	if strings.Contains(lower, "api") {
		score += 8
	}
	if strings.Contains(lower, "graphql") {
		score += 9
	}
	if strings.Contains(lower, "auth") {
		score += 7
	}
	if strings.Contains(lower, "upload") {
		score += 8
	}

	if ep.HasJSON {
		score += 5
	}
	if ep.CORS {
		score += 3
	}
	if ep.StatusCode >= 200 && ep.StatusCode < 300 {
		score += 5
	}
	if ep.StatusCode == 401 || ep.StatusCode == 403 {
		score += 7
	}
	return score
}

// Helper classification of endpoint type
func classifyURL(u string) string {
	u = strings.ToLower(u)
	switch {
	case strings.Contains(u, "admin"):
		return "admin"
	case strings.Contains(u, "api"):
		return "api"
	case strings.Contains(u, "graphql"):
		return "graphql"
	default:
		return "unknown"
	}
}
