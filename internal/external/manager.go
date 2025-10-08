package external

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	cacheValidDuration = 24 * time.Hour // Cache valid for 24 hours
)

// ServiceConfig holds configuration for an external Laravel service
type ServiceConfig struct {
	Name        string
	URL         string
	Description string
}

// ServiceMetadata holds cache metadata for a service
type ServiceMetadata struct {
	ServiceName string    `json:"service_name"`
	URL         string    `json:"url"`
	CachedAt    time.Time `json:"cached_at"`
	ContentSize int       `json:"content_size"`
}

// ExternalManager manages external Laravel service documentation
type ExternalManager struct {
	scraper   *WebScraper
	cachePath string
	services  map[string]ServiceConfig
}

// Service URL mappings for external Laravel services
var serviceURLs = map[string]ServiceConfig{
	"forge": {
		Name:        "Laravel Forge",
		URL:         "https://forge.laravel.com/docs",
		Description: "Server management and deployment platform",
	},
	"vapor": {
		Name:        "Laravel Vapor",
		URL:         "https://docs.vapor.build",
		Description: "Serverless deployment platform",
	},
	"envoyer": {
		Name:        "Laravel Envoyer",
		URL:         "https://envoyer.io/docs",
		Description: "Zero-downtime deployment platform",
	},
	"nova": {
		Name:        "Laravel Nova",
		URL:         "https://nova.laravel.com/docs",
		Description: "Administration panel for Laravel",
	},
}

// NewExternalManager creates a new external documentation manager
func NewExternalManager(cachePath string) *ExternalManager {
	// Ensure cache directory exists
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		// If we can't create cache dir, continue with empty cache
		cachePath = ""
	}

	return &ExternalManager{
		scraper:   NewWebScraper(),
		cachePath: cachePath,
		services:  serviceURLs,
	}
}

// UpdateService updates documentation for a specific service
func (m *ExternalManager) UpdateService(serviceName string, force bool) (string, error) {
	// Validate service name
	config, exists := m.services[serviceName]
	if !exists {
		return "", fmt.Errorf("unknown service: %s. Available services: forge, vapor, envoyer, nova", serviceName)
	}

	// Check if cache exists and is valid (unless forced)
	if !force {
		if valid, _ := m.isCacheValid(serviceName); valid {
			return fmt.Sprintf("Service '%s' documentation is already up to date (cached less than 24h ago). Use force=true to update anyway.", config.Name), nil
		}
	}

	// Fetch content from service URL
	content, err := m.scraper.FetchResource(config.URL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s documentation: %w", config.Name, err)
	}

	// Save to cache
	if err := m.saveToCache(serviceName, content); err != nil {
		return "", fmt.Errorf("failed to cache %s documentation: %w", config.Name, err)
	}

	// Save metadata
	metadata := ServiceMetadata{
		ServiceName: serviceName,
		URL:         config.URL,
		CachedAt:    time.Now(),
		ContentSize: len(content),
	}
	if err := m.saveMetadata(serviceName, metadata); err != nil {
		// Non-fatal, just log
		fmt.Fprintf(os.Stderr, "Warning: failed to save metadata for %s: %v\n", serviceName, err)
	}

	return fmt.Sprintf("Successfully updated %s documentation\n- URL: %s\n- Content Size: %d bytes\n- Cached at: %s",
		config.Name, config.URL, len(content), time.Now().Format(time.RFC3339)), nil
}

// UpdateServices updates documentation for multiple services
func (m *ExternalManager) UpdateServices(serviceNames []string, force bool) (string, error) {
	// If no services specified, update all
	if len(serviceNames) == 0 {
		serviceNames = []string{"forge", "vapor", "envoyer", "nova"}
	}

	var results []string
	var errors []string

	for _, serviceName := range serviceNames {
		_, err := m.UpdateService(serviceName, force)
		if err != nil {
			errors = append(errors, fmt.Sprintf("- %s: %v", serviceName, err))
		} else {
			results = append(results, fmt.Sprintf("✓ %s", serviceName))
		}
	}

	// Build response
	var response strings.Builder
	response.WriteString("# External Service Documentation Update\n\n")

	if len(results) > 0 {
		response.WriteString("## Successfully Updated:\n\n")
		for _, r := range results {
			response.WriteString(r + "\n")
		}
		response.WriteString("\n")
	}

	if len(errors) > 0 {
		response.WriteString("## Errors:\n\n")
		for _, e := range errors {
			response.WriteString(e + "\n")
		}
		response.WriteString("\n")
	}

	response.WriteString(fmt.Sprintf("**Total:** %d updated, %d errors\n", len(results), len(errors)))

	if len(errors) > 0 && len(results) == 0 {
		return "", fmt.Errorf("all updates failed:\n%s", strings.Join(errors, "\n"))
	}

	return response.String(), nil
}

// SearchServices searches through cached external service documentation
func (m *ExternalManager) SearchServices(query string, serviceNames []string) (string, error) {
	// If no services specified, search all
	if len(serviceNames) == 0 {
		serviceNames = []string{"forge", "vapor", "envoyer", "nova"}
	}

	query = strings.ToLower(query)
	var results []string
	totalMatches := 0

	for _, serviceName := range serviceNames {
		// Validate service
		config, exists := m.services[serviceName]
		if !exists {
			continue
		}

		// Get cached content
		content, err := m.getCachedContent(serviceName)
		if err != nil {
			// Service not cached, skip
			continue
		}

		// Search for query in content
		matches := m.searchInContent(content, query)
		if matches > 0 {
			results = append(results, fmt.Sprintf("**%s:** %d matches", config.Name, matches))
			totalMatches += matches
		}
	}

	// Build response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("# Search Results for '%s'\n\n", query))

	if len(results) == 0 {
		response.WriteString("No matches found in cached external service documentation.\n\n")
		response.WriteString("**Tip:** Try updating the service documentation first using `update_external_laravel_docs`.\n")
		return response.String(), nil
	}

	response.WriteString(fmt.Sprintf("Found **%d total matches** across %d services:\n\n", totalMatches, len(results)))
	for _, result := range results {
		response.WriteString("- " + result + "\n")
	}

	return response.String(), nil
}

