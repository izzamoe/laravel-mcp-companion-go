# JSON Schema Tag Format Fix

## Issue
Server panic dengan error:
```
panic: AddTool: tool "list_laravel_docs": input schema: ForType(server.ListDocsInput): 
tag must not begin with 'WORD=': "description=Specific Laravel version..."
```

## Root Cause
Library baru `github.com/modelcontextprotocol/go-sdk` menggunakan `github.com/google/jsonschema-go` yang memiliki format struct tag berbeda dari library lama.

## Format yang Salah ❌
```go
type Input struct {
    Field string `jsonschema:"description=Some description"`
    Required string `jsonschema:"required,description=Required field"`
    WithDefault bool `jsonschema:"default=true,description=Has default"`
}
```

## Format yang Benar ✅
```go
type Input struct {
    Field string `jsonschema:"Some description"`
    Required string `jsonschema:"required,Required field"`
    WithDefault *bool `jsonschema:"Has default value (pointer for optional with default)"`
}
```

## Key Changes Made

### 1. Removed `description=` prefix
**Before:**
```go
`jsonschema:"description=Specific Laravel version to list"`
```

**After:**
```go
`jsonschema:"Specific Laravel version to list"`
```

### 2. Removed `default=` syntax
**Before:**
```go
IncludeExternal bool `jsonschema:"default=true,description=Include external services"`
ContextLength float64 `jsonschema:"default=200,description=Context length"`
```

**After:**
```go
IncludeExternal *bool `jsonschema:"Include external services"`
ContextLength *int `jsonschema:"Context length (default: 200)"`
```

### 3. Handle defaults in code instead of tags
```go
// Default value handling for pointer types
includeExternal := true
if input.IncludeExternal != nil {
    includeExternal = *input.IncludeExternal
}

contextLength := 200
if input.ContextLength != nil {
    contextLength = *input.ContextLength
}

force := false
if input.Force != nil {
    force = *input.Force
}
```

## Files Fixed

### 1. internal/server/doc_tools.go
- Fixed `ListDocsInput`
- Fixed `ReadDocInput`
- Fixed `SearchDocsInput` - changed `bool` to `*bool`
- Fixed `SearchWithContextInput` - changed `float64` to `*int` and `bool` to `*bool`
- Fixed `GetStructureInput`
- Fixed `BrowseCategoryInput`
- Updated handler code to handle pointer defaults

### 2. internal/server/package_tools.go
- Fixed `RecommendPackageInput`
- Fixed `PackageInfoInput`
- Fixed `PackageCategoryInput`
- Fixed `PackageFeaturesInput`

### 3. internal/server/external_tools.go
- Fixed `UpdateDocsInput` - changed `bool` to `*bool`
- Fixed `DocsInfoInput`
- Fixed `UpdateExternalInput` - changed `bool` to `*bool`
- Fixed `SearchExternalInput`
- Fixed `ServiceInfoInput`
- Updated handler code to handle pointer defaults

## Result

✅ Server starts successfully without panic
✅ All 16 tools registered correctly
✅ JSON schema generation works properly
✅ Default values handled in code (cleaner approach)

## Verification Output
```
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Starting Laravel MCP Companion...
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Initialized documentation manager (path: ./docs, default: 12.x)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Initialized package catalog (path: ./configs/packages.json)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Initialized updater, web scraper, and external manager
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Created MCP server
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Registered documentation tools (6 tools)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Registered package tools (4 tools)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Registered update and info tools (2 tools)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Registered external service tools (4 tools)
[Laravel MCP] 2025/10/22 22:38:59 [INFO] Server ready with 16 total tools, starting event loop...
```

## Benefits of New Approach

1. **Cleaner struct tags** - No confusing `key=value` syntax
2. **Type safety** - Pointers make optional-with-defaults explicit
3. **Runtime flexibility** - Can distinguish between "not provided" and "provided with false/0"
4. **Better documentation** - Description text is clearer without prefixes

## Migration Complete! 🎉

Server now runs successfully with the official `github.com/modelcontextprotocol/go-sdk` v1.0.0 library!
