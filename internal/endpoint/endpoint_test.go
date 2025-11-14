package endpoint

import (
	"context"
	"testing"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestFingerprint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/status/403",
		"https://httpbin.org/json",
	}

	config := &common.Config{
		Workers: 5,
		MaxRate: 10,
	}

	results, err := Fingerprint(ctx, urls, config)
	if err != nil {
		t.Fatalf("Fingerprint failed: %v", err)
	}

	if len(results) != len(urls) {
		t.Errorf("Expected %d results, got %d", len(urls), len(results))
	}
}

func TestRunScorer(t *testing.T) {
	ctx := context.Background()

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/basic-auth/user/passwd",
		"https://httpbin.org/admin",
	}

	config := &common.Config{
		Workers: 3,
	}

	res, err := Run(ctx, config, urls)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if res.HighPriority == 0 {
		t.Errorf("Expected some high priority endpoints")
	}
}
