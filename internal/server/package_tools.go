package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/packages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input types for package tools
type RecommendPackageInput struct {
	UseCase string `json:"use_case" jsonschema:"required,Description of what the user wants to implement"`
}

type PackageInfoInput struct {
	PackageName string `json:"package_name" jsonschema:"required,The name of the package (e.g. 'laravel/cashier')"`
}

type PackageCategoryInput struct {
	Category string `json:"category" jsonschema:"required,The category to filter by (e.g. 'authentication' 'payment')"`
}

type PackageFeaturesInput struct {
	Package string `json:"package" jsonschema:"required,The Laravel package name (e.g. 'laravel/cashier')"`
}

// RegisterPackageTools registers all 4 package-related MCP tools
func (s *Server) RegisterPackageTools(catalog *packages.Catalog) {
	// Tool 7: get_laravel_package_recommendations
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_laravel_package_recommendations",
		Description: "Intelligently recommends Laravel packages based on described use cases or implementation needs.\n\nWhen to use:\n- Starting a new feature and need package suggestions\n- Looking for solutions to specific problems\n- Comparing available options for a use case\n- Discovering well-maintained packages",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input RecommendPackageInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.UseCase == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "use_case is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		recommendations := catalog.Recommend(input.UseCase, 5)

		var output strings.Builder
		output.WriteString(fmt.Sprintf("# Laravel Packages for: %s\n\n", input.UseCase))

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

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: output.String()}},
		}, EmptyOutput{}, nil
	})

	// Tool 8: get_laravel_package_info
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_laravel_package_info",
		Description: "Provides comprehensive details about a specific Laravel package including installation and use cases.\n\nWhen to use:\n- Learning about a specific package\n- Getting installation instructions\n- Understanding package capabilities\n- Checking maintenance status",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input PackageInfoInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.PackageName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "package_name is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		pkg, err := catalog.GetPackage(input.PackageName)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Package not found: %s", input.PackageName)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		formatted, err := packages.FormatPackageDetails(pkg)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to format package: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: formatted}},
		}, EmptyOutput{}, nil
	})

	// Tool 9: get_laravel_package_categories
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_laravel_package_categories",
		Description: "Lists all packages within a specific functional category.\n\nWhen to use:\n- Exploring packages by category\n- Finding all authentication/payment/testing packages\n- Discovering options in a domain\n- Browsing available solutions",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input PackageCategoryInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Category == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "category is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		cat, err := catalog.GetCategory(input.Category)
		if err != nil {
			availableCategories := catalog.ListCategories()
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(
					"No packages found in category: '%s'.\n\nAvailable categories: %s",
					input.Category,
					strings.Join(availableCategories, ", "),
				)}},
			}, EmptyOutput{}, nil
		}

		formatted, err := packages.FormatCategoryPackages(cat)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to format category: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: formatted}},
		}, EmptyOutput{}, nil
	})

	// Tool 10: get_features_for_laravel_package
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_features_for_laravel_package",
		Description: "Details common implementation features and patterns for a specific package.\n\nWhen to use:\n- Understanding what a package can do\n- Planning implementation\n- Comparing package capabilities\n- Learning package features",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input PackageFeaturesInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Package == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "package is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		pkg, err := catalog.GetPackage(input.Package)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Package not found: %s", input.Package)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		features := catalog.GetFeatures(input.Package)

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

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: output.String()}},
		}, EmptyOutput{}, nil
	})
}
