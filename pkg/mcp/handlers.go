package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
	"github.com/abdul-hamid-achik/fuego/pkg/generator"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) handleNew(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("name is required"), nil
	}

	apiOnly := req.GetBool("api_only", false)
	withProxy := req.GetBool("with_proxy", false)

	// Build command
	args := []string{"new", name, "--json"}
	if apiOnly {
		args = append(args, "--api-only")
	}
	if withProxy {
		args = append(args, "--with-proxy")
	}

	cmd := exec.CommandContext(ctx, "fuego", args...)
	cmd.Dir = s.workdir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create project: %s", string(output))), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleGenerateRoute(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := req.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError("path is required"), nil
	}

	methodsStr := req.GetString("methods", "GET")
	methods := strings.Split(strings.ToUpper(methodsStr), ",")
	for i, m := range methods {
		methods[i] = strings.TrimSpace(m)
	}

	appDir := filepath.Join(s.workdir, "app")
	result, err := generator.GenerateRoute(generator.RouteConfig{
		Path:    path,
		Methods: methods,
		AppDir:  appDir,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	output, _ := json.MarshalIndent(map[string]any{
		"success": true,
		"files":   result.Files,
		"pattern": result.Pattern,
		"methods": methods,
	}, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleGenerateMiddleware(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("name is required"), nil
	}

	path := req.GetString("path", "")
	template := req.GetString("template", "blank")

	appDir := filepath.Join(s.workdir, "app")
	result, err := generator.GenerateMiddleware(generator.MiddlewareConfig{
		Name:     name,
		Path:     path,
		Template: template,
		AppDir:   appDir,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	output, _ := json.MarshalIndent(map[string]any{
		"success":  true,
		"files":    result.Files,
		"path":     path,
		"template": template,
	}, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleGenerateProxy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	template, err := req.RequireString("template")
	if err != nil {
		return mcp.NewToolResultError("template is required"), nil
	}

	appDir := filepath.Join(s.workdir, "app")
	result, err := generator.GenerateProxy(generator.ProxyConfig{
		Template: template,
		AppDir:   appDir,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	output, _ := json.MarshalIndent(map[string]any{
		"success":  true,
		"files":    result.Files,
		"template": template,
	}, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleGeneratePage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := req.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError("path is required"), nil
	}

	withLayout := req.GetBool("with_layout", false)

	appDir := filepath.Join(s.workdir, "app")
	result, err := generator.GeneratePage(generator.PageConfig{
		Path:       path,
		AppDir:     appDir,
		WithLayout: withLayout,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	output, _ := json.MarshalIndent(map[string]any{
		"success": true,
		"files":   result.Files,
		"pattern": result.Pattern,
	}, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleListRoutes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	appDir := filepath.Join(s.workdir, "app")
	scanner := fuego.NewScanner(appDir)

	routes, err := scanner.ScanRouteInfo()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	middlewares, _ := scanner.ScanMiddlewareInfo()
	proxyInfo, _ := scanner.ScanProxyInfo()

	result := map[string]any{
		"routes":     routes,
		"middleware": middlewares,
		"total":      len(routes),
	}

	if proxyInfo != nil && proxyInfo.HasProxy {
		result["proxy"] = map[string]any{
			"enabled":  true,
			"file":     proxyInfo.FilePath,
			"matchers": proxyInfo.Matchers,
		}
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := map[string]any{
		"workdir": s.workdir,
	}

	// Check for fuego.yaml
	configPath := filepath.Join(s.workdir, "fuego.yaml")
	if _, err := os.Stat(configPath); err == nil {
		info["has_config"] = true
		info["config_path"] = configPath
	} else {
		info["has_config"] = false
	}

	// Check for go.mod
	goModPath := filepath.Join(s.workdir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		info["has_go_mod"] = true
	} else {
		info["has_go_mod"] = false
	}

	// Scan routes
	appDir := filepath.Join(s.workdir, "app")
	if _, err := os.Stat(appDir); err == nil {
		info["has_app_dir"] = true
		scanner := fuego.NewScanner(appDir)

		routes, _ := scanner.ScanRouteInfo()
		middlewares, _ := scanner.ScanMiddlewareInfo()
		proxyInfo, _ := scanner.ScanProxyInfo()

		info["routes"] = routes
		info["middleware"] = middlewares
		info["route_count"] = len(routes)

		if proxyInfo != nil && proxyInfo.HasProxy {
			info["proxy"] = map[string]any{
				"enabled":  true,
				"file":     proxyInfo.FilePath,
				"matchers": proxyInfo.Matchers,
			}
		}
	} else {
		info["has_app_dir"] = false
	}

	output, _ := json.MarshalIndent(info, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

func (s *Server) handleValidate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var issues []string
	var warnings []string

	// Check app directory
	appDir := filepath.Join(s.workdir, "app")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		issues = append(issues, "app/ directory not found")
	}

	// Check go.mod
	goModPath := filepath.Join(s.workdir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		issues = append(issues, "go.mod not found - not a Go project")
	}

	// Check main.go
	mainPath := filepath.Join(s.workdir, "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		warnings = append(warnings, "main.go not found in project root")
	}

	// Scan for route issues
	var routeCount int
	if _, err := os.Stat(appDir); err == nil {
		scanner := fuego.NewScanner(appDir)
		scanner.SetVerbose(false)

		routes, err := scanner.ScanRouteInfo()
		if err != nil {
			issues = append(issues, "Failed to scan routes: "+err.Error())
		} else {
			routeCount = len(routes)
			if routeCount == 0 {
				warnings = append(warnings, "No routes found in app/ directory")
			}
		}

		// Check middleware
		middlewares, err := scanner.ScanMiddlewareInfo()
		if err != nil {
			warnings = append(warnings, "Failed to scan middleware: "+err.Error())
		}
		_ = middlewares

		// Check proxy
		proxyInfo, err := scanner.ScanProxyInfo()
		if err != nil {
			warnings = append(warnings, "Failed to scan proxy: "+err.Error())
		}
		_ = proxyInfo
	}

	result := map[string]any{
		"valid":       len(issues) == 0,
		"issues":      issues,
		"warnings":    warnings,
		"route_count": routeCount,
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}
