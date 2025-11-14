package jsparser

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Run performs JS files downloading and endpoint extraction concurrently
func Run(ctx context.Context, config *common.Config, jsFiles []common.JSFile) (*common.JSResults, error) {
	results := &common.JSResults{
		Endpoints:  []string{},
		Parameters: []string{},
		Secrets:    []common.Secret{},
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, config.Workers)

	for _, jsFile := range jsFiles {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sem <- struct{}{}:
			wg.Add(1)
			go func(jsf common.JSFile) {
				defer wg.Done()
				defer func() { <-sem }()
				content, err := fetchJSContent(ctx, jsf.URL)
				if err != nil {
					if config.Verbose {
						fmt.Printf("Failed to fetch %s: %v\n", jsf.URL, err)
					}
					return
				}

				endpoints := ExtractEndpoints(content)
				secretsMap := ExtractSecrets(content)

				mu.Lock()
				results.Endpoints = append(results.Endpoints, endpoints...)
				for stype, secretValues := range secretsMap {
					for _, sv := range secretValues {
						results.Secrets = append(results.Secrets, common.Secret{
							Type:  stype,
							Value: sv,
							File:  jsf.URL,
						})
					}
				}
				mu.Unlock()
			}(jsFile)
		}
	}

	wg.Wait()

	// Deduplicate
	results.Endpoints = common.Deduplicate(results.Endpoints)

	return results, nil
}

// fetchJSContent downloads JavaScript file content from a URL
func fetchJSContent(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
