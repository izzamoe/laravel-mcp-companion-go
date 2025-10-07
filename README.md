# Laravel MCP Companion (Go Implementation)

A Model Context Protocol (MCP) server providing comprehensive Laravel documentation and package recommendations for AI assistants. This is a complete Go implementation using the official `github.com/mark3labs/mcp-go` library.

## 🚀 Features

### 16 MCP Tools Available

#### Documentation Tools (6 tools)
1. **`list_laravel_docs`** - List available documentation files
2. **`read_laravel_doc_content`** - Read complete documentation content
3. **`search_laravel_docs`** - Search across documentation with match counts
4. **`search_laravel_docs_with_context`** - Advanced search with surrounding context
5. **`get_doc_structure`** - Extract table of contents from documentation
6. **`browse_docs_by_category`** - Discover docs by category (frontend, database, etc.)

#### Package Tools (4 tools)
7. **`get_laravel_package_recommendations`** - Get package recommendations by use case
8. **`get_laravel_package_info`** - Detailed information about specific packages
9. **`get_laravel_package_categories`** - List packages within a category
10. **`get_features_for_laravel_package`** - Get common features and patterns

#### Update & Info Tools (2 tools)
11. **`update_laravel_docs`** - Update documentation from GitHub
12. **`laravel_docs_info`** - Get metadata about documentation versions

#### External Service Tools (4 tools)
13. **`update_external_laravel_docs`** - Update external service documentation
14. **`list_laravel_services`** - List available Laravel services (Forge, Vapor, etc.)
15. **`search_external_laravel_docs`** - Search external service documentation
16. **`get_laravel_service_info`** - Get detailed service information

## 📦 Installation

### Prerequisites
- Go 1.24 or later
- Git (for documentation updates)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/izzamoe/laravel-mcp-companion-go.git
cd laravel-mcp-companion-go

# Install dependencies
go mod download

# Build the server
go build -o bin/server ./cmd/server

# Run the server
./bin/server --docs-path ./docs --version 12.x
```

## 🔧 Configuration

### Command Line Flags

- `--docs-path` - Path to documentation directory (default: `./docs`)
- `--packages-path` - Path to packages catalog (default: `./configs/packages.json`)
- `--version` - Default Laravel version (default: `12.x`)
- `--log-level` - Logging level: debug, info, warn, error (default: `info`)

### Claude Desktop Configuration

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "laravel-companion": {
      "command": "/path/to/laravel-mcp-companion-go/bin/server",
      "args": [
        "--docs-path", "/path/to/laravel-mcp-companion-go/docs",
        "--version", "12.x",
        "--log-level", "info"
      ]
    }
  }
}
```

## 📚 Usage Examples

### 1. List Available Documentation

```
Tool: list_laravel_docs
Parameters:
  version: "12.x"
```

### 2. Read Documentation Content

```
Tool: read_laravel_doc_content
Parameters:
  filename: "routing.md"
  version: "12.x"
```

### 3. Search Documentation

```
Tool: search_laravel_docs
Parameters:
  query: "middleware"
  version: "12.x"
```

### 4. Search with Context

```
Tool: search_laravel_docs_with_context
Parameters:
  query: "middleware"
  context_length: 200
```

### 5. Get Package Recommendations

```
Tool: get_laravel_package_recommendations
Parameters:
  use_case: "implementing payment system"
```

### 6. Get Package Information

```
Tool: get_laravel_package_info
Parameters:
  package_name: "laravel/cashier"
```

### 7. Browse by Category

```
Tool: browse_docs_by_category
Parameters:
  category: "frontend"
  version: "12.x"
```

## 🏗️ Architecture

```
laravel-mcp-companion-go/
├── cmd/server/           # Main entry point
│   └── main.go
├── internal/
│   ├── docs/            # Documentation management
│   │   ├── manager.go   # Core doc operations
│   │   └── cache.go     # Documentation caching
│   ├── packages/        # Package catalog
│   │   ├── catalog.go   # Package search & recommendations
│   │   └── format.go    # Output formatting
│   ├── server/          # MCP server & tools
│   │   ├── server.go    # Server initialization
│   │   ├── doc_tools.go # Documentation tools (6)
│   │   ├── package_tools.go # Package tools (4)
│   │   └── external_tools.go # External tools (6)
│   ├── updater/         # GitHub documentation updater
│   ├── external/        # External resource handling
│   ├── logging/         # Logging utilities
│   └── models/          # Data models
├── docs/                # Documentation files
│   └── 12.x/           # Laravel 12.x docs
├── configs/             # Configuration files
│   └── packages.json   # Package catalog
└── bin/                # Built binaries
    └── server
```

## 🔍 Implementation Details

### MCP Library

This implementation uses the official MCP Go library:
- **Library:** `github.com/mark3labs/mcp-go` v0.41.1
- **Transport:** stdio (Standard Input/Output)
- **Protocol Version:** 2024-11-05

### Tool Registration

All tools are registered with:
- ✅ Proper parameter schemas (required/optional)
- ✅ Type validation (string, number, boolean, array)
- ✅ Default values for optional parameters
- ✅ Comprehensive descriptions with "When to use" guidance
- ✅ Error handling and response formatting

### Documentation Manager

Features:
- File-based documentation storage
- In-memory caching with TTL
- Path safety validation (prevents directory traversal)
- Version management (12.x, 11.x, etc.)
- Search with context extraction
- Structure parsing (TOC generation)
- Category-based browsing

### Package Catalog

Features:
- JSON-based package index
- Use case matching with relevance scoring
- Category organization
- Popularity scoring
- Maintenance status tracking
- Alternative package suggestions

## 🧪 Testing

### Manual Testing

