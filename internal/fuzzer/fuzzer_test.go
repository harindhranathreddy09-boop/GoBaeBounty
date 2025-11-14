package fuzzer

import (
	"context"
	"testing"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/endpoint"
)

func TestRunDirectoryFuzzing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	config := &common.Config{
		Workers: 5,
		MaxRate: 10,
	}

	baseURL := "https://httpbin.org"

	// Use minimal wordlist for testing
	wordlists := []string{"./pkg/wordlists/common.txt"}

	paths, err := RunDirectoryFuzzing(ctx, baseURL, config, wordlists)
	if err != nil {
		t.Fatalf("RunDirectoryFuzzing failed: %v", err)
	}

	if len(paths) == 0 {
		t.Error("Expected some valid paths")
	}
}

func TestRunParameterFuzzing(t *testing.T) {
	ctx := context.Background()

	config := &common.Config{Workers: 3}

	endpoints := []common.ScoredEndpoint{
		{URL: "https://httpbin.org/get"},
	}

	paramNames := []string{"id"}
	payloads := []string{"test", "' OR '1'='1"}

	results, err := RunParameterFuzzing(ctx, endpoints, config, paramNames, payloads)
	if err != nil {
		t.Fatalf("RunParameterFuzzing failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected some fuzzed parameters")
	}
}
