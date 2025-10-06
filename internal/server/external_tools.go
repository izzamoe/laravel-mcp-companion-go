package server

import (
	"context"
	"fmt"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/external"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/updater"
	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterExternalTools registers update and external resource tools
func (s *Server) RegisterExternalTools(upd *updater.GitHubUpdater, scraper *external.WebScraper) {
	// Tool 1: Update Laravel documentation
	updateDocsTool := mcp.NewTool("update_laravel_docs",
		mcp.WithDescription("Fetch and update Laravel documentation from GitHub for a specific version"),
		mcp.WithString("version",
			mcp.Required(),
			mcp.Description("Laravel version to update (e.g. '12.x', '11.x')"),
		),
	)

	s.mcp.AddTool(updateDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version, err := request.RequireString("version")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := upd.UpdateDocs(version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("update failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Tool 2: Get external resource
	getResourceTool := mcp.NewTool("get_external_resource",
		mcp.WithDescription("Fetch content from an external URL (for documentation, package info, etc.)"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL of the external resource to fetch (HTTP/HTTPS only)"),
		),
	)

	s.mcp.AddTool(getResourceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := request.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		content, err := scraper.FetchResource(url)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to fetch resource: %v", err)), nil
		}

		formatted := scraper.FormatResource(url, content)

		return mcp.NewToolResultText(formatted), nil
	})
}
