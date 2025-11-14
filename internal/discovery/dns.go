package discovery

import (
	"context"
	"fmt"
	"net"
)

// ResolveDNS performs DNS A record lookup
func ResolveDNS(ctx context.Context, domain string) ([]string, error) {
	resolver := &net.Resolver{PreferGo: true}

	ips, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed for %s: %w", domain, err)
	}

	return ips, nil
}

// LookupCNAME resolves CNAME record for a domain
func LookupCNAME(ctx context.Context, domain string) (string, error) {
	resolver := &net.Resolver{PreferGo: true}

	cname, err := resolver.LookupCNAME(ctx, domain)
	if err != nil {
		return "", fmt.Errorf("CNAME lookup failed for %s: %w", domain, err)
	}

	return cname, nil
}
