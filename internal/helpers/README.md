# Helpers Package

Package `helpers` provides reusable utility functions for the Laravel MCP Companion application.

## Path Helpers

### Overview

The path helpers provide OS-specific cache directory paths following best practices for each platform:

- **Linux**: `$XDG_CACHE_HOME/laravel-mcp-companion/` or `$HOME/.cache/laravel-mcp-companion/`
- **macOS**: `$HOME/Library/Caches/laravel-mcp-companion/`
- **Windows**: `%LocalAppData%/laravel-mcp-companion/`

### Functions

#### `GetDefaultDocsPath() (string, error)`

Returns the default path for documentation storage using OS-specific cache directories.

**Returns:**
- `string`: Full path to the docs directory
- `error`: Error if the cache directory cannot be determined

**Example:**
```go
docsPath, err := helpers.GetDefaultDocsPath()
if err != nil {
    log.Printf("Could not determine cache directory: %v", err)
    docsPath = "./docs" // fallback
}
```

**Result Examples:**
- Linux: `/home/user/.cache/laravel-mcp-companion/docs`
- macOS: `/Users/user/Library/Caches/laravel-mcp-companion/docs`
- Windows: `C:\Users\user\AppData\Local\laravel-mcp-companion\docs`

---

#### `GetDefaultExternalCachePath() (string, error)`

Returns the default path for external service cache storage.

**Returns:**
- `string`: Full path to the external cache directory
- `error`: Error if the cache directory cannot be determined

**Example:**
```go
cachePath, err := helpers.GetDefaultExternalCachePath()
if err != nil {
    log.Printf("Could not determine cache directory: %v", err)
    cachePath = "./cache/external" // fallback
}
```

**Result Examples:**
- Linux: `/home/user/.cache/laravel-mcp-companion/cache/external`
- macOS: `/Users/user/Library/Caches/laravel-mcp-companion/cache/external`
- Windows: `C:\Users\user\AppData\Local\laravel-mcp-companion\cache\external`

---

#### `EnsureDirExists(path string) error`

Ensures that a directory exists, creating it and all parent directories if necessary (equivalent to `mkdir -p`).

**Parameters:**
- `path`: The directory path to ensure exists

**Returns:**
- `error`: Error if the directory cannot be created

**Example:**
```go
if err := helpers.EnsureDirExists("/path/to/nested/directory"); err != nil {
    log.Fatalf("Failed to create directory: %v", err)
}
```

**Features:**
- Creates all parent directories as needed
- Idempotent (safe to call multiple times)
- Sets directory permissions to `0755`

---

## Usage in Main Application

The helpers are used in `cmd/server/main.go` to provide smart defaults:

```go
// Get default paths from cache directory
defaultDocsPath, err := helpers.GetDefaultDocsPath()
if err != nil {
    // Fallback to ./docs if cache dir cannot be determined
    defaultDocsPath = "./docs"
    logging.Warn("Could not determine cache directory, using ./docs as fallback: %v", err)
}

// Parse flags with smart defaults
docsPath := flag.String("docs-path", defaultDocsPath, "Path to documentation directory (default: OS-specific cache dir)")
flag.Parse()

// Ensure directory exists
if err := helpers.EnsureDirExists(*docsPath); err != nil {
    logging.Error("Failed to create docs directory: %v", err)
    os.Exit(1)
}
```

## Benefits

1. **No Configuration Required**: Works out-of-the-box without manual setup
2. **OS-Specific Best Practices**: Uses the correct cache location for each operating system
3. **Fallback Support**: Gracefully degrades to local directory if cache dir unavailable
4. **User Override**: Users can still use `--docs-path` flag to specify custom location
5. **Reusable**: Easy to extend for other cache/storage needs in the future

## Testing

The package includes comprehensive tests:

```bash
go test ./internal/helpers/... -v
```

Tests cover:
- Path generation correctness
- Directory structure validation
- Absolute path verification
- Directory creation (including nested directories)
- Idempotent behavior

---

## Constants

### `AppName`

```go
const AppName = "laravel-mcp-companion"
```

The application name used in cache directory paths. This ensures the application has its own isolated cache space.
