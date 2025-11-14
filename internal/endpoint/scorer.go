package endpoint

import (
	"errors"
	"strings"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Run scores endpoints and partitions by priority
func Run(ctx context.Context, config *common.Config, urls []string) (*common.EndpointResults, error) {
	if len(urls) == 0 {
		return nil, errors.New("no endpoints to score")
	}

	scoredEndpoints, err := Fingerprint(ctx, urls, config)
	if err != nil {
		return nil, err
	}

	results := &common.EndpointResults{
		All:            scoredEndpoints,
		HighPriority:   0,
		MediumPriority: 0,
		LowPriority:    0,
	}

	for _, ep := range scoredEndpoints {
		switch {
		case ep.Score >= 20:
			results.HighPriority++
		case ep.Score >= 10:
			results.MediumPriority++
		default:
			results.LowPriority++
		}
	}

	// Optionally, sort the results by score descending
	// Sorting omitted here for brevity

	return results, nil
}

// FilterEndpointsByStatus filters endpoints by status codes
func FilterEndpointsByStatus(endpoints []common.ScoredEndpoint, allowedStatusCodes []int) []common.ScoredEndpoint {
	filtered := []common.ScoredEndpoint{}
	for _, ep := range endpoints {
		for _, code := range allowedStatusCodes {
			if ep.StatusCode == code {
				filtered = append(filtered, ep)
				break
			}
		}
	}
	return filtered
}

// GetAdminEndpoints returns only those endpoints classified as admin
func GetAdminEndpoints(endpoints []common.ScoredEndpoint) []common.ScoredEndpoint {
	admins := []common.ScoredEndpoint{}
	for _, ep := range endpoints {
		if strings.Contains(strings.ToLower(ep.URL), "admin") {
			admins = append(admins, ep)
		}
	}
	return admins
}
