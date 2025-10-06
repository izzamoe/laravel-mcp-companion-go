# Laravel MCP Companion (Go)

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
