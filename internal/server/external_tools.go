package server

import (
	"context"
	"fmt"

	"github.com/izzamoe/laravel-mcp-companion-go/internal/external"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/updater"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input types for external tools
type UpdateDocsInput struct {
	VersionParam string `json:"version_param,omitempty" jsonschema:"Laravel version branch (e.g. '12.x')"`
	Force        *bool  `json:"force,omitempty" jsonschema:"Force update even if already up to date"`
}

type DocsInfoInput struct {
	Version string `json:"version,omitempty" jsonschema:"Specific Laravel version to get info for (e.g. '12.x'). If not provided shows all versions"`
}

type UpdateExternalInput struct {
	Services []string `json:"services,omitempty" jsonschema:"List of services to update (forge vapor envoyer nova). If None updates all"`
	Force    *bool    `json:"force,omitempty" jsonschema:"Force update even if cache is valid"`
}

type SearchExternalInput struct {
	Query    string   `json:"query" jsonschema:"required,Search term to look for"`
	Services []string `json:"services,omitempty" jsonschema:"List of services to search. If None searches all cached services"`
}

type ServiceInfoInput struct {
	Service string `json:"service" jsonschema:"required,Service name (forge vapor envoyer nova)"`
}

// RegisterExternalTools registers update and external resource tools
func (s *Server) RegisterExternalTools(upd *updater.GitHubUpdater, scraper *external.WebScraper) {
	// Tool 11: update_laravel_docs
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "update_laravel_docs",
		Description: "Updates documentation from the official Laravel GitHub repository. Ensures access to the latest documentation changes.\n\nWhen to use:\n- Getting the latest documentation\n- Ensuring documentation is up to date\n- After Laravel version release\n- When documentation seems outdated",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input UpdateDocsInput) (*mcp.CallToolResult, EmptyOutput, error) {
		version := input.VersionParam
		if version == "" {
			version = "12.x"
		}

		force := false
		if input.Force != nil {
			force = *input.Force
		}

		result, err := upd.UpdateDocs(version)
		if err != nil {
			errMsg := fmt.Sprintf("Update failed: %v", err)
			if !force {
				errMsg += ". Try with force=true to force update"
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: errMsg}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	// Tool 12: laravel_docs_info
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "laravel_docs_info",
		Description: "Provides metadata about documentation versions, including last update times and commit information.\n\nWhen to use:\n- Checking documentation freshness\n- Verifying which version is available\n- Getting documentation statistics\n- Planning documentation updates",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input DocsInfoInput) (*mcp.CallToolResult, EmptyOutput, error) {
		info, err := s.docManager.GetInfo(input.Version)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get info: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: info}},
		}, EmptyOutput{}, nil
	})
}

// RegisterExternalServiceTools registers external Laravel service tools (Tools 13-16)
func (s *Server) RegisterExternalServiceTools(externalManager *external.ExternalManager) {
	// Tool 13: update_external_laravel_docs
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "update_external_laravel_docs",
		Description: "Updates documentation for external Laravel services like Forge, Vapor, Envoyer, and Nova.\n\nWhen to use:\n- Getting latest external service docs\n- Setting up deployment workflows\n- Learning about Laravel services\n- Checking service features",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input UpdateExternalInput) (*mcp.CallToolResult, EmptyOutput, error) {
		force := false
		if input.Force != nil {
			force = *input.Force
		}

		result, err := externalManager.UpdateServices(input.Services, force)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Update failed: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	// Tool 14: list_laravel_services
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_laravel_services",
		Description: "Lists all available Laravel services with external documentation support.\n\nWhen to use:\n- Discovering Laravel services\n- Planning service integration\n- Learning about Laravel ecosystem\n- Checking available services",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, EmptyOutput, error) {
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

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	// Tool 15: search_external_laravel_docs
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search_external_laravel_docs",
		Description: "Searches through external Laravel service documentation.\n\nWhen to use:\n- Finding service-specific information\n- Learning about service features\n- Troubleshooting service issues\n- Comparing service capabilities",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input SearchExternalInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		result, err := externalManager.SearchServices(input.Query, input.Services)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Search failed: %v", err)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: result}},
		}, EmptyOutput{}, nil
	})

	// Tool 16: get_laravel_service_info
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_laravel_service_info",
		Description: "Provides detailed information about a specific Laravel service.\n\nWhen to use:\n- Learning about a service\n- Understanding service pricing\n- Checking service requirements\n- Planning service adoption",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input ServiceInfoInput) (*mcp.CallToolResult, EmptyOutput, error) {
		if input.Service == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "service is required"}},
				IsError: true,
			}, EmptyOutput{}, nil
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

		info, ok := serviceInfo[input.Service]
		if !ok {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Service not found: %s. Available services: forge, vapor, envoyer, nova", input.Service)}},
				IsError: true,
			}, EmptyOutput{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: info}},
		}, EmptyOutput{}, nil
	})
}
