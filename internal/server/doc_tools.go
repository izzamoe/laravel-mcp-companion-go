package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterDocTools registers all documentation-related MCP tools
func (s *Server) RegisterDocTools() error {
	// Tool 1: list_laravel_docs
	listDocsTool := mcp.NewTool("list_laravel_docs",
		mcp.WithDescription("List available Laravel documentation files for a specific version"),
		mcp.WithString("version",
			mcp.Description("Specific Laravel version (e.g. '12.x', '11.x'). Leave empty for default 12.x"),
		),
	)

	s.mcp.AddTool(listDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := request.GetString("version", "")

		files, err := s.docManager.ListDocs(version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		var result strings.Builder
		if version == "" {
			version = "12.x" // default
		}
		result.WriteString(fmt.Sprintf("# Laravel %s Documentation\n\n", version))
		result.WriteString(fmt.Sprintf("Found %d documentation files:\n\n", len(files)))

		for _, file := range files {
			result.WriteString(fmt.Sprintf("- %s\n", file))
		}

		return mcp.NewToolResultText(result.String()), nil
	})

	// Tool 2: read_laravel_doc_content
	readDocTool := mcp.NewTool("read_laravel_doc_content",
		mcp.WithDescription("Read the content of a specific Laravel documentation file"),
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Description("Documentation file name (e.g. 'routing.md')"),
		),
		mcp.WithString("version",
			mcp.Description("Laravel version (e.g. '12.x'). Leave empty for default"),
		),
	)

	s.mcp.AddTool(readDocTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		version := request.GetString("version", "")

		content, err := s.docManager.ReadDoc(version, filename)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(content), nil
	})

	// Tool 3: search_laravel_docs
	searchDocsTool := mcp.NewTool("search_laravel_docs",
		mcp.WithDescription("Search for a term across Laravel documentation files"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search term or keyword"),
		),
		mcp.WithString("version",
			mcp.Description("Laravel version to search in (e.g. '12.x'). Leave empty for default"),
		),
	)

	s.mcp.AddTool(searchDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		version := request.GetString("version", "")

		results, err := s.docManager.SearchDocs(query, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(results), nil
	})

	return nil
}
