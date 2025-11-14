package common

import (
    "net/url"
    "regexp"
    "strings"
)

// Deduplicate removes duplicate strings from slice
func Deduplicate(items []string) []string {
    seen := make(map[string]bool)
    var result []string
    for _, s := range items {
        s = strings.TrimSpace(s)
        if s != "" && !seen[s] {
            seen[s] = true
            result = append(result, s)
        }
    }
    return result
}

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

// ExtractDomain gets the hostname from a URL
func ExtractDomain(rawURL string) string {
    u, err := url.Parse(rawURL)
    if err != nil {
        return ""
    }
    return u.Hostname()
}

// IsJavaScriptFile returns true if it looks like a JS file/link
func IsJavaScriptFile(rawURL string) bool {
    re := regexp.MustCompile(`(?i)\.js(\?|$)`)
    return re.MatchString(rawURL)
}

// IsInScope returns true if the provided host is within allowed scope (add your logic)
func IsInScope(host string, scope []string) bool {
    // Example: simple check, update as needed
    for _, s := range scope {
        if host == s || strings.HasSuffix(host, "."+s) {
            return true
        }
    }
    return false
}
// Add any remaining utility functions here...
