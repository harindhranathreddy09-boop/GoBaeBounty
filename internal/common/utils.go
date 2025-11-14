package common

import (
	"crypto/md5"
	"encoding/hex"
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
func Deduplicate(strings []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strings {
		if _, exists := seen[s]; !exists {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// IsJavaScriptFile returns true if URL path looks like a JS file
func IsJavaScriptFile(rawURL string) bool {
	re := regexp.MustCompile(`(?i)\.js(\?.*)?$`)
	return re.MatchString(rawURL)
}
