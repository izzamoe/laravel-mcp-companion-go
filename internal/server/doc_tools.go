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
		mcp.WithDescription("List all available Laravel documentation files across versions. Essential for discovering what documentation exists before diving into specific topics.\n\nWhen to use:\n- Initial exploration of Laravel documentation\n- Finding available documentation files\n- Checking which versions have specific documentation\n- Getting an overview of documentation coverage"),
		mcp.WithString("version",
			mcp.Description("Specific Laravel version to list (e.g., '12.x'). If not provided, lists all versions"),
		),
	)

	s.mcp.AddTool(listDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := mcp.ParseString(request, "version", "")

		files, err := s.docManager.ListDocs(version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list docs: %v", err)), nil
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
		mcp.WithDescription("Reads the complete content of a specific Laravel documentation file. This is the primary tool for accessing actual documentation content.\n\nWhen to use:\n- Reading full documentation for a feature\n- Getting complete implementation details\n- Accessing code examples from docs\n- Understanding concepts in depth"),
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Description("Name of the file (e.g., 'mix.md', 'vite.md')"),
		),
		mcp.WithString("version",
			mcp.Description("Laravel version (e.g., '12.x'). Defaults to latest"),
		),
	)

	s.mcp.AddTool(readDocTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError("filename is required"), nil
		}

		version := mcp.ParseString(request, "version", "")

		content, err := s.docManager.ReadDoc(version, filename)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read doc: %v", err)), nil
		}

		return mcp.NewToolResultText(content), nil
	})

	// Tool 3: search_laravel_docs
	searchDocsTool := mcp.NewTool("search_laravel_docs",
		mcp.WithDescription("Searches for specific terms across all Laravel documentation files. Returns file names and match counts.\n\nWhen to use:\n- Finding which files contain specific topics\n- Getting quick overview of where a concept is mentioned\n- Discovering related documentation\n- Checking documentation coverage for a feature"),
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

	s.mcp.AddTool(searchDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("query is required"), nil
		}

		version := mcp.ParseString(request, "version", "")
		// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search

		results, err := s.docManager.SearchDocs(query, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		return mcp.NewToolResultText(results), nil
	})

	// Tool 4: search_laravel_docs_with_context
	searchWithContextTool := mcp.NewTool("search_laravel_docs_with_context",
		mcp.WithDescription("Advanced search that returns matching text with surrounding context. Shows exactly how terms are used in documentation.\n\nWhen to use:\n- Understanding how a term is used in context\n- Getting code examples that use specific features\n- Finding usage patterns\n- Quick answers without reading full docs"),
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

	s.mcp.AddTool(searchWithContextTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("query is required"), nil
		}

		version := mcp.ParseString(request, "version", "")
		contextLength := int(mcp.ParseFloat64(request, "context_length", 200))
		// includeExternal := mcp.ParseBoolean(request, "include_external", true) // TODO: implement external search

		result, err := s.docManager.SearchWithContext(query, version, contextLength)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Tool 5: get_doc_structure
	getStructureTool := mcp.NewTool("get_doc_structure",
		mcp.WithDescription("Extracts the table of contents and structure from a documentation file. Shows headers and brief content previews.\n\nWhen to use:\n- Getting an overview of documentation organization\n- Finding specific sections quickly\n- Understanding document layout\n- Navigation planning"),
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Description("Documentation file name"),
		),
		mcp.WithString("version",
			mcp.Description("Laravel version"),
		),
	)

	s.mcp.AddTool(getStructureTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError("filename is required"), nil
		}

		version := mcp.ParseString(request, "version", "")

		structure, err := s.docManager.GetStructure(filename, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get structure: %v", err)), nil
		}

		return mcp.NewToolResultText(structure), nil
	})

	// Tool 6: browse_docs_by_category
	browseCategoryTool := mcp.NewTool("browse_docs_by_category",
		mcp.WithDescription("Discovers documentation files related to specific categories like 'frontend', 'database', or 'authentication'.\n\nWhen to use:\n- Exploring documentation by topic area\n- Finding all related documentation for a feature\n- Learning about a domain (frontend, database, etc.)\n- Discovery of related concepts"),
		mcp.WithString("category",
			mcp.Required(),
			mcp.Description("Category like 'frontend', 'database', 'authentication', etc."),
		),
		mcp.WithString("version",
			mcp.Description("Laravel version"),
		),
	)

	s.mcp.AddTool(browseCategoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category, err := request.RequireString("category")
		if err != nil {
			return mcp.NewToolResultError("category is required"), nil
		}

		version := mcp.ParseString(request, "version", "")

		result, err := s.docManager.BrowseByCategory(category, version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to browse: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	return nil
}
