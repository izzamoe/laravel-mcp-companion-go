package server

import (
	"github.com/izzamoe/laravel-mcp-companion-go/internal/docs"
	"github.com/mark3labs/mcp-go/server"
)

// Server wraps MCP server with documentation manager
type Server struct {
	mcp        *server.MCPServer
	docManager *docs.Manager
}

// NewServer creates a new server instance
func NewServer(docManager *docs.Manager) *Server {
	mcpServer := server.NewMCPServer(
		"Laravel MCP Companion",
		"1.0.0",
		server.WithInstructions("Laravel documentation and package recommendations for AI assistants"),
	)

	return &Server{
		mcp:        mcpServer,
		docManager: docManager,
	}
}

// GetMCPServer returns the underlying MCP server for ServeStdio
func (s *Server) GetMCPServer() *server.MCPServer {
	return s.mcp
}
