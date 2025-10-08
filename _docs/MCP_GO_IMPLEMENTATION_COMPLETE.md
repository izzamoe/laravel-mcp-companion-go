# Laravel MCP Companion - Complete Go Implementation Guide

## Menggunakan Library `github.com/mark3labs/mcp-go`

Dokumen ini adalah **complete reference** untuk implementasi Go yang **persis sama** dengan Python MCP server. Semua 16 tools, parameter schema, dan response format harus identik.

---

## 📦 Dependencies

```bash
go get github.com/mark3labs/mcp-go@latest
```

---

## 🏗️ Complete Server Implementation

### Main Entry Point dengan `mcp-go`

```go
// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // Setup logging ke stderr (stdout digunakan MCP stdio)
    log.SetOutput(os.Stderr)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // Create MCP server
    s := server.NewMCPServer(
        "Laravel MCP Companion",
        "1.0.0",
        server.WithToolCapabilities(true), // Enable tool listing notifications
    )
    
    // Initialize managers
    docManager, err := NewDocManager("./docs", "12.x")
    if err != nil {
        log.Fatalf("Failed to create doc manager: %v", err)
    }
    
    pkgCatalog := NewPackageCatalog()
    
    // Register all 16 tools
    registerAllTools(s, docManager, pkgCatalog)
    
    // Start stdio server
    log.Println("Starting Laravel MCP Companion...")
    if err := server.ServeStdio(s); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

---

## 🛠️ All 16 Tools - Complete Implementation

### Tool 1: `list_laravel_docs`

**Python Signature:**
```python
def list_laravel_docs(version: Optional[str] = None) -> str
```

**Go Implementation:**
```go
func registerListLaravelDocs(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("list_laravel_docs",
        mcp.WithDescription("List all available Laravel documentation files across versions. Essential for discovering what documentation exists before diving into specific topics."),
        mcp.WithString("version",
            mcp.Description("Specific Laravel version to list (e.g., '12.x'). If not provided, lists all versions"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        version := mcp.ParseString(request, "version", "")
        
        result, err := dm.ListDocs(version)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to list docs: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

**Response Format (harus sama persis):**
```markdown
Available Laravel Documentation Versions:

## Version 12.x
Last updated: 2024-01-15 10:30:00
Commit: abc1234
Files: 95 documentation files

## Version 11.x
Last updated: 2024-01-10 08:15:00
Commit: def5678
Files: 90 documentation files
```

---

### Tool 2: `read_laravel_doc_content`

**Python Signature:**
```python
def read_laravel_doc_content(filename: str, version: Optional[str] = None) -> str
```

**Go Implementation:**
```go
func registerReadLaravelDocContent(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("read_laravel_doc_content",
        mcp.WithDescription("Reads the complete content of a specific Laravel documentation file. This is the primary tool for accessing actual documentation content."),
        mcp.WithString("filename",
            mcp.Required(),
            mcp.Description("Name of the file (e.g., 'mix.md', 'vite.md')"),
        ),
        mcp.WithString("version",
            mcp.Description("Laravel version (e.g., '12.x'). Defaults to latest"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        filename, err := request.RequireString("filename")
        if err != nil {
            return mcp.NewToolResultError("filename is required"), nil
        }
        
        version := mcp.ParseString(request, "version", "")
        
        content, err := dm.ReadDoc(filename, version)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to read doc: %v", err)), nil
        }
        
        return mcp.NewToolResultText(content), nil
    })
}
```

---

### Tool 3: `search_laravel_docs`

**Python Signature:**
```python
def search_laravel_docs(
    query: str, 
    version: Optional[str] = None, 
    include_external: bool = True
) -> str
```

**Go Implementation:**
```go
func registerSearchLaravelDocs(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("search_laravel_docs",
        mcp.WithDescription("Searches for specific terms across all Laravel documentation files. Returns file names and match counts."),
        mcp.WithString("query",
            mcp.Required(),
            mcp.Description("Search term to look for"),
        ),
        mcp.WithString("version",
            mcp.Description("Specific Laravel version to search (e.g., '12.x'). If not provided, searches all versions"),
        ),
        mcp.WithBoolean("include_external",
            mcp.DefaultBool(true),
            mcp.Description("Whether to include external Laravel services documentation in search"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        query, err := request.RequireString("query")
        if err != nil {
            return mcp.NewToolResultError("query is required"), nil
        }
        
        version := mcp.ParseString(request, "version", "")
        includeExternal := mcp.ParseBoolean(request, "include_external", true)
        
        result, err := dm.Search(query, version, includeExternal)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

**Response Format:**
```markdown
Search results for 'middleware':

**Core Laravel Documentation (15 files):**
  - 12.x/middleware.md (45 matches)
  - 12.x/routing.md (12 matches)
  - 11.x/middleware.md (43 matches)

**External Laravel Services (2 services):**
  - Forge: deployment.md (3 matches)
```

---

### Tool 4: `search_laravel_docs_with_context`

**Python Signature:**
```python
def search_laravel_docs_with_context(
    query: str,
    version: Optional[str] = None,
    context_length: int = 200,
    include_external: bool = True
) -> str
```

**Go Implementation:**
```go
func registerSearchLaravelDocsWithContext(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("search_laravel_docs_with_context",
        mcp.WithDescription("Advanced search that returns matching text with surrounding context. Shows exactly how terms are used in documentation."),
        mcp.WithString("query",
            mcp.Required(),
            mcp.Description("Search term"),
        ),
        mcp.WithString("version",
            mcp.Description("Specific version or None for all"),
        ),
        mcp.WithNumber("context_length",
            mcp.DefaultNumber(200),
            mcp.Description("Characters of context to show (default: 200)"),
        ),
        mcp.WithBoolean("include_external",
            mcp.DefaultBool(true),
            mcp.Description("Include external services"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        query, err := request.RequireString("query")
        if err != nil {
            return mcp.NewToolResultError("query is required"), nil
        }
        
        version := mcp.ParseString(request, "version", "")
        contextLength := int(mcp.ParseFloat64(request, "context_length", 200))
        includeExternal := mcp.ParseBoolean(request, "include_external", true)
        
        result, err := dm.SearchWithContext(query, version, contextLength, includeExternal)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 5: `get_doc_structure`

**Python Signature:**
```python
def get_doc_structure(filename: str, version: Optional[str] = None) -> str
```

**Go Implementation:**
```go
func registerGetDocStructure(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("get_doc_structure",
        mcp.WithDescription("Extracts the table of contents and structure from a documentation file. Shows headers and brief content previews."),
        mcp.WithString("filename",
            mcp.Required(),
            mcp.Description("Documentation file name"),
        ),
        mcp.WithString("version",
            mcp.Description("Laravel version"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        filename, err := request.RequireString("filename")
        if err != nil {
            return mcp.NewToolResultError("filename is required"), nil
        }
        
        version := mcp.ParseString(request, "version", "")
        
        structure, err := dm.GetStructure(filename, version)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to get structure: %v", err)), nil
        }
        
        return mcp.NewToolResultText(structure), nil
    })
}
```

---

### Tool 6: `browse_docs_by_category`

**Python Signature:**
```python
def browse_docs_by_category(category: str, version: Optional[str] = None) -> str
```

**Go Implementation:**
```go
func registerBrowseDocsByCategory(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("browse_docs_by_category",
        mcp.WithDescription("Discovers documentation files related to specific categories like 'frontend', 'database', or 'authentication'."),
        mcp.WithString("category",
            mcp.Required(),
            mcp.Description("Category like 'frontend', 'database', 'authentication', etc."),
        ),
        mcp.WithString("version",
            mcp.Description("Laravel version"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        category, err := request.RequireString("category")
        if err != nil {
            return mcp.NewToolResultError("category is required"), nil
        }
        
        version := mcp.ParseString(request, "version", "")
        
        result, err := dm.BrowseByCategory(category, version)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to browse: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 7: `get_laravel_package_recommendations`

**Python Signature:**
```python
def get_laravel_package_recommendations(use_case: str) -> str
```

**Go Implementation:**
```go
func registerGetLaravelPackageRecommendations(s *server.MCPServer, pc *PackageCatalog) {
    tool := mcp.NewTool("get_laravel_package_recommendations",
        mcp.WithDescription("Intelligently recommends Laravel packages based on described use cases or implementation needs."),
        mcp.WithString("use_case",
            mcp.Required(),
            mcp.Description("Description of what the user wants to implement"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        useCase, err := request.RequireString("use_case")
        if err != nil {
            return mcp.NewToolResultError("use_case is required"), nil
        }
        
        recommendations := pc.RecommendByUseCase(useCase)
        result := formatPackageRecommendations(useCase, recommendations)
        
        return mcp.NewToolResultText(result), nil
    })
}
```

**Response Format:**
```markdown
# Laravel Packages for: implementing payment system

## 1. Laravel Cashier
Laravel Cashier provides an expressive, fluent interface to Stripe's...

**Use Cases:**
- Implementing subscription billing
- Processing one-time payments

**Installation:**
```bash
composer require laravel/cashier
```

**Documentation:** laravel://packages/cashier.md
```

---

### Tool 8: `get_laravel_package_info`

**Python Signature:**
```python
def get_laravel_package_info(package_name: str) -> str
```

**Go Implementation:**
```go
func registerGetLaravelPackageInfo(s *server.MCPServer, pc *PackageCatalog) {
    tool := mcp.NewTool("get_laravel_package_info",
        mcp.WithDescription("Provides comprehensive details about a specific Laravel package including installation and use cases."),
        mcp.WithString("package_name",
            mcp.Required(),
            mcp.Description("The name of the package (e.g., 'laravel/cashier')"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        packageName, err := request.RequireString("package_name")
        if err != nil {
            return mcp.NewToolResultError("package_name is required"), nil
        }
        
        pkg, err := pc.GetPackage(packageName)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Package not found: %s", packageName)), nil
        }
        
        result := formatPackageInfo(pkg)
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 9: `get_laravel_package_categories`

**Python Signature:**
```python
def get_laravel_package_categories(category: str) -> str
```

**Go Implementation:**
```go
func registerGetLaravelPackageCategories(s *server.MCPServer, pc *PackageCatalog) {
    tool := mcp.NewTool("get_laravel_package_categories",
        mcp.WithDescription("Lists all packages within a specific functional category."),
        mcp.WithString("category",
            mcp.Required(),
            mcp.Description("The category to filter by (e.g., 'authentication', 'payment')"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        category, err := request.RequireString("category")
        if err != nil {
            return mcp.NewToolResultError("category is required"), nil
        }
        
        packages := pc.GetByCategory(category)
        if len(packages) == 0 {
            availableCategories := pc.GetAllCategories()
            return mcp.NewToolResultText(fmt.Sprintf(
                "No packages found in category: '%s'. Available categories: %s",
                category,
                strings.Join(availableCategories, ", "),
            )), nil
        }
        
        result := formatCategoryPackages(category, packages)
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 10: `get_features_for_laravel_package`

**Python Signature:**
```python
def get_features_for_laravel_package(package: str) -> str
```

**Go Implementation:**
```go
func registerGetFeaturesForLaravelPackage(s *server.MCPServer, pc *PackageCatalog) {
    tool := mcp.NewTool("get_features_for_laravel_package",
        mcp.WithDescription("Details common implementation features and patterns for a specific package."),
        mcp.WithString("package",
            mcp.Required(),
            mcp.Description("The Laravel package name (e.g., 'laravel/cashier')"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        packageName, err := request.RequireString("package")
        if err != nil {
            return mcp.NewToolResultError("package is required"), nil
        }
        
        pkg, err := pc.GetPackage(packageName)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Package not found: %s", packageName)), nil
        }
        
        features := pc.GetFeatures(packageName)
        result := formatPackageFeatures(pkg, features)
        
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 11: `update_laravel_docs`

**Python Signature:**
```python
def update_laravel_docs(version_param: Optional[str] = None, force: bool = False) -> str
```

**Go Implementation:**
```go
func registerUpdateLaravelDocs(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("update_laravel_docs",
        mcp.WithDescription("Updates documentation from the official Laravel GitHub repository. Ensures access to the latest documentation changes."),
        mcp.WithString("version_param",
            mcp.Description("Laravel version branch (e.g., '12.x')"),
        ),
        mcp.WithBoolean("force",
            mcp.DefaultBool(false),
            mcp.Description("Force update even if already up to date"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        version := mcp.ParseString(request, "version_param", "")
        force := mcp.ParseBoolean(request, "force", false)
        
        if version == "" {
            version = dm.DefaultVersion
        }
        
        result, err := dm.Update(version, force)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Update failed: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

**Response Format:**
```markdown
Documentation updated successfully to 12.x
Commit: abc1234
Date: 2024-01-15
Message: Updated middleware documentation
```

---

### Tool 12: `laravel_docs_info`

**Python Signature:**
```python
def laravel_docs_info(version: Optional[str] = None) -> str
```

**Go Implementation:**
```go
func registerLaravelDocsInfo(s *server.MCPServer, dm *DocManager) {
    tool := mcp.NewTool("laravel_docs_info",
        mcp.WithDescription("Provides metadata about documentation versions, including last update times and commit information."),
        mcp.WithString("version",
            mcp.Description("Specific Laravel version to get info for (e.g., '12.x'). If not provided, shows all versions"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        version := mcp.ParseString(request, "version", "")
        
        info, err := dm.GetInfo(version)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to get info: %v", err)), nil
        }
        
        return mcp.NewToolResultText(info), nil
    })
}
```

---

### Tool 13: `update_external_laravel_docs`

**Python Signature:**
```python
def update_external_laravel_docs(services: Optional[List[str]] = None, force: bool = False) -> str
```

**Go Implementation:**
```go
func registerUpdateExternalLaravelDocs(s *server.MCPServer, em *ExternalManager) {
    tool := mcp.NewTool("update_external_laravel_docs",
        mcp.WithDescription("Updates documentation for external Laravel services like Forge, Vapor, Envoyer, and Nova."),
        mcp.WithArray("services",
            mcp.Description("List of services to update (forge, vapor, envoyer, nova). If None, updates all"),
            mcp.WithStringItems(),
        ),
        mcp.WithBoolean("force",
            mcp.DefaultBool(false),
            mcp.Description("Force update even if cache is valid"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Parse services array
        servicesMap := mcp.ParseStringMap(request, "services", nil)
        var services []string
        if servicesMap != nil {
            // Extract services from map
            for _, v := range servicesMap {
                if str, ok := v.(string); ok {
                    services = append(services, str)
                }
            }
        }
        
        force := mcp.ParseBoolean(request, "force", false)
        
        result, err := em.UpdateServices(services, force)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Update failed: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 14: `list_laravel_services`

**Python Signature:**
```python
def list_laravel_services() -> str
```

**Go Implementation:**
```go
func registerListLaravelServices(s *server.MCPServer, em *ExternalManager) {
    tool := mcp.NewTool("list_laravel_services",
        mcp.WithDescription("Lists all available Laravel services with external documentation support."),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        result := em.ListServices()
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 15: `search_external_laravel_docs`

**Python Signature:**
```python
def search_external_laravel_docs(query: str, services: Optional[List[str]] = None) -> str
```

**Go Implementation:**
```go
func registerSearchExternalLaravelDocs(s *server.MCPServer, em *ExternalManager) {
    tool := mcp.NewTool("search_external_laravel_docs",
        mcp.WithDescription("Searches through external Laravel service documentation."),
        mcp.WithString("query",
            mcp.Required(),
            mcp.Description("Search term to look for"),
        ),
        mcp.WithArray("services",
            mcp.Description("List of services to search. If None, searches all cached services"),
            mcp.WithStringItems(),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        query, err := request.RequireString("query")
        if err != nil {
            return mcp.NewToolResultError("query is required"), nil
        }
        
        // Parse services
        servicesMap := mcp.ParseStringMap(request, "services", nil)
        var services []string
        if servicesMap != nil {
            for _, v := range servicesMap {
                if str, ok := v.(string); ok {
                    services = append(services, str)
                }
            }
        }
        
        result, err := em.Search(query, services)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
        }
        
        return mcp.NewToolResultText(result), nil
    })
}
```

---

### Tool 16: `get_laravel_service_info`

**Python Signature:**
```python
def get_laravel_service_info(service: str) -> str
```

**Go Implementation:**
```go
func registerGetLaravelServiceInfo(s *server.MCPServer, em *ExternalManager) {
    tool := mcp.NewTool("get_laravel_service_info",
        mcp.WithDescription("Provides detailed information about a specific Laravel service."),
        mcp.WithString("service",
            mcp.Required(),
            mcp.Description("Service name (forge, vapor, envoyer, nova)"),
        ),
    )
    
    s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        service, err := request.RequireString("service")
        if err != nil {
            return mcp.NewToolResultError("service is required"), nil
        }
        
        info, err := em.GetServiceInfo(service)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Service not found: %s", service)), nil
        }
        
        return mcp.NewToolResultText(info), nil
    })
}
```

---

## 📝 Complete Tool Registration Function

```go
func registerAllTools(
    s *server.MCPServer,
    dm *DocManager,
    pc *PackageCatalog,
) {
    // Documentation tools (1-6)
    registerListLaravelDocs(s, dm)
    registerReadLaravelDocContent(s, dm)
    registerSearchLaravelDocs(s, dm)
    registerSearchLaravelDocsWithContext(s, dm)
    registerGetDocStructure(s, dm)
    registerBrowseDocsByCategory(s, dm)
    
    // Package tools (7-10)
    registerGetLaravelPackageRecommendations(s, pc)
    registerGetLaravelPackageInfo(s, pc)
    registerGetLaravelPackageCategories(s, pc)
    registerGetFeaturesForLaravelPackage(s, pc)
    
    // Update tools (11-12)
    registerUpdateLaravelDocs(s, dm)
    registerLaravelDocsInfo(s, dm)
    
    // External service tools (13-16)
    em := NewExternalManager("./docs/external")
    registerUpdateExternalLaravelDocs(s, em)
    registerListLaravelServices(s, em)
    registerSearchExternalLaravelDocs(s, em)
    registerGetLaravelServiceInfo(s, em)
    
    log.Printf("Registered all 16 MCP tools successfully")
}
```

---

## ✅ Validation Checklist

Untuk memastikan implementasi Go **persis sama** dengan Python:

### Parameter Schema Validation
- [ ] Semua 16 tools terdaftar dengan nama yang sama
- [ ] Setiap parameter memiliki tipe yang benar (string, boolean, number, array)
- [ ] Required parameters ditandai dengan `mcp.Required()`
- [ ] Optional parameters memiliki default value yang sama
- [ ] Deskripsi parameter identik dengan Python

### Response Format Validation
- [ ] Semua response dalam format Markdown string
- [ ] Struktur Markdown identik (heading levels, bullet points, code blocks)
- [ ] Error messages format sama
- [ ] Success messages format sama

### Behavior Validation
- [ ] Path safety checks untuk file operations
- [ ] Cache behavior sama (TTL, eviction policy)
- [ ] Search algorithm sama (case-insensitive, regex)
- [ ] Scoring algorithm sama untuk package recommendations

---

## 🚀 Build & Test

```bash
# Build
go build -o laravel-mcp-companion ./cmd/server

# Test dengan Claude Desktop
# Edit config: ~/Library/Application Support/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "laravel-companion": {
      "command": "/path/to/laravel-mcp-companion",
      "args": ["--docs-path", "./docs", "--version", "12.x"]
    }
  }
}

# Restart Claude Desktop dan test setiap tool
```

---

## 📚 Complete Tool Descriptions (Copy dari Python)

Setiap tool harus memiliki description yang **persis sama** dengan Python `TOOL_DESCRIPTIONS`:

```go
var ToolDescriptions = map[string]string{
    "list_laravel_docs": "Lists all available Laravel documentation files across versions. Essential for discovering what documentation exists before diving into specific topics.\n\nWhen to use:\n- Initial exploration of Laravel documentation\n- Finding available documentation files\n- Checking which versions have specific documentation\n- Getting an overview of documentation coverage",
    
    "read_laravel_doc_content": "Reads the complete content of a specific Laravel documentation file. This is the primary tool for accessing actual documentation content.\n\nWhen to use:\n- Reading full documentation for a feature\n- Getting complete implementation details\n- Accessing code examples from docs\n- Understanding concepts in depth",
    
    // ... (all 16 descriptions)
}
```

---

Dokumen ini adalah **complete reference** untuk memastikan implementasi Go 100% kompatibel dengan Python MCP server.
