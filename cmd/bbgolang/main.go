package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/auth"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/crawler"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/discovery"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/endpoint"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/fuzzer"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/jsparser"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/ratelimit"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/reporter"
	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/vulncheck"
)

func main() {
	var target, authFile, output string
	var workers, maxRate, depth int
	var ignoreRobots, intrusive, dryRun, verbose bool

	flag.StringVar(&target, "target", "", "Target domain (required)")
	flag.StringVar(&authFile, "auth-file", "", "Authorization file path (required)")
	flag.StringVar(&output, "o", "results", "Output directory")
	flag.IntVar(&workers, "workers", 50, "Number of concurrent workers")
	flag.IntVar(&maxRate, "max-rate", 100, "Max requests per second")
	flag.IntVar(&depth, "depth", 3, "Crawling depth")
	flag.BoolVar(&ignoreRobots, "ignore-robots", false, "Ignore robots.txt rules")
	flag.BoolVar(&intrusive, "intrusive", false, "Enable intrusive vulnerability checks")
	flag.BoolVar(&dryRun, "dry-run", false, "Authorization only")
	flag.BoolVar(&verbose, "v", false, "Verbose output")

	flag.Parse()

	if target == "" || authFile == "" {
		fmt.Println("Usage: --target and --auth-file are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	printBanner()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on interrupt
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		fmt.Println("\nReceived interrupt; shutting down...")
		cancel()
	}()

	// Validate auth
	if err := auth.ValidateAuthFile(authFile, target); err != nil {
		fmt.Fprintf(os.Stderr, "Authorization failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Authorization validated")

	if dryRun {
		fmt.Println("Dry run complete; exiting.")
		return
	}

	outDirAbs, err := filepath.Abs(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve output dir: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(outDirAbs, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	limiter := ratelimit.NewLimiter(maxRate, workers)

	cfg := &common.Config{
		Target:       target,
		OutputDir:    outDirAbs,
		Workers:      workers,
		MaxRate:      maxRate,
		CrawlDepth:   depth,
		IgnoreRobots: ignoreRobots,
		Intrusive:    intrusive,
		Verbose:      verbose,
		Limiter:      limiter,
	}

	fmt.Println("[1/8] Asset Discovery...")
	discoverResults, err := discovery.Run(ctx, cfg)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Discovered %d subdomains and %d URLs\n", len(discoverResults.Subdomains), len(discoverResults.HistoricalURLs))

	fmt.Println("[2/8] Crawling...")
	crawlResults, err := crawler.Run(ctx, cfg, discoverResults)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Crawled %d pages; found %d JavaScript files\n", len(crawlResults.Pages), len(crawlResults.JSFiles))

	fmt.Println("[3/8] JavaScript Parsing...")
	jsResults, err := jsparser.Run(ctx, cfg, crawlResults.JSFiles)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Extracted %d endpoints and %d parameters\n", len(jsResults.Endpoints), len(jsResults.Parameters))

	fmt.Println("[4/8] Endpoint Fingerprinting...")
	endpointResults, err := endpoint.Run(ctx, cfg, jsResults.Endpoints)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Prioritized %d endpoints\n", len(endpointResults.All))

	fmt.Println("[5/8] Fuzzing endpoints...")
	fuzzResults, err := fuzzer.Run(ctx, cfg, endpointResults)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Found %d valid paths\n", len(fuzzResults.ValidPaths))

	fmt.Println("[6/8] Vulnerability scanning...")
	vulnResults, err := vulncheck.Run(ctx, cfg, fuzzResults)
	if err != nil {
		exitError(err)
	}
	fmt.Printf("Found %d vulnerabilities\n", len(vulnResults.Findings))

	fmt.Println("[7/8] Generating reports...")
	if err := reporter.Generate(cfg, vulnResults); err != nil {
		exitError(err)
	}

	fmt.Println("[8/8] Scan complete. Reports generated.")

}

func printBanner() {
	fmt.Println("GoBaeBounty - Bug Bounty Automation Framework")
	fmt.Println("© 2025 by harindhranathreddy09-boop")
	fmt.Println("Run with --help for usage")
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
