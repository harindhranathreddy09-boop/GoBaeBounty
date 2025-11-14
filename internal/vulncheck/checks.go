package vulncheck

import (
	"context"
	"fmt"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func init() {
	RegisterPlugin(&OpenRedirectCheck{})
	RegisterPlugin(&AuthBypassCheck{})
	// Add other checks later...
}

// OpenRedirectCheck detects open redirect vulnerability patterns
type OpenRedirectCheck struct{}

func (o *OpenRedirectCheck) Name() string {
	return "OpenRedirect"
}

func (o *OpenRedirectCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	var findings []common.Finding
	for _, ep := range endpoints {
		// Simulate by checking if URL contains "redirect" param - placeholder for actual detection
		if containsRedirectParam(ep.URL) {
			f := common.Finding{
				ID:          "openredirect-" + ep.URL,
				Title:       "Potential Open Redirect",
				Severity:    "medium",
				CVSSScore:   6.1,
				Type:        "openredirect",
				URL:         ep.URL,
				Description: "URL contains redirect parameter that might be vulnerable.",
				Remediation: "Validate and sanitize redirect parameters.",
			}
			findings = append(findings, f)
		}
	}
	return findings, nil
}

func containsRedirectParam(url string) bool {
	return (len(url) > 0) && (containsStringCaseInsensitive(url, "redirect=") || containsStringCaseInsensitive(url, "redir="))
}

func containsStringCaseInsensitive(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (len(needle) == 0 || containsFold(haystack, needle))
}

func containsFold(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if strings.EqualFold(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

// AuthBypassCheck checks for common auth bypass patterns
type AuthBypassCheck struct{}

func (a *AuthBypassCheck) Name() string {
	return "AuthBypass"
}

func (a *AuthBypassCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	var findings []common.Finding
	// This is a placeholder stub; real checks would attempt HTTP requests
	return findings, nil
}
