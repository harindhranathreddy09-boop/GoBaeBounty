package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/temoto/robotstxt"
)

// RobotsManager manages robots.txt parsing and caching per domain
type RobotsManager struct {
	cache      map[string]*robotstxt.RobotsData
	cacheMutex sync.RWMutex
	userAgent  string
	delayCache map[string]time.Duration
	delayMutex sync.RWMutex
}

// NewRobotsManager creates a new RobotsManager
func NewRobotsManager(userAgent string) *RobotsManager {
	return &RobotsManager{
		cache:      make(map[string]*robotstxt.RobotsData),
		userAgent:  userAgent,
		delayCache: make(map[string]time.Duration),
	}
}

// GetRobotsTxt fetches and parses robots.txt for given URL
func (rm *RobotsManager) GetRobotsTxt(urlStr string) (*robotstxt.RobotsData, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	host := u.Scheme + "://" + u.Host

	rm.cacheMutex.RLock()
	if data, ok := rm.cache[host]; ok {
		rm.cacheMutex.RUnlock()
		return data, nil
	}
	rm.cacheMutex.RUnlock()

	robotsURL := host + "/robots.txt"
	resp, err := http.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	rm.cacheMutex.Lock()
	rm.cache[host] = data
	rm.cacheMutex.Unlock()

	return data, nil
}

// Allowed returns true if given URL is allowed to be crawled according to robots.txt
func (rm *RobotsManager) Allowed(urlStr string) (bool, error) {
	data, err := rm.GetRobotsTxt(urlStr)
	if err != nil {
		// Assume allowed if robots.txt fetch fails
		return true, nil
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}
	return data.TestAgent(u.Path, rm.userAgent), nil
}

// CrawlDelay returns crawl delay for the domain of given URL or zero if none specified
func (rm *RobotsManager) CrawlDelay(urlStr string) time.Duration {
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0
	}
	host := u.Scheme + "://" + u.Host

	rm.delayMutex.RLock()
	delay, ok := rm.delayCache[host]
	rm.delayMutex.RUnlock()
	if ok {
		return delay
	}

	data, err := rm.GetRobotsTxt(urlStr)
	if err != nil {
		return 0
	}

	delay = data.CrawlDelay(rm.userAgent)

	rm.delayMutex.Lock()
	rm.delayCache[host] = delay
	rm.delayMutex.Unlock()

	return delay
}

// IsNoIndex checks if the URL path is disallowed by robots.txt as noindex or nofollow
func (rm *RobotsManager) IsNoIndex(urlStr string) (bool, error) {
	data, err := rm.GetRobotsTxt(urlStr)
	if err != nil {
		return false, err
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}
	group := data.FindGroup(rm.userAgent)
	if group == nil {
		return false, nil
	}
	// Check Disallow fields
	for _, path := range group.Disallow {
		if strings.HasPrefix(u.Path, path) {
			return true, nil
		}
	}
	return false, nil
}

// Example usage:
//
// rm := NewRobotsManager("GoBaeBounty")
// allowed, err := rm.Allowed("https://example.com/page")
//
// This should be integrated into the crawler to respect politeness and avoid forbidden paths.
