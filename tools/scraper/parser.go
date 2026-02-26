package main

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// EndpointDef represents a discovered CUPI REST endpoint.
type EndpointDef struct {
	Chapter     string   `json:"chapter"`
	ChapterURL  string   `json:"chapterUrl"`
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	QueryParams []string `json:"queryParams"`
}

// endpointRegex matches HTTP methods followed by /vmrest paths in text.
var endpointRegex = regexp.MustCompile(`(?i)\b(GET|POST|PUT|DELETE)\b[^/]*/vmrest(/[^\s"<>{}\[\]]+)`)

// parseEndpoints extracts endpoints from an HTML page body.
func parseEndpoints(chapterName, chapterURL, htmlBody string) []EndpointDef {
	var endpoints []EndpointDef
	seen := make(map[string]bool)

	// Method 1: regex scan over raw HTML text
	for _, match := range endpointRegex.FindAllStringSubmatch(htmlBody, -1) {
		method := strings.ToUpper(match[1])
		path := normalizePath(match[2])
		key := method + ":" + path
		if seen[key] {
			continue
		}
		seen[key] = true
		endpoints = append(endpoints, EndpointDef{
			Chapter:    chapterName,
			ChapterURL: chapterURL,
			Method:     method,
			Path:       "/vmrest" + path,
		})
	}

	// Method 2: HTML table parsing for description context
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err == nil {
		tableEndpoints := extractFromTables(doc, chapterName, chapterURL, seen)
		endpoints = append(endpoints, tableEndpoints...)
	}

	return endpoints
}

// extractFromTables finds <table> elements and extracts method+path+description triples.
func extractFromTables(n *html.Node, chapterName, chapterURL string, seen map[string]bool) []EndpointDef {
	var results []EndpointDef

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			cells := collectTDText(node)
			if len(cells) >= 2 {
				for _, match := range endpointRegex.FindAllStringSubmatch(cells[0], -1) {
					method := strings.ToUpper(match[1])
					path := normalizePath(match[2])
					key := method + ":" + path
					if seen[key] {
						continue
					}
					seen[key] = true
					desc := ""
					if len(cells) >= 2 {
						desc = strings.TrimSpace(cells[1])
					}
					results = append(results, EndpointDef{
						Chapter:     chapterName,
						ChapterURL:  chapterURL,
						Method:      method,
						Path:        "/vmrest" + path,
						Description: desc,
					})
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return results
}

// collectTDText returns text content of <td> cells in a <tr> node.
func collectTDText(tr *html.Node) []string {
	var cells []string
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			cells = append(cells, textContent(c))
		}
	}
	return cells
}

// textContent returns the concatenated text content of an HTML node.
func textContent(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(sb.String())
}

// normalizePath normalizes a CUPI path: lowercases {objectId} placeholders.
func normalizePath(path string) string {
	// Remove trailing punctuation
	path = strings.TrimRight(path, ".,;:")
	// Normalize <objectid> → {objectId}
	path = regexp.MustCompile(`(?i)<objectid>`).ReplaceAllString(path, "{objectId}")
	// Normalize {objectid} → {objectId}
	path = regexp.MustCompile(`\{objectid\}`).ReplaceAllStringFunc(path, func(s string) string {
		return "{objectId}"
	})
	return path
}
