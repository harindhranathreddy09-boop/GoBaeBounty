package vulncheck

import (
	"context"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

type IDORCheck struct{}

func (i *IDORCheck) Name() string {
	return "IDOR"
}

func (i *IDORCheck) Run(ctx context.Context, config *common.Config, endpoints []common.ScoredEndpoint) ([]common.Finding, error) {
	// Placeholder for IDOR vulnerability checks
	return nil, nil
}
