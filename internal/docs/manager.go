package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/models"
)

// Manager handles documentation operations
type Manager struct {
	DocsPath       string
	defaultVersion string
	versions       []string
	cache          *Cache
	mu             sync.RWMutex
}

// NewManager creates a new documentation manager
func NewManager(docsPath, defaultVersion string) *Manager {
	return &Manager{
		DocsPath:       docsPath,
		defaultVersion: defaultVersion,
		versions:       models.SupportedVersions,
		cache:          NewCache(),
	}
}

// ListDocs returns list of available documentation files
func (m *Manager) ListDocs(version string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if version == "" {
		version = m.defaultVersion
	}

	versionPath := filepath.Join(m.DocsPath, version)

	// Check if version directory exists
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("version %s not found", version)
	}

	// Read directory
	entries, err := os.ReadDir(versionPath)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// ReadDoc reads a documentation file
func (m *Manager) ReadDoc(version, filename string) (string, error) {
	if version == "" {
		version = m.defaultVersion
	}

	// Build full path
	fullPath := filepath.Join(m.DocsPath, version, filename)

	// Security check - prevent directory traversal
	if !isPathSafe(m.DocsPath, fullPath) {
		return "", fmt.Errorf("invalid path: directory traversal detected")
	}

	// Check cache
	cacheKey := fmt.Sprintf("%s:%s", version, filename)
	if content, found := m.cache.Get(cacheKey); found {
		return content, nil
	}

	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("document not found: %s", filename)
		}
		return "", fmt.Errorf("read file: %w", err)
	}

	content := string(data)

	// Cache content
	m.cache.Set(cacheKey, content)

	return content, nil
}

// SearchDocs searches across documentation files
func (m *Manager) SearchDocs(query, version string) (string, error) {
	if version == "" {
		version = m.defaultVersion
	}

	// Check cache
	cacheKey := fmt.Sprintf("search:%s:%s", version, query)
	if results, found := m.cache.GetSearch(cacheKey); found {
		return results, nil
	}

	versionPath := filepath.Join(m.DocsPath, version)

	// Check if version exists
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return "", fmt.Errorf("version %s not found", version)
	}

	// Read all files
	entries, err := os.ReadDir(versionPath)
	if err != nil {
		return "", fmt.Errorf("read directory: %w", err)
	}

	// Search results
	type searchResult struct {
		filename string
		matches  int
	}
	var results []searchResult

	queryLower := strings.ToLower(query)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(versionPath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		content := strings.ToLower(string(data))
		matches := strings.Count(content, queryLower)

		if matches > 0 {
			results = append(results, searchResult{
				filename: entry.Name(),
				matches:  matches,
			})
		}
	}

	// Format results
	var output strings.Builder
	output.WriteString(fmt.Sprintf("# Search Results for '%s' in Laravel %s\n\n", query, version))

	if len(results) == 0 {
		output.WriteString("No matches found.\n")
	} else {
		output.WriteString(fmt.Sprintf("Found %d files with matches:\n\n", len(results)))
		for _, result := range results {
			output.WriteString(fmt.Sprintf("- **%s**: %d matches\n", result.filename, result.matches))
		}
	}

	resultStr := output.String()

	// Cache results
	m.cache.SetSearch(cacheKey, resultStr)

	return resultStr, nil
}

// SearchWithContext searches and returns results with surrounding context
func (m *Manager) SearchWithContext(query, version string, contextLength int) (string, error) {
	if version == "" {
		version = m.defaultVersion
	}

	versionPath := filepath.Join(m.DocsPath, version)
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return "", fmt.Errorf("version %s not found", version)
	}

	entries, err := os.ReadDir(versionPath)
	if err != nil {
		return "", fmt.Errorf("read directory: %w", err)
	}

	type contextMatch struct {
		filename string
		context  string
		position int
	}
	var matches []contextMatch

	queryLower := strings.ToLower(query)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(versionPath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		content := string(data)
		contentLower := strings.ToLower(content)

		// Find all occurrences
		index := 0
		for {
			pos := strings.Index(contentLower[index:], queryLower)
			if pos == -1 {
				break
			}

			actualPos := index + pos

			// Extract context
			start := actualPos - contextLength
			if start < 0 {
				start = 0
			}
			end := actualPos + len(query) + contextLength
			if end > len(content) {
				end = len(content)
			}

			contextStr := content[start:end]
			// Add ellipsis if truncated
			if start > 0 {
				contextStr = "..." + contextStr
			}
			if end < len(content) {
				contextStr = contextStr + "..."
			}

			matches = append(matches, contextMatch{
				filename: entry.Name(),
				context:  contextStr,
				position: actualPos,
			})

			index = actualPos + len(query)
			if index >= len(content) {
				break
			}
		}
	}

	// Format results
	var output strings.Builder
	output.WriteString(fmt.Sprintf("# Search Results with Context for '%s' in Laravel %s\n\n", query, version))

	if len(matches) == 0 {
		output.WriteString("No matches found.\n")
	} else {
		output.WriteString(fmt.Sprintf("Found %d matches:\n\n", len(matches)))

		currentFile := ""
		for _, match := range matches {
			if match.filename != currentFile {
				output.WriteString(fmt.Sprintf("\n### %s\n\n", match.filename))
				currentFile = match.filename
			}
			output.WriteString(fmt.Sprintf("```\n%s\n```\n\n", match.context))
		}
	}

	return output.String(), nil
}

