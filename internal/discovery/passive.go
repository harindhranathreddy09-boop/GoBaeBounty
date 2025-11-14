package discovery

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Run performs passive asset discovery: subdomains, Wayback URLs, GAU integration
func Run(ctx context.Context, config *common.Config) (*common.DiscoveryResults, error) {
	results := &common.DiscoveryResults{
		Subdomains:     []string{config.Target},
		HistoricalURLs: []string{},
		DNSRecords:     make(map[string][]string),
	}

	// Collect URLs via waybackurls if installed
	if urls, err := runWaybackURLs(ctx, config.Target); err == nil {
		results.HistoricalURLs = append(results.HistoricalURLs, urls...)
	}

	// Collect URLs via gau if installed
	if urls, err := runGAU(ctx, config.Target); err == nil {
		results.HistoricalURLs = append(results.HistoricalURLs, urls...)
	}

	// Deduplicate URLs
	results.HistoricalURLs = common.Deduplicate(results.HistoricalURLs)

	// Extract subdomains from URLs
	for _, url := range results.HistoricalURLs {
		domain := common.ExtractDomain(url)
		if domain != "" && common.IsInScope(url, config.Target) {
			results.Subdomains = append(results.Subdomains, domain)
		}
	}

	results.Subdomains = common.Deduplicate(results.Subdomains)

	return results, nil
}

// runWaybackURLs uses the 'waybackurls' external tool to get URLs
func runWaybackURLs(ctx context.Context, target string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "waybackurls", target)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("waybackurls not available or failed: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var urls []string
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls, nil
}

// runGAU uses the 'gau' external tool to get URLs
func runGAU(ctx context.Context, target string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "gau", target)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gau not available or failed: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var urls []string
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls, nil
}
