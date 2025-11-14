package vulncheck

import (
	"context"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// VulnCheck defines interface for a vulnerability check plugin
type VulnCheck interface {
	Name() string
	Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error)
}

var (
	registeredChecks = []VulnCheck{}
)

// RegisterPlugin adds a vulnerability check plugin
func RegisterPlugin(c VulnCheck) {
	registeredChecks = append(registeredChecks, c)
}

// RunAllChecks runs all registered vulnerability checks and aggregates findings
func RunAllChecks(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	var allFindings []common.Finding
	for _, check := range registeredChecks {
		findings, err := check.Run(ctx, config, endpoints)
		if err != nil {
			if config.Verbose {
				// Log error but continue with others
				config.Verbose = true
			}
			continue
		}
		allFindings = append(allFindings, findings...)
	}
	return allFindings, nil
}
