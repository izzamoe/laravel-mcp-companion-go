package docs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/docs"
)

func TestManager_ListDocs(t *testing.T) {
	// Create temporary docs directory
	tmpDir := t.TempDir()
	versionDir := filepath.Join(tmpDir, "12.x")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	testFiles := []string{"routing.md", "middleware.md", "controllers.md"}
	for _, file := range testFiles {
		path := filepath.Join(versionDir, file)
		if err := os.WriteFile(path, []byte("# Test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create manager
	manager := docs.NewManager(tmpDir, "12.x")

	// List docs
	files, err := manager.ListDocs("12.x")
	if err != nil {
		t.Fatalf("ListDocs failed: %v", err)
	}

	// Verify
	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(files))
	}
}

func TestManager_ReadDoc(t *testing.T) {
	// Create temporary docs directory
	tmpDir := t.TempDir()
	versionDir := filepath.Join(tmpDir, "12.x")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test file
	testContent := "# Laravel Routing\n\nThis is routing documentation."
	testFile := filepath.Join(versionDir, "routing.md")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create manager
	manager := docs.NewManager(tmpDir, "12.x")

	// Read doc
	content, err := manager.ReadDoc("12.x", "routing.md")
	if err != nil {
		t.Fatalf("ReadDoc failed: %v", err)
	}

	// Verify
	if content != testContent {
		t.Errorf("Content mismatch.\nExpected: %s\nGot: %s", testContent, content)
	}
}

func TestManager_SearchDocs(t *testing.T) {
	// Create temporary docs directory
	tmpDir := t.TempDir()
	versionDir := filepath.Join(tmpDir, "12.x")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files with searchable content
	files := map[string]string{
		"routing.md":    "# Routing\n\nLearn about HTTP routing in Laravel.",
		"database.md":   "# Database\n\nQuery builder and Eloquent ORM.",
		"middleware.md": "# Middleware\n\nHTTP middleware for request filtering.",
	}

	for name, content := range files {
		path := filepath.Join(versionDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create manager
	manager := docs.NewManager(tmpDir, "12.x")

	// Search for "HTTP"
	results, err := manager.SearchDocs("HTTP", "12.x")
	if err != nil {
		t.Fatalf("SearchDocs failed: %v", err)
	}

	// Should return results string
	if len(results) == 0 {
		t.Error("Expected search results, got empty string")
	}

	// Should mention HTTP
	if !contains(results, "HTTP") && !contains(results, "http") {
		t.Error("Expected results to mention HTTP")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || s[0:len(substr)] == substr ||
			len(s) > len(substr) && findInString(s, substr))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
