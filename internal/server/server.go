package server

import (
	"github.com/izzamoe/laravel-mcp-companion-go/internal/docs"
	"github.com/izzamoe/laravel-mcp-companion-go/internal/external"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps MCP server with documentation manager
type Server struct {
	mcp             *mcp.Server
	docManager      *docs.Manager
	externalManager *external.ExternalManager
}

// NewServer creates a new server instance
func NewServer(docManager *docs.Manager) *Server {
	impl := &mcp.Implementation{
		Name:    "Laravel MCP Companion",
		Version: "1.0.0",
	}

	opts := &mcp.ServerOptions{
		Instructions: "Laravel documentation and package recommendations for AI assistants",
		HasTools:     true,
	}

	mcpServer := mcp.NewServer(impl, opts)

	return &Server{
		mcp:        mcpServer,
		docManager: docManager,
	}
}

// SetExternalManager sets the external manager for the server
func (s *Server) SetExternalManager(em *external.ExternalManager) {
	s.externalManager = em
}

// GetMCPServer returns the underlying MCP server
func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcp
}
