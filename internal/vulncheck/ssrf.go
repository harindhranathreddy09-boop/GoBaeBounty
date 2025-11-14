package vulncheck

import (
	"context"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

type SSRFCheck struct{}

func (s *SSRFCheck) Name() string {
	return "SSRF"
}

func (s *SSRFCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	// Placeholder for SSRF detection logic
	return nil, nil
}
