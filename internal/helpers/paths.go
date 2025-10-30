package helpers

import (
	"os"
	"path/filepath"
)

const (
	// AppName is the application name used for cache directories
	AppName = "laravel-mcp-companion"
)

// GetDefaultDocsPath returns the default path for documentation storage.
// It uses os.UserCacheDir() to get the OS-specific cache directory and
// creates an application-specific subdirectory within it.
//
// Returns:
//   - string: The full path to the docs directory
//   - error: An error if the cache directory cannot be determined
//
// The resulting path will be:
//   - Linux: $XDG_CACHE_HOME/laravel-mcp-companion/docs or $HOME/.cache/laravel-mcp-companion/docs
//   - macOS: $HOME/Library/Caches/laravel-mcp-companion/docs
//   - Windows: %LocalAppData%/laravel-mcp-companion/docs
func GetDefaultDocsPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	// Create application-specific subdirectory path
	docsPath := filepath.Join(cacheDir, AppName, "docs")
	return docsPath, nil
}

// GetDefaultExternalCachePath returns the default path for external service cache storage.
// It uses the same base cache directory as GetDefaultDocsPath().
//
// Returns:
//   - string: The full path to the external cache directory
//   - error: An error if the cache directory cannot be determined
func GetDefaultExternalCachePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	// Create application-specific subdirectory path for external cache
	cachePath := filepath.Join(cacheDir, AppName, "cache", "external")
	return cachePath, nil
}

// EnsureDirExists ensures that a directory exists, creating it if necessary.
// It creates all parent directories as needed (like mkdir -p).
//
// Parameters:
//   - path: The directory path to ensure exists
//
// Returns:
//   - error: An error if the directory cannot be created
func EnsureDirExists(path string) error {
	return os.MkdirAll(path, 0755)
}
