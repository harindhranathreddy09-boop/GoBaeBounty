package vulncheck

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

type SQLiCheck struct{}

func (s *SQLiCheck) Name() string {
	return "SQLInjection"
}

func (s *SQLiCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	var findings []common.Finding

	injectionPayload := "' OR '1'='1"

	for _, ep := range endpoints {
		fullURL, err := addInjectionToURL(ep.URL, injectionPayload)
		if err != nil {
			continue
		}

		client := common.NewHTTPClient(config)
		resp, err := client.Get(ctx, fullURL)
		if err != nil {
			continue
		}
		resp.Body.Close()

		// Simplified logic: if response contains certain strings, flag
		if checkSQLiIndicators(resp) {
			findings = append(findings, common.Finding{
				ID:         "sqli-" + ep.URL,
				Title:      "Possible SQL Injection",
				Severity:   "high",
				Type:       "sqli",
				URL:        ep.URL,
				Description: fmt.Sprintf("Endpoint %s might be vulnerable to SQL injection.", ep.URL),
				Remediation: "Use parameterized queries and sanitize inputs.",
			})
		}
	}

	return findings, nil
}

// addInjectionToURL adds payload as a query param to URL
func addInjectionToURL(rawURL, payload string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for param := range q {
		q.Set(param, payload)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func checkSQLiIndicators(resp *http.Response) bool {
	// Placeholder; you would parse body and look for SQL errors
	return false
}
