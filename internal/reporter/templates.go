package reporter

// Template strings for markdown and JSON reports

const ReportMarkdownTemplate = `# Bug Bounty Report

{{range .Findings}}
## Issue: {{.Title}} ({{.Severity | title}})
- URL: {{.URL}}
- Type: {{.Type}}
- Severity: {{.Severity}}
- CVSS Score: {{printf "%.1f" .CVSSScore}}
- Description: {{.Description}}
- Remediation: {{.Remediation}}

### Proof of Concept
{{.Payload}}

---

{{end}}
`

const ReportJSONTemplate = `
{
  "findings": [
    {{range $index, $f := .Findings -}}
    {{if $index}},{{end}}
    {
      "id": "{{$f.ID}}",
      "title": "{{$f.Title}}",
      "type": "{{$f.Type}}",
      "severity": "{{$f.Severity}}",
      "cvss_score": {{$f.CVSSScore}},
      "url": "{{$f.URL}}",
      "description": "{{$f.Description}}",
      "remediation": "{{$f.Remediation}}"
    }
    {{end}}
  ]
}
`