// GetStructure extracts the table of contents from a documentation file
func (m *Manager) GetStructure(filename, version string) (string, error) {
	if version == "" {
		version = m.defaultVersion
	}

	content, err := m.ReadDoc(version, filename)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("# Structure of %s (Laravel %s)\n\n", filename, version))

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Extract headers (markdown style)
		if strings.HasPrefix(trimmed, "#") {
			// Count header level
			level := 0
			for i := 0; i < len(trimmed) && trimmed[i] == '#'; i++ {
				level++
			}

			if level > 0 && level <= 6 {
				headerText := strings.TrimSpace(trimmed[level:])
				indent := strings.Repeat("  ", level-1)
				output.WriteString(fmt.Sprintf("%s- %s\n", indent, headerText))
			}
		}
	}

	return output.String(), nil
}

// BrowseByCategory returns documentation files related to a specific category
func (m *Manager) BrowseByCategory(category, version string) (string, error) {
	if version == "" {
		version = m.defaultVersion
	}

	files, err := m.ListDocs(version)
	if err != nil {
		return "", err
	}

	// Category mappings
	categoryMap := map[string][]string{
		"frontend":       {"blade", "vite", "mix", "frontend", "views"},
		"database":       {"database", "migrations", "seeding", "eloquent", "queries"},
		"authentication": {"authentication", "authorization", "sanctum", "passport", "fortify"},
		"testing":        {"testing", "http-tests", "console-tests", "dusk", "database-testing"},
		"security":       {"authentication", "authorization", "csrf", "encryption", "hashing", "verification"},
		"deployment":     {"deployment", "octane", "horizon", "envoy", "sail", "homestead"},
		"packages":       {"packages", "cashier", "socialite", "scout", "telescope", "horizon"},
	}

	categoryLower := strings.ToLower(category)
	keywords, ok := categoryMap[categoryLower]
	if !ok {
		// Try fuzzy matching
		keywords = []string{categoryLower}
	}

	// Filter files by category
	var matches []string
	for _, file := range files {
		fileLower := strings.ToLower(file)
		for _, keyword := range keywords {
			if strings.Contains(fileLower, keyword) {
				matches = append(matches, file)
				break
			}
		}
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("# Documentation for Category: %s (Laravel %s)\n\n", category, version))

	if len(matches) == 0 {
		output.WriteString(fmt.Sprintf("No documentation files found for category '%s'.\n\n", category))
		output.WriteString("Available categories: frontend, database, authentication, testing, security, deployment, packages\n")
	} else {
		output.WriteString(fmt.Sprintf("Found %d files:\n\n", len(matches)))
		for _, file := range matches {
			output.WriteString(fmt.Sprintf("- %s\n", file))
		}
	}

	return output.String(), nil
}

// GetInfo returns metadata about documentation versions
func (m *Manager) GetInfo(version string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var output strings.Builder

	if version == "" {
		// Show all versions
		output.WriteString("# Laravel Documentation Information\n\n")

		for _, ver := range m.versions {
			versionPath := filepath.Join(m.DocsPath, ver)
			info, err := os.Stat(versionPath)
			if err != nil {
				continue
			}

			files, err := m.ListDocs(ver)
			if err != nil {
				continue
			}

			output.WriteString(fmt.Sprintf("## Version %s\n", ver))
			output.WriteString(fmt.Sprintf("Last updated: %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
			output.WriteString(fmt.Sprintf("Files: %d documentation files\n\n", len(files)))
		}
	} else {
		// Show specific version
		versionPath := filepath.Join(m.DocsPath, version)
		info, err := os.Stat(versionPath)
		if err != nil {
			return "", fmt.Errorf("version %s not found", version)
		}

		files, err := m.ListDocs(version)
		if err != nil {
			return "", err
		}

		output.WriteString(fmt.Sprintf("# Laravel %s Documentation\n\n", version))
		output.WriteString(fmt.Sprintf("Last updated: %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
		output.WriteString(fmt.Sprintf("Files: %d documentation files\n", len(files)))
	}

	return output.String(), nil
}

// ClearCache clears all cached documentation
func (m *Manager) ClearCache() {
	m.cache.Clear()
}

// isPathSafe checks if the path is within the base directory
func isPathSafe(base, path string) bool {
	// Get absolute paths
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if path starts with base
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return false
	}

	// Path should not contain ".."
	return !strings.HasPrefix(rel, "..")
}
