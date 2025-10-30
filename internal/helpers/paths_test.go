package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetDefaultDocsPath(t *testing.T) {
	path, err := GetDefaultDocsPath()
	if err != nil {
		t.Fatalf("GetDefaultDocsPath() returned error: %v", err)
	}

	// Path should not be empty
	if path == "" {
		t.Error("GetDefaultDocsPath() returned empty path")
	}

	// Path should contain the app name
	if !strings.Contains(path, AppName) {
		t.Errorf("GetDefaultDocsPath() path %q does not contain app name %q", path, AppName)
	}

	// Path should end with "docs"
	if !strings.HasSuffix(path, "docs") {
		t.Errorf("GetDefaultDocsPath() path %q does not end with 'docs'", path)
	}

	// Path should be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("GetDefaultDocsPath() returned relative path: %q", path)
	}
}

func TestGetDefaultExternalCachePath(t *testing.T) {
	path, err := GetDefaultExternalCachePath()
	if err != nil {
		t.Fatalf("GetDefaultExternalCachePath() returned error: %v", err)
	}

	// Path should not be empty
	if path == "" {
		t.Error("GetDefaultExternalCachePath() returned empty path")
	}

	// Path should contain the app name
	if !strings.Contains(path, AppName) {
		t.Errorf("GetDefaultExternalCachePath() path %q does not contain app name %q", path, AppName)
	}

	// Path should contain "cache/external"
	if !strings.Contains(path, filepath.Join("cache", "external")) {
		t.Errorf("GetDefaultExternalCachePath() path %q does not contain 'cache/external'", path)
	}

	// Path should be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("GetDefaultExternalCachePath() returned relative path: %q", path)
	}
}

func TestEnsureDirExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test", "nested", "directories")

	// Ensure the directory doesn't exist yet
	if _, err := os.Stat(testPath); !os.IsNotExist(err) {
		t.Fatalf("Test path %q should not exist yet", testPath)
	}

	// Create the directory
	err := EnsureDirExists(testPath)
	if err != nil {
		t.Fatalf("EnsureDirExists() returned error: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Path %q is not a directory", testPath)
	}

	// Calling again should not error (idempotent)
	err = EnsureDirExists(testPath)
	if err != nil {
		t.Errorf("EnsureDirExists() should be idempotent, but returned error: %v", err)
	}
}

func TestGetDefaultDocsPath_Structure(t *testing.T) {
	path, err := GetDefaultDocsPath()
	if err != nil {
		t.Fatalf("GetDefaultDocsPath() returned error: %v", err)
	}

	// Verify expected structure: <cache-dir>/<app-name>/docs
	expectedSuffix := filepath.Join(AppName, "docs")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("GetDefaultDocsPath() path %q does not have expected suffix %q", path, expectedSuffix)
	}
}

func TestGetDefaultExternalCachePath_Structure(t *testing.T) {
	path, err := GetDefaultExternalCachePath()
	if err != nil {
		t.Fatalf("GetDefaultExternalCachePath() returned error: %v", err)
	}

	// Verify expected structure: <cache-dir>/<app-name>/cache/external
	expectedSuffix := filepath.Join(AppName, "cache", "external")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("GetDefaultExternalCachePath() path %q does not have expected suffix %q", path, expectedSuffix)
	}
}
