package common

import (
	"net/url"
	"regexp"
	"strings"
)

// NormalizeURL removes fragments and standardizes URL for deduplication
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	return u.String(), nil
}

// ExtractDomain extracts domain part of URL
func ExtractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// IsInScope checks if URL domain belongs to scope domain
func IsInScope(targetURL, scopeDomain string) bool {
	domain := ExtractDomain(targetURL)
	return domain == scopeDomain || strings.HasSuffix(domain, "."+scopeDomain)
}

// Deduplicate removes duplicate strings from slice
func Deduplicate(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range items {
		if !seen[s] {
			result = append(result, s)
			seen[s] = true
		}
	}
	return result
}

// IsJavaScriptFile returns true if URL path looks like a JS file
func IsJavaScriptFile(rawURL string) bool {
	re := regexp.MustCompile(`(?i)\.js(\?.*)?$`)
	return re.MatchString(rawURL)
}
package common

import (
	"net/url"
	"regexp"
	"strings"
)

// NormalizeURL removes fragments and standardizes URL for deduplication
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	return u.String(), nil
}

// ExtractDomain extracts domain part of URL
func ExtractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// IsInScope checks if URL domain belongs to scope domain
func IsInScope(targetURL, scopeDomain string) bool {
	domain := ExtractDomain(targetURL)
	return domain == scopeDomain || strings.HasSuffix(domain, "."+scopeDomain)
}

// Deduplicate removes duplicate strings from slice
func Deduplicate(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range items {
		if !seen[s] {
			result = append(result, s)
			seen[s] = true
		}
	}
	return result
}

// IsJavaScriptFile returns true if URL path looks like a JS file
func IsJavaScriptFile(rawURL string) bool {
	re := regexp.MustCompile(`(?i)\.js(\?.*)?$`)
	return re.MatchString(rawURL)
}
