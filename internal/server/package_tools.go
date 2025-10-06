package server

import (
	"context"
	"fmt"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/packages"
	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterPackageTools registers all package-related MCP tools
func (s *Server) RegisterPackageTools(catalog *packages.Catalog) {
	// Tool 1: List package categories
	listCatTool := mcp.NewTool("list_package_categories",
		mcp.WithDescription("List all available Laravel package categories"),
	)

	s.mcp.AddTool(listCatTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categories := catalog.ListCategories()

		formatted, err := packages.FormatCategoriesList(categories)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to format categories: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})

	// Tool 2: Recommend Laravel packages
	recommendTool := mcp.NewTool("recommend_laravel_packages",
		mcp.WithDescription("Get package recommendations based on use case"),
		mcp.WithString("use_case",
			mcp.Required(),
			mcp.Description("What you want to accomplish (e.g. 'authentication', 'file uploads', 'api')"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of recommendations (default: 5)"),
		),
	)

	s.mcp.AddTool(recommendTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		useCase, err := request.RequireString("use_case")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		limit := request.GetInt("limit", 5)

		recommendations := catalog.Recommend(useCase, limit)

		formatted, err := packages.FormatPackageList(
			recommendations,
			fmt.Sprintf("Recommendations for '%s'", useCase),
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to format recommendations: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})

	// Tool 3: Search packages
	searchTool := mcp.NewTool("search_packages",
		mcp.WithDescription("Search for Laravel packages with filters"),
		mcp.WithString("query",
			mcp.Description("Search query (package name, description, or tags)"),
		),
		mcp.WithString("category",
			mcp.Description("Filter by category name"),
		),
		mcp.WithBoolean("maintained_only",
			mcp.Description("Show only actively maintained packages"),
		),
		mcp.WithNumber("min_popularity",
			mcp.Description("Minimum popularity score (0-100)"),
		),
		mcp.WithArray("tags",
			mcp.Description("Filter by tags (array of tag strings)"),
			mcp.WithStringItems(),
		),
	)

	s.mcp.AddTool(searchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		category := request.GetString("category", "")

		// Handle category filter
		if category != "" {
			cat, err := catalog.GetCategory(category)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("category not found: %s", category)), nil
			}

			formatted, err := packages.FormatCategoryPackages(cat)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to format category: %v", err)), nil
			}

			return mcp.NewToolResultText(formatted), nil
		}

		// Build filters
		filters := make(map[string]interface{})
		if maintainedOnly := request.GetBool("maintained_only", false); maintainedOnly {
			filters["maintained"] = maintainedOnly
		}
		if minPop := request.GetInt("min_popularity", 0); minPop > 0 {
			filters["min_popularity"] = float64(minPop)
		}
		if tags := request.GetStringSlice("tags", []string{}); len(tags) > 0 {
			filters["tags"] = tags
		}

		results := catalog.Search(query, filters)

		title := "Search Results"
		if query != "" {
			title = fmt.Sprintf("Search Results for '%s'", query)
		}

		formatted, err := packages.FormatPackageList(results, title)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to format search results: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})

	// Tool 4: Get package details
	detailsTool := mcp.NewTool("get_package_details",
		mcp.WithDescription("Get detailed information about a specific package"),
		mcp.WithString("composer_name",
			mcp.Required(),
			mcp.Description("Composer package name (e.g. 'spatie/laravel-permission')"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: 'markdown' or 'json' (default: markdown)"),
		),
	)

	s.mcp.AddTool(detailsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		composerName, err := request.RequireString("composer_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		pkg, err := catalog.GetPackage(composerName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		format := request.GetString("format", "markdown")

		var formatted string
		if format == "json" {
			formatted, err = packages.FormatPackageJSON(pkg)
		} else {
			formatted, err = packages.FormatPackageDetails(pkg)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to format package details: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})
}
