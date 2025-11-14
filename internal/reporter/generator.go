package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/harindhranathreddy09-boop/GoBaeBounty/internal/common"
)

// Generate creates report files in Markdown and JSON formats in output directory
func Generate(config *common.Config, vulnResults *common.VulnResults) error {
	outputDir := config.OutputDir

	// Sort findings by severity descending
	vulnResults.SortBySeverity()

	// Prepare Markdown report
	mdReport, err := generateMarkdownReport(vulnResults)
	if err != nil {
		return fmt.Errorf("failed to generate markdown report: %w", err)
	}

	mdFile := filepath.Join(outputDir, "report.md")
	if err := os.WriteFile(mdFile, []byte(mdReport), 0644); err != nil {
		return fmt.Errorf("failed to write markdown report: %w", err)
	}

	// Prepare JSON report
	jsonReport, err := generateJSONReport(vulnResults)
	if err != nil {
		return fmt.Errorf("failed to generate JSON report: %w", err)
	}

	jsonFile := filepath.Join(outputDir, "report.json")
	if err := os.WriteFile(jsonFile, []byte(jsonReport), 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

func generateMarkdownReport(vulnResults *common.VulnResults) (string, error) {
	tmpl, err := template.New("mdreport").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(ReportMarkdownTemplate)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, vulnResults)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func generateJSONReport(vulnResults *common.VulnResults) (string, error) {
	data, err := json.MarshalIndent(vulnResults, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
