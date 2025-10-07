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
	// Tool 11: update_laravel_docs
	updateDocsTool := mcp.NewTool("update_laravel_docs",
		mcp.WithDescription("Updates documentation from the official Laravel GitHub repository. Ensures access to the latest documentation changes.\n\nWhen to use:\n- Getting the latest documentation\n- Ensuring documentation is up to date\n- After Laravel version release\n- When documentation seems outdated"),
		mcp.WithString("version_param",
			mcp.Description("Laravel version branch (e.g., '12.x')"),
		),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Force update even if already up to date"),
		),
	)

	s.mcp.AddTool(updateDocsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := mcp.ParseString(request, "version_param", "")
		force := mcp.ParseBoolean(request, "force", false)

		if version == "" {
			version = "12.x" // default
		}

		result, err := upd.UpdateDocs(version)
		if err != nil {
			if !force {
				return mcp.NewToolResultError(fmt.Sprintf("Update failed: %v. Try with force=true to force update", err)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Update failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Tool 12: laravel_docs_info
	docsInfoTool := mcp.NewTool("laravel_docs_info",
		mcp.WithDescription("Provides metadata about documentation versions, including last update times and commit information.\n\nWhen to use:\n- Checking documentation freshness\n- Verifying which version is available\n- Getting documentation statistics\n- Planning documentation updates"),
		mcp.WithString("version",
			mcp.Description("Specific Laravel version to get info for (e.g., '12.x'). If not provided, shows all versions"),
		),
	)

	s.mcp.AddTool(docsInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		version := mcp.ParseString(request, "version", "")

		info, err := s.docManager.GetInfo(version)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get info: %v", err)), nil
		}

		return mcp.NewToolResultText(info), nil
	})
}

// RegisterExternalServiceTools registers external Laravel service tools (Tools 13-16)
func (s *Server) RegisterExternalServiceTools(scraper *external.WebScraper) {
	// Tool 13: update_external_laravel_docs
	updateExternalTool := mcp.NewTool("update_external_laravel_docs",
		mcp.WithDescription("Updates documentation for external Laravel services like Forge, Vapor, Envoyer, and Nova.\n\nWhen to use:\n- Getting latest external service docs\n- Setting up deployment workflows\n- Learning about Laravel services\n- Checking service features"),
		mcp.WithArray("services",
			mcp.Description("List of services to update (forge, vapor, envoyer, nova). If None, updates all"),
			mcp.WithStringItems(),
		),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Force update even if cache is valid"),
		),
	)

	s.mcp.AddTool(updateExternalTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Parse services array
		servicesMap := mcp.ParseStringMap(request, "services", nil)
		var services []string
		if servicesMap != nil {
			for _, v := range servicesMap {
				if str, ok := v.(string); ok {
					services = append(services, str)
				}
			}
		}

		force := mcp.ParseBoolean(request, "force", false)

		// TODO: Implement external manager
		_ = services
		_ = force

		return mcp.NewToolResultText("External service documentation update is not yet implemented. Available services: forge, vapor, envoyer, nova"), nil
	})

	// Tool 14: list_laravel_services
	listServicesTool := mcp.NewTool("list_laravel_services",
		mcp.WithDescription("Lists all available Laravel services with external documentation support.\n\nWhen to use:\n- Discovering Laravel services\n- Planning service integration\n- Learning about Laravel ecosystem\n- Checking available services"),
	)

	s.mcp.AddTool(listServicesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result := `# Available Laravel Services

## 1. Laravel Forge
Automated server management and deployment platform for Laravel applications.
**URL:** https://forge.laravel.com

## 2. Laravel Vapor
Serverless deployment platform for Laravel on AWS Lambda.
**URL:** https://vapor.laravel.com

## 3. Laravel Envoyer
Zero-downtime deployment platform for PHP applications.
**URL:** https://envoyer.io

## 4. Laravel Nova
Administration panel for Laravel applications.
**URL:** https://nova.laravel.com`

		return mcp.NewToolResultText(result), nil
	})

	// Tool 15: search_external_laravel_docs
	searchExternalTool := mcp.NewTool("search_external_laravel_docs",
		mcp.WithDescription("Searches through external Laravel service documentation.\n\nWhen to use:\n- Finding service-specific information\n- Learning about service features\n- Troubleshooting service issues\n- Comparing service capabilities"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search term to look for"),
		),
		mcp.WithArray("services",
			mcp.Description("List of services to search. If None, searches all cached services"),
			mcp.WithStringItems(),
		),
	)

	s.mcp.AddTool(searchExternalTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("query is required"), nil
		}

		// Parse services
		servicesMap := mcp.ParseStringMap(request, "services", nil)
		var services []string
		if servicesMap != nil {
			for _, v := range servicesMap {
				if str, ok := v.(string); ok {
					services = append(services, str)
				}
			}
		}

		// TODO: Implement external search
		_ = services

		return mcp.NewToolResultText(fmt.Sprintf("External service search for '%s' is not yet implemented. Try using the main documentation search or visit the service websites directly.", query)), nil
	})

	// Tool 16: get_laravel_service_info
	serviceInfoTool := mcp.NewTool("get_laravel_service_info",
		mcp.WithDescription("Provides detailed information about a specific Laravel service.\n\nWhen to use:\n- Learning about a service\n- Understanding service pricing\n- Checking service requirements\n- Planning service adoption"),
		mcp.WithString("service",
			mcp.Required(),
			mcp.Description("Service name (forge, vapor, envoyer, nova)"),
		),
	)

	s.mcp.AddTool(serviceInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service, err := request.RequireString("service")
		if err != nil {
			return mcp.NewToolResultError("service is required"), nil
		}

		serviceInfo := map[string]string{
			"forge": `# Laravel Forge

**Type:** Server Management & Deployment Platform

**Description:**
Laravel Forge is a server management and deployment platform that makes it easy to deploy and manage your Laravel applications on any cloud provider.

**Key Features:**
- Automated server provisioning
- Zero-downtime deployments
- SSL certificate management
- Database management
- Queue worker monitoring
- Scheduled job management

**Website:** https://forge.laravel.com
**Documentation:** https://forge.laravel.com/docs`,

			"vapor": `# Laravel Vapor

**Type:** Serverless Deployment Platform

**Description:**
Laravel Vapor is an auto-scaling, serverless deployment platform for Laravel, powered by AWS Lambda.

**Key Features:**
- Serverless infrastructure
- Auto-scaling capabilities
- Global CDN integration
- Database management
- Cache and queue services
- Environment management

**Website:** https://vapor.laravel.com
**Documentation:** https://docs.vapor.build`,

			"envoyer": `# Laravel Envoyer

**Type:** Zero-Downtime Deployment

**Description:**
Envoyer is a zero-downtime deployment platform for PHP applications, including Laravel.

**Key Features:**
- Zero-downtime deployments
- Health checks
- Deployment hooks
- Team collaboration
- Deployment monitoring
- Rollback capabilities

**Website:** https://envoyer.io
**Documentation:** https://envoyer.io/docs`,

			"nova": `# Laravel Nova

**Type:** Administration Panel

**Description:**
Laravel Nova is a beautifully-designed administration panel for Laravel applications.

**Key Features:**
- Resource management
- Custom metrics and cards
- Action execution
- Advanced filters
- Relationship management
- Authorization

**Website:** https://nova.laravel.com
**Documentation:** https://nova.laravel.com/docs`,
		}

		info, ok := serviceInfo[service]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Service not found: %s. Available services: forge, vapor, envoyer, nova", service)), nil
		}

		return mcp.NewToolResultText(info), nil
	})
}
