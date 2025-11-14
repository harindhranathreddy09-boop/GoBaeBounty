package common

import "github.com/harindhranathreddy09-boop/GoBaeBounty/internal/ratelimit"

// Config holds global configuration for the tool
type Config struct {
	Target        string
	OutputDir     string
	Workers       int
	MaxRate       int
	CrawlDepth    int
	IgnoreRobots  bool
	Intrusive     bool
	Verbose       bool
	Limiter       *ratelimit.Limiter
	CustomHeaders map[string]string
	Cookies       map[string]string
}

// DiscoveryResults from passive discovery
type DiscoveryResults struct {
	Subdomains     []string
	HistoricalURLs []string
	DNSRecords     map[string][]string
}

// CrawlResults from crawling
type CrawlResults struct {
	Pages   []Page
	JSFiles []JSFile
	Forms   []Form
}

// Page represents HTML page crawled
type Page struct {
	URL        string
	StatusCode int
	Headers    map[string][]string
	BodySize   int
	CrawlTime  int64 // milliseconds
}

// JSFile represents JS resource found
type JSFile struct {
	URL string
}

// Form represents HTML form found
type Form struct {
	URL    string
	Action string
	Method string
	Inputs []FormInput
}

// FormInput represents HTML form input
type FormInput struct {
	Name  string
	Type  string
	Value string
}

// JSResults from JS parser
type JSResults struct {
	Endpoints  []string
	Parameters []string
	Secrets    []Secret
}

// Secret represents extracted sensitive string
type Secret struct {
	Type  string
	Value string
	File  string
}

// ScoredEndpoint represents an endpoint with metadata and score
type ScoredEndpoint struct {
	URL          string
	Score        int
	Type         string
	Methods      []string
	StatusCode   int
	ContentType  string
	Headers      map[string]string
	ResponseSize int
	HasJSON      bool
	CORS        bool
	CORSOrigins  []string
	AuthRequired bool
	Parameters   []string
}

// EndpointResults holds prioritized endpoints
type EndpointResults struct {
	All            []ScoredEndpoint
	HighPriority   int
	MediumPriority int
	LowPriority    int
}

// FuzzedParam represents fuzzed parameter input
type FuzzedParam struct {
	Name     string
	Endpoint string
	Method   string
	Value    string
	Response string
}

// FuzzResults holds results of fuzzing
type FuzzResults struct {
	ValidPaths []string
	Parameters []FuzzedParam
}

// VulnResults holds vulnerability findings
type VulnResults struct {
	Findings []Finding
}

// Finding represents discovered vulnerability
type Finding struct {
	ID          string
	Title       string
	Description string
	Severity    string
	CVSSScore   float64
	Type        string
	URL         string
	Parameter   string
	Payload     string
	Evidence    string
	Request     string
	Response    string
	Impact      string
	Remediation string
	References  []string
	Timestamp   string
	Verified    bool
	Confidence  float64
	Tags        []string
}

// CountBySeverity counts findings of severity
func (vr *VulnResults) CountBySeverity(sev string) int {
	count := 0
	for _, f := range vr.Findings {
		if f.Severity == sev {
			count++
		}
	}
	return count
}

// Deduplicate string slice
func Deduplicate(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}
	return result
}
