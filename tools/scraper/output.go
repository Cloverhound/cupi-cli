package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// writeJSON writes endpoints to a JSON file.
func writeJSON(path string, endpoints []EndpointDef) error {
	data, err := json.MarshalIndent(endpoints, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Wrote %d endpoints to %s\n", len(endpoints), path)
	return nil
}

// writeMarkdown writes endpoints to a Markdown file.
func writeMarkdown(path string, endpoints []EndpointDef) error {
	var sb strings.Builder
	sb.WriteString("# CUPI API Reference\n\n")
	sb.WriteString("Auto-generated from Cisco CUPI API documentation.\n\n")
	sb.WriteString("| Method | Path | Description | Chapter |\n")
	sb.WriteString("|--------|------|-------------|--------|\n")

	for _, ep := range endpoints {
		desc := strings.ReplaceAll(ep.Description, "|", "\\|")
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		sb.WriteString(fmt.Sprintf("| %s | `%s` | %s | [%s](%s) |\n",
			ep.Method, ep.Path, desc, ep.Chapter, ep.ChapterURL))
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write Markdown file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Wrote %s\n", path)
	return nil
}
