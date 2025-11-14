package jsparser

import (
	"context"
	"testing"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestExtractEndpoints(t *testing.T) {
	js := `
		var apiUrl = "https://api.example.com/v1/get";
		fetch("/api/users");
		let graphql = "/graphql";
	`
	got := ExtractEndpoints(js)
	expected := []string{
		"https://api.example.com/v1/get",
		"/api/users",
		"/graphql",
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d endpoints, got %d", len(expected), len(got))
	}
	for i, ep := range expected {
		if got[i] != ep {
			t.Errorf("expected endpoint %s, got %s", ep, got[i])
		}
	}
}

func TestExtractSecrets(t *testing.T) {
	js := `
		const apiKey = "secret123";
		const token = "abc.def.ghi";
	`
	secrets := ExtractSecrets(js)
	if len(secrets["api_key"]) == 0 {
		t.Error("expected to find api_key secret")
	}
	if len(secrets["token"]) == 0 {
		t.Error("expected to find token secret")
	}
}

func TestRun(t *testing.T) {
	jsFiles := []common.JSFile{
		{URL: "https://example.com/test.js"},
	}
	config := &common.Config{
		Workers: 2,
		Verbose: false,
	}
	ctx := context.Background()
	results, err := Run(ctx, config, jsFiles)
	if err != nil {
		t.Errorf("Run failed: %v", err)
	}
	if results == nil {
		t.Error("Expected results")
	}
}
