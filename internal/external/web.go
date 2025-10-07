package external

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxContentSize = 5 * 1024 * 1024 // 5MB max
	requestTimeout = 15 * time.Second
)

// WebScraper handles fetching external web resources
type WebScraper struct {
	httpClient *http.Client
}

// NewWebScraper creates a new web scraper
func NewWebScraper() *WebScraper {
	return &WebScraper{
		httpClient: &http.Client{
			Timeout: requestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// FetchResource fetches content from a URL
func (w *WebScraper) FetchResource(urlStr string) (string, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow HTTP/HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("only HTTP/HTTPS URLs are supported")
	}

	// Perform request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Laravel-MCP-Companion/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Limit content size
	limitedReader := io.LimitReader(resp.Body, maxContentSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check if we hit the size limit
	if len(content) >= maxContentSize {
		return "", fmt.Errorf("content too large (max: %d bytes)", maxContentSize)
	}

	// Convert to string and basic cleanup
	contentStr := string(content)
	contentType := resp.Header.Get("Content-Type")

	// If HTML, try to extract main content
	if strings.Contains(contentType, "text/html") {
		contentStr = w.extractMainContent(contentStr)
	}

	return contentStr, nil
}

// extractMainContent attempts to extract main content from HTML
func (w *WebScraper) extractMainContent(html string) string {
	// Simple heuristic: remove script and style tags
	result := html

	// Remove script tags
	result = removeTagsWithContent(result, "<script", "</script>")

	// Remove style tags
	result = removeTagsWithContent(result, "<style", "</style>")

	// Remove excessive whitespace
	lines := strings.Split(result, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// removeTagsWithContent removes HTML tags and their content
func removeTagsWithContent(text, startTag, endTag string) string {
	result := text
	for {
		startIdx := strings.Index(result, startTag)
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(result[startIdx:], endTag)
		if endIdx == -1 {
			break
		}

		// Remove from start tag to end of end tag
		result = result[:startIdx] + result[startIdx+endIdx+len(endTag):]
	}
	return result
}

// FormatResource formats fetched content for display
func (w *WebScraper) FormatResource(url, content string) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("# External Resource\n\n"))
	result.WriteString(fmt.Sprintf("**URL:** %s\n\n", url))
	result.WriteString(fmt.Sprintf("**Content Length:** %d bytes\n\n", len(content)))
	result.WriteString("---\n\n")
	result.WriteString(content)

	return result.String()
}
