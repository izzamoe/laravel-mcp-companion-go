package main

import (
	"context"
	"flag"
	"os"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/docs"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/external"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/logging"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/packages"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/server"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/updater"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Parse command line flags
	docsPath := flag.String("docs-path", "./docs", "Path to documentation directory")
	packagesPath := flag.String("packages-path", "./configs/packages.json", "Path to packages catalog")
	defaultVersion := flag.String("version", "12.x", "Default Laravel version")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Configure logging
	switch *logLevel {
	case "debug":
		logging.SetLevel(logging.LevelDebug)
	case "warn":
		logging.SetLevel(logging.LevelWarn)
	case "error":
		logging.SetLevel(logging.LevelError)
	default:
		logging.SetLevel(logging.LevelInfo)
	}

	logging.Info("Starting Laravel MCP Companion...")

	// Initialize documentation manager
	docManager := docs.NewManager(*docsPath, *defaultVersion)
	logging.Info("Initialized documentation manager (path: %s, default: %s)", *docsPath, *defaultVersion)

	// Initialize package catalog
	catalog, err := packages.NewCatalog(*packagesPath)
	if err != nil {
		logging.Error("Failed to initialize package catalog: %v", err)
		os.Exit(1)
	}
	logging.Info("Initialized package catalog (path: %s)", *packagesPath)

	// Initialize updater and scraper
	upd := updater.NewGitHubUpdater(*docsPath)
	scraper := external.NewWebScraper()

	// Initialize external manager with cache path
	externalCachePath := "./cache/external"
	externalManager := external.NewExternalManager(externalCachePath)
	logging.Info("Initialized updater, web scraper, and external manager")

	// Create server
	srv := server.NewServer(docManager)
	srv.SetExternalManager(externalManager)
	srv.SetUpdater(upd)
	logging.Info("Created MCP server")

	// Register documentation tools
	if err := srv.RegisterDocTools(); err != nil {
		logging.Error("Failed to register doc tools: %v", err)
		os.Exit(1)
	}
	logging.Info("Registered documentation tools (6 tools)")

	// Register package tools
	srv.RegisterPackageTools(catalog)
	logging.Info("Registered package tools (4 tools)")

	// Register external tools (update & info)
	srv.RegisterExternalTools(upd, scraper)
	logging.Info("Registered update and info tools (2 tools)")

	// Register external service tools with the external manager
	srv.RegisterExternalServiceTools(externalManager)
	logging.Info("Registered external service tools (4 tools)")

	// Start the server (blocking call)
	logging.Info("Server ready with 16 total tools, starting event loop...")

	// Start the MCP server over stdio using new SDK
	if err := srv.GetMCPServer().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logging.Error("Server error: %v", err)
		os.Exit(1)
	}
}
