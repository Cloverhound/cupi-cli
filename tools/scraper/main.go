package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	jsonOut := flag.String("json", "api_reference.json", "Output JSON file path")
	mdOut := flag.String("md", "api_reference.md", "Output Markdown file path")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Fetching %d CUPI API documentation chapters...\n", len(chapters))

	var allEndpoints []EndpointDef

	httpClient := &http.Client{Timeout: 30 * time.Second}

	for i, chapter := range chapters {
		url := baseURL + chapter.File
		fmt.Fprintf(os.Stderr, "[%d/%d] Fetching: %s\n", i+1, len(chapters), chapter.Name)

		body, err := fetchPage(httpClient, url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  WARNING: failed to fetch %s: %v\n", url, err)
			continue
		}

		endpoints := parseEndpoints(chapter.Name, url, body)
		fmt.Fprintf(os.Stderr, "  Found %d endpoints\n", len(endpoints))
		allEndpoints = append(allEndpoints, endpoints...)

		// Polite delay between requests
		if i < len(chapters)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Fprintf(os.Stderr, "\nTotal endpoints discovered: %d\n", len(allEndpoints))

	if *jsonOut != "" {
		if err := writeJSON(*jsonOut, allEndpoints); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
	}

	if *mdOut != "" {
		if err := writeMarkdown(*mdOut, allEndpoints); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
	}
}

func fetchPage(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "cupi-cli-scraper/1.0 (API documentation scraper)")
	req.Header.Set("Accept", "text/html")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
