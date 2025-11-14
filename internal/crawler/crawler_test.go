package crawler

import (
	"context"
	"testing"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestRunCrawler(t *testing.T) {
	config := &common.Config{
		Target:      "example.com",
		Workers:     2,
		CrawlDepth:  2,
		IgnoreRobots: true,
		Limiter:     nil, // no rate limit for test
	}

	discoveryResults := &common.DiscoveryResults{
		Subdomains: []string{"example.com"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := Run(ctx, config, discoveryResults)
	if err != nil {
		t.Fatalf("Crawler run failed: %v", err)
	}

	if len(results.Pages) == 0 {
		t.Error("Expected to crawl some HTML pages")
	}

	if len(results.JSFiles) == 0 {
		t.Error("Expected to find some JavaScript files")
	}
}