```bash
# Build the server
go build -o bin/server ./cmd/server

# Run test script
./test_tools.sh

# Or manually test with Claude Desktop
```

### Verify Tool Count

The server should log on startup:
```
Registered documentation tools (6 tools)
Registered package tools (4 tools)
Registered update and info tools (2 tools)
Registered external service tools (4 tools)
Server ready with 16 total tools, starting event loop...
```

## 📝 Tool Specifications

See [TOOLS_VERIFICATION.md](TOOLS_VERIFICATION.md) for complete tool specifications and compliance checklist.

See [MCP_GO_IMPLEMENTATION_COMPLETE.md](MCP_GO_IMPLEMENTATION_COMPLETE.md) for detailed implementation guide.

## 🚀 Development

### Building

```bash
go build -o bin/server ./cmd/server
```

### Running

```bash
./bin/server --docs-path ./docs --version 12.x --log-level debug
```

### Logging

Logs are output to stderr (stdout is reserved for MCP protocol):
- **Debug:** Detailed operation logs
- **Info:** General operation logs (default)
- **Warn:** Warning messages
- **Error:** Error messages

## 🤝 Contributing

Contributions are welcome! Please ensure:
- All 16 tools remain functional
- Parameter schemas match the specification
- Response formats are Markdown-compatible
- Tests pass
- Code follows Go conventions

## 📄 License

MIT License

## 🙏 Acknowledgments

- Laravel Framework team for excellent documentation
- Anthropic for the MCP protocol specification
- mark3labs for the mcp-go library

## 📞 Support

For issues and questions:
- Create an issue on GitHub
- Check the documentation in `docs/`
- Review the implementation guide in `MCP_GO_IMPLEMENTATION_COMPLETE.md`

---

**Status:** ✅ All 16 tools implemented and tested  
**Version:** 1.0.0  
**Go Version:** 1.24.0  
**Binary Size:** 9.7MB

> 🚀 MCP Server for Laravel documentation and package recommendations

[![Build Status](https://github.com/izzamoe/laravel-mcp-companion-go/workflows/Test/badge.svg)](https://github.com/izzamoe/laravel-mcp-companion-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/izzamoe/laravel-mcp-companion-go)](https://goreportcard.com/report/github.com/izzamoe/laravel-mcp-companion-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ Features

- 📚 **Complete Laravel Documentation** - Versions 6.x through 12.x
- 🔍 **Smart Package Recommendations** - AI-powered suggestions based on your use case
- 🌐 **External Services** - Forge, Vapor, Nova, Envoyer documentation
- ⚡ **High Performance** - Built with Go for speed and efficiency
- 💾 **Intelligent Caching** - Fast responses with automatic cache management
- 🎯 **16 MCP Tools** - Comprehensive API for Laravel development
- 🔗 **Resource URIs** - Direct access via `laravel://` and `laravel-external://`

## 📋 Prerequisites

- Go 1.24 or higher
- Claude Desktop app or VSCode with MCP support
- macOS, Linux, or Windows

## 🚀 Installation

### Download Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/izzamoe/laravel-mcp-companion-go/releases) for your platform:

- **Linux (AMD64/ARM64)**: `laravel-mcp-companion-go-linux-{arch}`
- **macOS (AMD64/ARM64)**: `laravel-mcp-companion-go-darwin-{arch}`
- **Windows (AMD64)**: `laravel-mcp-companion-go-windows-amd64.exe`

Make the binary executable and move to your preferred location:

```bash
# For Linux/macOS
chmod +x laravel-mcp-companion-go-*
sudo mv laravel-mcp-companion-go-* /usr/local/bin/laravel-mcp-companion-go

# For Windows, just move the .exe file to your desired location
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/izzamoe/laravel-mcp-companion-go.git
cd laravel-mcp-companion-go

# Build
go build -o laravel-mcp-companion-go cmd/server/main.go

# Move to your preferred location
mv laravel-mcp-companion-go /usr/local/bin/
```

## ⚙️ Configuration

### Claude Desktop Setup

Add to your Claude Desktop config:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "laravel-companion": {
      "command": "/usr/local/bin/laravel-mcp-companion-go",
      "args": []
    }
  }
}
```

### VSCode Setup

1. **Download or build the binary** (see Installation section above).

2. **Create MCP configuration**:
   - Create a `.vscode` folder in your workspace root (if it doesn't exist)
   - Create a file named `mcp.json` inside `.vscode/`
   - Add the following content to `mcp.json`:

   ```json
   {
     "servers": {
       "laravel_docs": {
         "type": "stdio",
         "command": "/absolute/path/to/laravel-mcp-companion-go",
         "args": []
       }
     },
     "inputs": []
   }
   ```

   **Note**: Replace `/absolute/path/to/laravel-mcp-companion-go` with the actual absolute path to your binary (e.g., `/usr/local/bin/laravel-mcp-companion-go`).

3. **Restart VSCode** or reload the window to apply the MCP configuration.

4. **Verify setup**: The MCP server should now be available in VSCode's MCP-enabled extensions.

## 📖 Usage

### In Claude Desktop

Simply ask Claude about Laravel! Examples:

- **Documentation**: "Show me Laravel 11.x routing documentation"
- **Packages**: "Recommend packages for implementing payment processing"
- **Features**: "What features does Laravel Sanctum provide?"
- **Search**: "Search for middleware in Laravel 11.x"

## 🏗️ Development

### Run Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Run with Debug Logging

```bash
LOG_LEVEL=debug ./laravel-mcp-companion-go
```

## 📝 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Laravel Framework](https://laravel.com/) - The PHP framework
- [MCP Protocol](https://modelcontextprotocol.io/) - Model Context Protocol
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Go MCP library

---

**Made with ❤️ for the Laravel community**
