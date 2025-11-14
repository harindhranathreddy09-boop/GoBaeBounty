package reporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

func TestGenerateReports(t *testing.T) {
	tmpDir := os.TempDir()
	config := &common.Config{
		OutputDir: tmpDir,
	}
	findings := []common.Finding{
		{
			ID:          "1",
			Title:       "Test XSS",
			Severity:    "high",
			CVSSScore:   7.5,
			Type:        "xss",
			URL:         "https://example.com/test",
			Description: "Test description",
			Remediation: "Test remediation",
			Payload:     "<script>alert(1)</script>",
		},
	}
	vulnResults := &common.VulnResults{Findings: findings}

	err := Generate(config, vulnResults)
	if err != nil {
		t.Errorf("Generate failed: %v", err)
	}

	mdPath := filepath.Join(tmpDir, "report.md")
	jsonPath := filepath.Join(tmpDir, "report.json")

	if _, err := os.Stat(mdPath); err != nil {
		t.Errorf("Markdown report not found: %v", err)
	}

	if _, err := os.Stat(jsonPath); err != nil {
		t.Errorf("JSON report not found: %v", err)
	}
}
