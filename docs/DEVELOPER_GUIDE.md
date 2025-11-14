# GoBaeBounty Developer Guide

Welcome to the GoBaeBounty development guide. This document explains project architecture, coding standards, and how to extend or maintain the tool.

## Architecture Overview

GoBaeBounty is composed of modular packages under `internal/` with clear responsibilities:

- `auth/` - Authorization file validation to enforce legal usage.
- `discovery/` - Passive asset discovery via historical sources and DNS.
- `crawler/` - Polite and concurrent web crawler with robots.txt respect.
- `jsparser/` - Downloading and parsing JavaScript to extract endpoints/secrets.
- `endpoint/` - Endpoint fingerprinting and scoring heuristics.
- `fuzzer/` - Endpoint-driven fuzzing against directories and parameters.
- `vulncheck/` - Plugin-based vulnerability detection framework.
- `reporter/` - Reporting in Markdown and JSON formats for bug bounty submissions.
- `ratelimit/` - Token-bucket rate limiter supporting concurrency limits.
- `common/` - Shared types and utilities.

## Code Conventions

- Idiomatic Go with clear naming conventions.
- Concurrency implemented with worker pools and channels.
- Context usage for cancellation and timeouts.
- Minimal third-party dependencies.
- Thorough unit testing of all critical paths.

## How to Add Vulnerability Checks

1. Implement a new struct conforming to `vulncheck.VulnCheck` interface with `Run` and `Name`.
2. Add your struct to registration in `checks.go` inside `init()`.
3. Implement scanning logic focusing on safe requests.
4. Provide test suites covering positive and negative cases.
5. Document usage and false positives.

## Extending JS Parser

- Add new regex patterns or AST-based heuristics inside `internal/jsparser/extractor.go`.
- Use concurrency-safe appending to results.
- Consider decoding encoded JS, handling dynamic imports, or minified code prettification.
- Include new secret detectors with test coverage.

## Reporting Formats

- Markdown report template located in `internal/reporter/templates.go`.
- JSON is automatically marshaled from `common.VulnResults`.
- Extend or add new formats by implementing methods in `reporter/generator.go`.

## Running Tests

go test ./... -v



- Integration tests in `test/` run a full scan cycle on local servers.
- Unit tests validate isolated modules.

## Contact & Contribution

- Fork, develop on feature branches, and submit well-documented pull requests.
- Adhere to the code style and provide tests.
- Report bugs or feature requests via GitHub issues.

---

Thank you for contributing to GoBaeBounty!
File 49: docs/EXTENDING.md
text
# Extending GoBaeBounty

This guide helps you extend GoBaeBounty with new features and vulnerability checks.

## Adding Vulnerability Plugins

1. Create new Go file under `internal/vulncheck/`.
2. Define a struct implementing `vulncheck.VulnCheck` interface:

type MyCheck struct{}

func (c *MyCheck) Name() string {
return "my-check"
}

func (c *MyCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
findings := []common.Finding{}
// Implement your check logic here
return findings, nil
}



3. Register your plugin inside `checks.go`:

func init() {
RegisterPlugin(&MyCheck{})
}

text

4. Write comprehensive tests covering both positive and false positive scenarios.

## Extending JS Parsing Logic

- Modify or add regex patterns in `internal/jsparser/extractor.go`.
- Implement AST parsing if required for deeper analysis.
- Add handlers for obfuscated or dynamically loaded JS.

## Adding New Wordlists

Place wordlists in `pkg/wordlists/` and reference them in the configuration YAML.

## Customizing Reports

- Modify or add templates in `internal/reporter/templates.go`.
- Extend `internal/reporter/generator.go` to support new formats like HTML or CSV.

## Working with Rate Limits and Concurrency

- Tune values in config YAML.
- Extend the rate limiter if new strategies are needed, e.g., per-target limits.

## Testing Your Changes

- Use `go test ./...` to run tests.
- Add new tests in `_test.go` files alongside the module you're extending.
- Use integration tests in `test/` for full pipeline verification.

---

Happy hacking and coding with GoBaeBounty!
