package test

import (
	"context"
	"testing"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/discovery"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/jsparser"
)

func TestFullScanPipeline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	config := &common.Config{
		Target: "localhost:8080",
		Workers: 5,
		CrawlDepth: 2,
		IgnoreRobots: true,
	}

	// Passive discovery (mocked)
	discoveryResults, err := discovery.Run(ctx, config)
	if err != nil {
		t.Fatalf("Discovery failed: %v", err)
	}
	if len(discoveryResults.Subdomains) == 0 {
		t.Fatal("Expected at least one subdomain")
	}

	// Crawling
	crawlResults, err := jsparser.Run(ctx, config, []common.JSFile{
		{URL: "http://localhost:8080/static/app.js"},
	})
	if err != nil {
		t.Fatalf("JS parsing failed: %v", err)
	}
	if len(crawlResults.Endpoints) == 0 {
		t.Error("Expected endpoints to be extracted")
	}
}
