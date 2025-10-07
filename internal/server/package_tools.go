package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/packages"
	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterPackageTools registers all package-related MCP tools
func (s *Server) RegisterPackageTools(catalog *packages.Catalog) {
	// Tool 7: get_laravel_package_recommendations
	recommendTool := mcp.NewTool("get_laravel_package_recommendations",
		mcp.WithDescription("Intelligently recommends Laravel packages based on described use cases or implementation needs.\n\nWhen to use:\n- Starting a new feature and need package suggestions\n- Looking for solutions to specific problems\n- Comparing available options for a use case\n- Discovering well-maintained packages"),
		mcp.WithString("use_case",
			mcp.Required(),
			mcp.Description("Description of what the user wants to implement"),
		),
	)

	s.mcp.AddTool(recommendTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		useCase, err := request.RequireString("use_case")
		if err != nil {
			return mcp.NewToolResultError("use_case is required"), nil
		}

		recommendations := catalog.Recommend(useCase, 5)

		var output strings.Builder
		output.WriteString(fmt.Sprintf("# Laravel Packages for: %s\n\n", useCase))

		if len(recommendations) == 0 {
			output.WriteString("No packages found matching your use case. Try different keywords or browse categories.\n")
		} else {
			for i, pkg := range recommendations {
				output.WriteString(fmt.Sprintf("## %d. %s\n", i+1, pkg.Name))
				output.WriteString(fmt.Sprintf("%s\n\n", pkg.Description))

				if len(pkg.UseCase) > 0 {
					output.WriteString("**Use Cases:**\n")
					for _, uc := range pkg.UseCase {
						output.WriteString(fmt.Sprintf("- %s\n", uc))
					}
					output.WriteString("\n")
				}

				output.WriteString("**Installation:**\n")
				output.WriteString("```bash\n")
				output.WriteString(fmt.Sprintf("composer require %s\n", pkg.ComposerName))
				output.WriteString("```\n\n")
			}
		}

		return mcp.NewToolResultText(output.String()), nil
	})

	// Tool 8: get_laravel_package_info
	packageInfoTool := mcp.NewTool("get_laravel_package_info",
		mcp.WithDescription("Provides comprehensive details about a specific Laravel package including installation and use cases.\n\nWhen to use:\n- Learning about a specific package\n- Getting installation instructions\n- Understanding package capabilities\n- Checking maintenance status"),
		mcp.WithString("package_name",
			mcp.Required(),
			mcp.Description("The name of the package (e.g., 'laravel/cashier')"),
		),
	)

	s.mcp.AddTool(packageInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName, err := request.RequireString("package_name")
		if err != nil {
			return mcp.NewToolResultError("package_name is required"), nil
		}

		pkg, err := catalog.GetPackage(packageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Package not found: %s", packageName)), nil
		}

		formatted, err := packages.FormatPackageDetails(pkg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format package: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})

	// Tool 9: get_laravel_package_categories
	categoriesToolNew := mcp.NewTool("get_laravel_package_categories",
		mcp.WithDescription("Lists all packages within a specific functional category.\n\nWhen to use:\n- Exploring packages by category\n- Finding all authentication/payment/testing packages\n- Discovering options in a domain\n- Browsing available solutions"),
		mcp.WithString("category",
			mcp.Required(),
			mcp.Description("The category to filter by (e.g., 'authentication', 'payment')"),
		),
	)

	s.mcp.AddTool(categoriesToolNew, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category, err := request.RequireString("category")
		if err != nil {
			return mcp.NewToolResultError("category is required"), nil
		}

		cat, err := catalog.GetCategory(category)
		if err != nil {
			availableCategories := catalog.ListCategories()
			return mcp.NewToolResultText(fmt.Sprintf(
				"No packages found in category: '%s'.\n\nAvailable categories: %s",
				category,
				strings.Join(availableCategories, ", "),
			)), nil
		}

		formatted, err := packages.FormatCategoryPackages(cat)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format category: %v", err)), nil
		}

		return mcp.NewToolResultText(formatted), nil
	})

	// Tool 10: get_features_for_laravel_package
	featuresTool := mcp.NewTool("get_features_for_laravel_package",
		mcp.WithDescription("Details common implementation features and patterns for a specific package.\n\nWhen to use:\n- Understanding what a package can do\n- Planning implementation\n- Comparing package capabilities\n- Learning package features"),
		mcp.WithString("package",
			mcp.Required(),
			mcp.Description("The Laravel package name (e.g., 'laravel/cashier')"),
		),
	)

	s.mcp.AddTool(featuresTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName, err := request.RequireString("package")
		if err != nil {
			return mcp.NewToolResultError("package is required"), nil
		}

		pkg, err := catalog.GetPackage(packageName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Package not found: %s", packageName)), nil
		}

		features := catalog.GetFeatures(packageName)

		var output strings.Builder
		output.WriteString(fmt.Sprintf("# Features: %s\n\n", pkg.Name))
		output.WriteString(fmt.Sprintf("**Package:** `%s`\n\n", pkg.ComposerName))

		if len(features) == 0 {
			output.WriteString("No specific features documented. Check the package documentation for details.\n")
		} else {
			output.WriteString("## Key Features\n\n")
			for _, feature := range features {
				output.WriteString(fmt.Sprintf("- %s\n", feature))
			}
			output.WriteString("\n")
		}

		if len(pkg.Alternatives) > 0 {
			output.WriteString("## Alternatives\n\n")
			for _, alt := range pkg.Alternatives {
				output.WriteString(fmt.Sprintf("- %s\n", alt))
			}
		}

		return mcp.NewToolResultText(output.String()), nil
	})
}
