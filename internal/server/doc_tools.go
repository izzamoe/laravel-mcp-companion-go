package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input/output types for all 6 documentation tools
type ListDocsInput struct {
	Version string `json:"version,omitempty" jsonschema:"Specific Laravel version to list (e.g. '12.x'). If not provided lists all versions"`
}

type ReadDocInput struct {
	Filename string `json:"filename" jsonschema:"required,Name of the file (e.g. 'mix.md' 'vite.md')"`
	Version  string `json:"version,omitempty" jsonschema:"Laravel version (e.g. '12.x'). Defaults to latest"`
}

type SearchDocsInput struct {
	Query           string `json:"query" jsonschema:"required,Search term to look for"`
	Version         string `json:"version,omitempty" jsonschema:"Specific Laravel version to search (e.g. '12.x'). If not provided searches all versions"`
	IncludeExternal *bool  `json:"include_external,omitempty" jsonschema:"Whether to include external Laravel services documentation in search"`
}

type SearchWithContextInput struct {
	Query           string `json:"query" jsonschema:"required,Search term"`
	Version         string `json:"version,omitempty" jsonschema:"Specific version or None for all"`
	ContextLength   *int   `json:"context_length,omitempty" jsonschema:"Characters of context to show (default: 200)"`
	IncludeExternal *bool  `json:"include_external,omitempty" jsonschema:"Include external services"`
}

type GetStructureInput struct {
	Filename string `json:"filename" jsonschema:"required,Documentation file name"`
	Version  string `json:"version,omitempty" jsonschema:"Laravel version"`
}

type BrowseCategoryInput struct {
	Category string `json:"category" jsonschema:"required,Category like 'frontend' 'database' 'authentication' etc."`
	Version  string `json:"version,omitempty" jsonschema:"Laravel version"`
}

// Empty output types - we return text content
type EmptyOutput struct{}

// RegisterDocTools registers all 6 documentation-related MCP tools
func (s *Server) RegisterDocTools() error {
	// Tool 1: list_laravel_docs
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_laravel_docs",
		Description: "List all available Laravel documentation files across versions. Essential for discovering what documentation exists before diving into specific topics.\n\nWhen to use:\n- Initial exploration of Laravel documentation\n- Finding available documentation files\n- Checking which versions have specific documentation\n- Getting an overview of documentation coverage",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input ListDocsInput) (*mcp.CallToolResult, EmptyOutput, error) {
		version := input.Version
		if version == "" {
			version = "12.x"
		}

		files, err := s.docManager.ListDocs(version)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list docs: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("# Laravel %s Documentation\n\n", version))
		result.WriteString(fmt.Sprintf("Found %d documentation files:\n\n", len(files)))
		for _, file := range files {
			result.WriteString(fmt.Sprintf("- %s\n", file))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
		}, EmptyOutput{}, nil
	})

	// Tool 2: read_laravel_doc_content
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "read_laravel_doc_content",
		Description: "Reads the complete content of a specific Laravel documentation file. This is the primary tool for accessing actual documentation content.\n\nWhen to use:\n- Reading full documentation for a feature\n- Getting complete implementation details\n- Accessing code examples from docs\n- Understanding concepts in depth",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input ReadDocInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Filename == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "filename is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		content, err := s.docManager.ReadDoc(input.Version, input.Filename)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to read doc: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: content}},
		}, EmptyOutput{}, nil
	})

	// Tool 3: search_laravel_docs
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search_laravel_docs",
		Description: "Searches for specific terms across all Laravel documentation files. Returns file names and match counts.\n\nWhen to use:\n- Finding which files contain specific topics\n- Getting quick overview of where a concept is mentioned\n- Discovering related documentation\n- Checking documentation coverage for a feature",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input SearchDocsInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		version := input.Version
		includeExternal := true
		if input.IncludeExternal != nil {
			includeExternal = *input.IncludeExternal
		}

		results, err := s.docManager.SearchDocs(input.Query, version)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Search failed: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		if includeExternal && s.externalManager != nil {
			externalResults, err := s.externalManager.SearchServices(input.Query, nil)
			if err == nil && externalResults != "" {
				results += "\n\n---\n\n" + externalResults
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: results}},
		}, EmptyOutput{}, nil
	})

	// Tool 4: search_laravel_docs_with_context
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search_laravel_docs_with_context",
		Description: "Advanced search that returns matching text with surrounding context. Shows exactly how terms are used in documentation.\n\nWhen to use:\n- Understanding how a term is used in context\n- Getting code examples that use specific features\n- Finding usage patterns\n- Quick answers without reading full docs",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input SearchWithContextInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		contextLength := 200
		if input.ContextLength != nil {
			contextLength = *input.ContextLength
		}

		includeExternal := true
		if input.IncludeExternal != nil {
			includeExternal = *input.IncludeExternal
		}

		result, err := s.docManager.SearchWithContext(input.Query, input.Version, contextLength)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Search failed: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		if includeExternal && s.externalManager != nil {
			externalResults, err := s.externalManager.SearchServicesWithContext(input.Query, nil, contextLength)
			if err == nil && externalResults != "" {
				result += "\n\n---\n\n" + externalResults
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	// Tool 5: get_doc_structure
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_doc_structure",
		Description: "Extracts the table of contents and structure from a documentation file. Shows headers and brief content previews.\n\nWhen to use:\n- Getting an overview of documentation organization\n- Finding specific sections quickly\n- Understanding document layout\n- Navigation planning",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input GetStructureInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Filename == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "filename is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		structure, err := s.docManager.GetStructure(input.Filename, input.Version)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get structure: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: structure}},
		}, EmptyOutput{}, nil
	})

	// Tool 6: browse_docs_by_category
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "browse_docs_by_category",
		Description: "Discovers documentation files related to specific categories like 'frontend', 'database', or 'authentication'.\n\nWhen to use:\n- Exploring documentation by topic area\n- Finding all related documentation for a feature\n- Learning about a domain (frontend, database, etc.)\n- Discovery of related concepts",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input BrowseCategoryInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Category == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "category is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		result, err := s.docManager.BrowseByCategory(input.Category, input.Version)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to browse: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	return nil
}
