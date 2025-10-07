package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/models"
)

const (
	githubAPIBase = "https://api.github.com"
	docsRepo      = "laravel/docs"
	timeout       = 30 * time.Second
)

// GitHubUpdater handles documentation updates from GitHub
type GitHubUpdater struct {
	httpClient *http.Client
	basePath   string
}

// NewGitHubUpdater creates a new GitHub updater
func NewGitHubUpdater(basePath string) *GitHubUpdater {
	return &GitHubUpdater{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		basePath: basePath,
	}
}

// UpdateDocs updates documentation for a specific version
func (u *GitHubUpdater) UpdateDocs(version string) (string, error) {
	// Verify version is supported
	supported := false
	for _, v := range models.SupportedVersions {
		if v == version {
			supported = true
			break
		}
	}
	if !supported {
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	// Fetch commit SHA
	commitSHA, err := u.getLatestCommitSHA(version)
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit: %w", err)
	}

	// Fetch tree of markdown files
	files, err := u.getMarkdownFiles(version, commitSHA)
	if err != nil {
		return "", fmt.Errorf("failed to get file list: %w", err)
	}

	// Create version directory
	versionPath := filepath.Join(u.basePath, version)
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create version directory: %w", err)
	}

	// Download each file
	downloadedCount := 0
	for _, file := range files {
		if err := u.downloadFile(version, file, versionPath); err != nil {
			return "", fmt.Errorf("failed to download %s: %w", file, err)
		}
		downloadedCount++
	}

	// Save metadata
	metadata := models.DocMetadata{
		Version:   version,
		CommitSHA: commitSHA,
		SyncTime:  time.Now(),
		FileCount: downloadedCount,
	}

	if err := u.saveMetadata(versionPath, metadata); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	return fmt.Sprintf("Successfully updated %s documentation: %d files downloaded (commit: %s)",
		version, downloadedCount, commitSHA[:7]), nil
}

// getLatestCommitSHA fetches the latest commit SHA for a branch
func (u *GitHubUpdater) getLatestCommitSHA(version string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/commits/%s", githubAPIBase, docsRepo, version)

	resp, err := u.doRequest(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var commit struct {
		SHA string `json:"sha"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return "", fmt.Errorf("failed to decode commit response: %w", err)
	}

	return commit.SHA, nil
}

// getMarkdownFiles fetches list of markdown files in the docs repo
func (u *GitHubUpdater) getMarkdownFiles(version, commitSHA string) ([]string, error) {
	url := fmt.Sprintf("%s/repos/%s/git/trees/%s", githubAPIBase, docsRepo, commitSHA)

	resp, err := u.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tree struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("failed to decode tree response: %w", err)
	}

	// Filter for .md files
	var files []string
	for _, item := range tree.Tree {
		if item.Type == "blob" && filepath.Ext(item.Path) == ".md" {
			files = append(files, item.Path)
		}
	}

	return files, nil
}

// downloadFile downloads a single markdown file
func (u *GitHubUpdater) downloadFile(version, filename, destPath string) error {
	// Raw content URL
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", docsRepo, version, filename)

	resp, err := u.doRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Write to file
	filePath := filepath.Join(destPath, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// doRequest performs HTTP request with GitHub API headers
func (u *GitHubUpdater) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Laravel-MCP-Companion")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

// saveMetadata saves update metadata to file
func (u *GitHubUpdater) saveMetadata(versionPath string, metadata models.DocMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metaPath := filepath.Join(versionPath, ".metadata.json")
	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}
