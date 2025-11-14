package vulncheck

import (
	"context"
	"testing"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestPluginRegistration(t *testing.T) {
	if len(registeredChecks) == 0 {
		t.Error("No vulnerability plugins registered")
	}
}

func TestRunAllChecksEmpty(t *testing.T) {
	ctx := context.Background()
	config := &common.Config{}
	eps := []common.ScoredEndpoint{}
	findings, err := RunAllChecks(ctx, config, eps)
	if err != nil {
		t.Errorf("RunAllChecks failed: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("Expected zero findings, got %d", len(findings))
	}
}
