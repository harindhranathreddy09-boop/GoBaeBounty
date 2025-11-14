package jsparser

import (
	"regexp"
)

// Patterns to extract API endpoints and relevant strings from JS content
var (
	// Regexes for URLs and API endpoints
	urlPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)["'](https?://[^\s"'\\]+)["']`),
		regexp.MustCompile(`(?i)["'](/api/[a-zA-Z0-9_\-/]+)["']`),
		regexp.MustCompile(`(?i)["'](/v[0-9]+[^\s"'\\]*)["']`),
		regexp.MustCompile(`(?i)["'](/graphql[^\s"'\\]*)["']`),
	}

	// Regex for secret-like strings
	secretPatterns = map[string]*regexp.Regexp{
		"api_key":     regexp.MustCompile(`(?i)api[_-]?key["']?\s*[:=]\s*["']([a-zA-Z0-9_\-]+)["']`),
		"token":       regexp.MustCompile(`(?i)(?:token|auth)["']?\s*[:=]\s*["']([a-zA-Z0-9_\-\.]+)["']`),
		"csrf":        regexp.MustCompile(`(?i)csrf["']?\s*[:=]\s*["'](.+?)["']`),
		"jwt":         regexp.MustCompile(`eyJ[A-Za-z0-9\-_]*\.[A-Za-z0-9\-_]*\.[A-Za-z0-9\-_]*`),
		"password":    regexp.MustCompile(`(?i)password["']?\s*[:=]\s*["']([^"']+)["']`),
		"private_key": regexp.MustCompile(`-----BEGIN (?:RSA|EC|OPENSSH) PRIVATE KEY-----`),
	}
)

// ExtractEndpoints extracts and returns endpoints from JavaScript code content
func ExtractEndpoints(jsContent string) []string {
	endpoints := []string{}
	matched := map[string]bool{}

	for _, pattern := range urlPatterns {
		matches := pattern.FindAllStringSubmatch(jsContent, -1)
		for _, match := range matches {
			if len(match) > 1 && !matched[match[1]] {
				endpoints = append(endpoints, match[1])
				matched[match[1]] = true
			}
		}
	}
	return endpoints
}

// ExtractSecrets extracts secret-like strings from JavaScript content
func ExtractSecrets(jsContent string) map[string][]string {
	secrets := make(map[string][]string)
	for secretType, pattern := range secretPatterns {
		matches := pattern.FindAllStringSubmatch(jsContent, -1)
		for _, match := range matches {
			if len(match) > 1 {
				secrets[secretType] = append(secrets[secretType], match[1])
			}
		}
	}
	return secrets
}
