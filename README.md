# Laravel MCP Companion (Go Implementation)

A high-performance MCP server providing comprehensive Laravel documentation and intelligent package recommendations for AI assistants. Built with Go using the official MCP SDK.

## ✨ Features

- 📚 **16 MCP Tools** - Complete Laravel development toolkit
- 🔍 **Smart Documentation** - Search across Laravel 6.x-12.x docs
- 📦 **Package Intelligence** - AI-powered recommendations by use case
- 🌐 **External Services** - Forge, Vapor, Nova, Envoyer integration
- ⚡ **Go Performance** - Fast, efficient, and lightweight
- 💾 **Intelligent Caching** - Optimized response times

### MCP Tools Overview
- **Documentation** (6): Browse, search, and extract Laravel docs
- **Packages** (4): Recommendations, info, and category browsing
- **Updates** (2): Documentation and metadata management
- **External** (4): Laravel ecosystem service documentation

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Claude Desktop or VSCode with MCP support

### Install & Run

```bash
# Clone and build
git clone https://github.com/izzamoe/laravel-mcp-companion-go.git
cd laravel-mcp-companion-go
go build -o bin/server ./cmd/server

# Run server
./bin/server --docs-path ./docs --version 12.x
```

### Claude Desktop Config

Add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "laravel-companion": {
      "command": "/path/to/bin/server",
      "args": ["--docs-path", "/path/to/docs", "--version", "12.x"]
    }
  }
}
```

### VSCode Setup

1. Install MCP extension
2. Create `.vscode/mcp.json`:
```json
{
  "servers": {
    "laravel-companion": {
      "type": "stdio",
      "command": "/path/to/bin/server",
      "args": ["--docs-path", "/path/to/docs"]
    }
  }
}
```

## 📖 Usage Examples

### In Claude Desktop
Simply ask Claude about Laravel development:

- *"Show me Laravel 12.x routing documentation"*
- *"Recommend packages for payment processing"*
- *"Search for middleware examples in Laravel"*
- *"What are Laravel Sanctum features?"*
- *"How to implement Laravel Cashier?"*

### In VSCode
Once configured, ask your AI assistant:

- *"Browse Laravel documentation categories"*
- *"Find authentication packages"*
- *"Show Laravel Forge deployment features"*
- *"Search for database migration docs"*
- *"Get package info for laravel/socialite"*

### Tool Examples

**Documentation Search:**
```
Search: "middleware authentication"
Result: Links to relevant docs with context
```

**Package Recommendations:**
```
Use case: "implementing user notifications"
Result: Recommended packages with descriptions
```

**External Services:**
```
Service: "Laravel Forge"
Result: Features, pricing, and documentation
```

## 🏗️ Architecture

```
├── cmd/server/          # Main entry point
├── internal/
│   ├── docs/           # Documentation management
│   ├── packages/       # Package catalog
│   ├── server/         # MCP tools (16 total)
│   ├── external/       # Laravel ecosystem services
│   └── models/         # Data structures
├── docs/               # Laravel documentation
└── configs/            # Package catalog
```

## 🧪 Development

```bash
# Run tests
go test ./...

# Build
go build -o bin/server ./cmd/server

# Debug mode
./bin/server --log-level debug
```

## 📋 Command Line Options

- `--docs-path` - Documentation directory (default: `./docs`)
- `--packages-path` - Package catalog (default: `./configs/packages.json`)
- `--version` - Default Laravel version (default: `12.x`)
- `--log-level` - Logging: debug, info, warn, error (default: `info`)

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Laravel Framework team
- Anthropic MCP Protocol
- Model Context Protocol Go SDK

---

**Status:** ✅ 16/16 tools implemented | **Version:** 1.0.0 | **Go:** 1.24+

*Made with ❤️ for the Laravel community*