// SearchServicesWithContext searches and returns matching text with context
func (m *ExternalManager) SearchServicesWithContext(query string, serviceNames []string, contextLength int) (string, error) {
	// If no services specified, search all
	if len(serviceNames) == 0 {
		serviceNames = []string{"forge", "vapor", "envoyer", "nova"}
	}

	if contextLength <= 0 {
		contextLength = 200
	}

	query = strings.ToLower(query)
	var results []string
	totalMatches := 0

	for _, serviceName := range serviceNames {
		// Validate service
		config, exists := m.services[serviceName]
		if !exists {
			continue
		}

		// Get cached content
		content, err := m.getCachedContent(serviceName)
		if err != nil {
			// Service not cached, skip
			continue
		}

		// Find matches with context
		contexts := m.findContexts(content, query, contextLength)
		if len(contexts) > 0 {
			results = append(results, fmt.Sprintf("\n## %s (%d matches)\n\n%s",
				config.Name, len(contexts), strings.Join(contexts, "\n\n---\n\n")))
			totalMatches += len(contexts)
		}
	}

	// Build response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("# Search Results with Context for '%s'\n\n", query))

	if len(results) == 0 {
		response.WriteString("No matches found in cached external service documentation.\n\n")
		response.WriteString("**Tip:** Try updating the service documentation first using `update_external_laravel_docs`.\n")
		return response.String(), nil
	}

	response.WriteString(fmt.Sprintf("Found **%d total matches** across services:\n", totalMatches))
	for _, result := range results {
		response.WriteString(result)
	}

	return response.String(), nil
}

// GetCachedServices returns list of services with cached documentation
func (m *ExternalManager) GetCachedServices() []string {
	var cached []string

	for serviceName := range m.services {
		if valid, _ := m.isCacheValid(serviceName); valid {
			cached = append(cached, serviceName)
		}
	}

	return cached
}

// --- Private helper methods ---

// isCacheValid checks if cache exists and is still valid
func (m *ExternalManager) isCacheValid(serviceName string) (bool, error) {
	if m.cachePath == "" {
		return false, nil
	}

	metadataPath := filepath.Join(m.cachePath, fmt.Sprintf("%s_metadata.json", serviceName))
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return false, err
	}

	var metadata ServiceMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return false, err
	}

	// Check if cache is still valid
	age := time.Since(metadata.CachedAt)
	return age < cacheValidDuration, nil
}

// saveToCache saves content to cache file
func (m *ExternalManager) saveToCache(serviceName, content string) error {
	if m.cachePath == "" {
		return fmt.Errorf("cache path not configured")
	}

	cachePath := filepath.Join(m.cachePath, fmt.Sprintf("%s_docs.txt", serviceName))
	return os.WriteFile(cachePath, []byte(content), 0644)
}

// saveMetadata saves service metadata
func (m *ExternalManager) saveMetadata(serviceName string, metadata ServiceMetadata) error {
	if m.cachePath == "" {
		return fmt.Errorf("cache path not configured")
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	metadataPath := filepath.Join(m.cachePath, fmt.Sprintf("%s_metadata.json", serviceName))
	return os.WriteFile(metadataPath, data, 0644)
}

// getCachedContent reads cached content for a service
func (m *ExternalManager) getCachedContent(serviceName string) (string, error) {
	if m.cachePath == "" {
		return "", fmt.Errorf("cache path not configured")
	}

	cachePath := filepath.Join(m.cachePath, fmt.Sprintf("%s_docs.txt", serviceName))
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return "", fmt.Errorf("service documentation not cached: %w", err)
	}

	return string(data), nil
}

// searchInContent counts matches of query in content
func (m *ExternalManager) searchInContent(content, query string) int {
	contentLower := strings.ToLower(content)
	return strings.Count(contentLower, query)
}

// findContexts finds all occurrences of query with surrounding context
func (m *ExternalManager) findContexts(content, query string, contextLength int) []string {
	contentLower := strings.ToLower(content)
	var contexts []string
	maxContexts := 10 // Limit to prevent too many results

	searchPos := 0
	for i := 0; i < maxContexts; i++ {
		idx := strings.Index(contentLower[searchPos:], query)
		if idx == -1 {
			break
		}

		// Adjust index to absolute position
		idx += searchPos

		// Calculate context boundaries
		start := idx - contextLength
		if start < 0 {
			start = 0
		}
		end := idx + len(query) + contextLength
		if end > len(content) {
			end = len(content)
		}

		// Extract context
		context := content[start:end]

		// Add ellipsis if truncated
		if start > 0 {
			context = "..." + context
		}
		if end < len(content) {
			context = context + "..."
		}

		// Highlight the match (using markdown bold)
		highlightedContext := strings.ReplaceAll(
			context,
			content[idx:idx+len(query)],
			"**"+content[idx:idx+len(query)]+"**",
		)

		contexts = append(contexts, highlightedContext)

		// Move search position forward
		searchPos = idx + len(query)
	}

	return contexts
}
