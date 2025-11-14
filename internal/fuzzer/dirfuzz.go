package fuzzer

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// RunDirectoryFuzzing performs recursive directory and file fuzzing on endpoints
func RunDirectoryFuzzing(ctx context.Context, baseURL string, config *common.Config, wordlists []string) ([]string, error) {
	validPaths := []string{}
	client := common.NewHTTPClient(config)

	// Load combined wordlist entries
	words, err := loadWordlists(wordlists)
	if err != nil {
		return nil, fmt.Errorf("failed to load wordlists: %w", err)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Workers)

	for _, word := range words {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sem <- struct{}{}:
			wg.Add(1)
			go func(w string) {
				defer wg.Done()
				defer func() { <-sem }()

				url := strings.TrimRight(baseURL, "/") + "/" + w
				reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				resp, err := client.Head(reqCtx, url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
					mu.Lock()
					validPaths = append(validPaths, url)
					mu.Unlock()
				}
			}(word)
		}
	}

	wg.Wait()
	return validPaths, nil
}

// loadWordlists loads words from multiple wordlist file paths
func loadWordlists(paths []string) ([]string, error) {
	wordsMap := map[string]bool{}
	for _, path := range paths {
		file, err := os.Open(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("could not open wordlist %s: %v", path, err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			w := strings.TrimSpace(scanner.Text())
			if w != "" {
				wordsMap[w] = true
			}
		}
		file.Close()
	}
	var words []string
	for w := range wordsMap {
		words = append(words, w)
	}
	return words, nil
}
