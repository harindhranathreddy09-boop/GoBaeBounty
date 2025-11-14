package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestRunPassiveDiscovery(t *testing.T) {
	config := &common.Config{
		Target: "example.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := Run(ctx, config)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if len(results.Subdomains) == 0 {
		t.Errorf("Expected some subdomains, got none")
	}
	if results.Subdomains[0] != "example.com" {
		t.Logf("First domain: %s", results.Subdomains[0])
	}
}

func TestResolveDNS(t *testing.T) {
	ctx := context.Background()
	ips, err := ResolveDNS(ctx, "example.com")
	if err != nil {
		t.Errorf("DNS Resolve failed: %v", err)
	}
	if len(ips) == 0 {
		t.Error("Expected at least one IP address")
	}
}

func TestLookupCNAME(t *testing.T) {
	ctx := context.Background()
	cname, err := LookupCNAME(ctx, "www.example.com")
	if err != nil {
		t.Errorf("CNAME Lookup failed: %v", err)
	}
	if cname == "" {
		t.Error("Expected a CNAME record")
	}
}
