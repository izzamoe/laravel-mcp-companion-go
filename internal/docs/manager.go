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
	DocsPath string
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
