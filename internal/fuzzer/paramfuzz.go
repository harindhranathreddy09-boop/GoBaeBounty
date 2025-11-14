package fuzzer

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// RunParameterFuzzing runs parameter fuzzing for discovered endpoints
func RunParameterFuzzing(ctx context.Context, endpoints []common.ScoredEndpoint, config *common.Config, paramNames []string, payloads []string) ([]common.FuzzedParam, error) {
	fuzzedParams := []common.FuzzedParam{}
	client := common.NewHTTPClient(config)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Workers)

	for _, ep := range endpoints {
		baseURL := ep.URL
		for _, param := range paramNames {
			for _, payload := range payloads {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case sem <- struct{}{}:
					wg.Add(1)

					go func(e common.ScoredEndpoint, pName, pl string) {
						defer wg.Done()
						defer func() { <-sem }()
						fullURL, err := addOrReplaceQueryParam(e.URL, pName, pl)
						if err != nil {
							return
						}

						reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
						defer cancel()

						resp, err := client.Get(reqCtx, fullURL)
						if err != nil {
							return
						}
						defer resp.Body.Close()

						if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
							mu.Lock()
							fuzzedParams = append(fuzzedParams, common.FuzzedParam{
								Endpoint: e.URL,
								Name:     pName,
								Method:   "GET",
								Value:    pl,
								Response: resp.Status,
							})
							mu.Unlock()
						}
					}(ep, param, payload)
				}
			}
		}
	}

	wg.Wait()
	return fuzzedParams, nil
}

func addOrReplaceQueryParam(rawURL, param, value string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set(param, value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
