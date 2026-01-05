// Package mcp provides an MCP (Model Context Protocol) server for Nexo.
package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server is the Nexo MCP server.
type Server struct {
	mcpServer *server.MCPServer
	workdir   string
}

// NewServer creates a new Nexo MCP server.
func NewServer(workdir string) *Server {
	s := server.NewMCPServer(
		"nexo",
		"0.2.1",
		server.WithToolCapabilities(true),
	)

	srv := &Server{
		mcpServer: s,
		workdir:   workdir,
	}

	srv.registerTools()
	return srv
}

func (s *Server) registerTools() {
	// nexo_new - Create new project
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_new",
			mcp.WithDescription("Create a new Nexo project"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
			mcp.WithBoolean("api_only", mcp.Description("Create API-only project without templ templates")),
			mcp.WithBoolean("with_proxy", mcp.Description("Include proxy.go example")),
		),
		s.handleNew,
	)

	// nexo_generate_route - Generate a route
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_generate_route",
			mcp.WithDescription("Generate a new route file with handler functions"),
			mcp.WithString("path", mcp.Required(), mcp.Description("Route path (e.g., 'users/[id]', 'posts/[...slug]')")),
			mcp.WithString("methods", mcp.Description("HTTP methods comma-separated (default: GET)")),
		),
		s.handleGenerateRoute,
	)

	// nexo_generate_middleware - Generate middleware
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_generate_middleware",
			mcp.WithDescription("Generate a middleware file"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Middleware name")),
			mcp.WithString("path", mcp.Description("Path prefix (e.g., 'api/protected')")),
			mcp.WithString("template", mcp.Description("Template: blank, auth, logging, timing, cors")),
		),
		s.handleGenerateMiddleware,
	)

	// nexo_generate_proxy - Generate proxy
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_generate_proxy",
			mcp.WithDescription("Generate a proxy.go file for request interception"),
			mcp.WithString("template", mcp.Required(), mcp.Description("Template: blank, auth-check, rate-limit, maintenance, redirect-www")),
		),
		s.handleGenerateProxy,
	)

	// nexo_generate_page - Generate page
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_generate_page",
			mcp.WithDescription("Generate a page.templ file"),
			mcp.WithString("path", mcp.Required(), mcp.Description("Page path (e.g., 'dashboard', 'admin/settings')")),
			mcp.WithBoolean("with_layout", mcp.Description("Also generate a layout.templ for this section")),
		),
		s.handleGeneratePage,
	)

	// nexo_list_routes - List all routes
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_list_routes",
			mcp.WithDescription("List all routes, middleware, and proxy in the project"),
		),
		s.handleListRoutes,
	)

	// nexo_info - Get project info
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_info",
			mcp.WithDescription("Get comprehensive project information"),
		),
		s.handleInfo,
	)

	// nexo_validate - Validate project
	s.mcpServer.AddTool(
		mcp.NewTool("nexo_validate",
			mcp.WithDescription("Validate project structure and handler signatures"),
		),
		s.handleValidate,
	)
}

// ServeStdio starts the MCP server over stdio.
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.mcpServer)
}
