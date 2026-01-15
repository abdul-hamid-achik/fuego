package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List all registered routes and pages",
	Long: `Display all routes and pages discovered in the app directory.

This command scans the app/ directory and displays:
- API routes (route.go files) with their HTTP methods and patterns
- Pages (page.templ files) with their URL patterns and associated layouts

Examples:
  nexo routes
  nexo routes --json
  nexo routes --app-dir custom/app`,
	Run: runRoutes,
}

var (
	routesAppDir string
)

func init() {
	routesCmd.Flags().StringVarP(&routesAppDir, "app-dir", "d", "app", "App directory to scan")
}

func runRoutes(cmd *cobra.Command, args []string) {
	// Check if app directory exists
	if _, err := os.Stat(routesAppDir); os.IsNotExist(err) {
		if jsonOutput {
			printSuccess(RoutesOutput{
				Routes:      []RouteOutput{},
				TotalRoutes: 0,
			})
		} else {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("\n  %s No app directory found at %s\n\n", yellow("Warning:"), routesAppDir)
		}
		return
	}

	// Scan for routes
	scanner := nexo.NewScanner(routesAppDir)

	// Check for proxy
	proxyInfo, proxyErr := scanner.ScanProxyInfo()

	// Scan for middleware
	middlewares, mwErr := scanner.ScanMiddlewareInfo()

	// Scan for routes
	routes, routeErr := scanner.ScanRouteInfo()
	if routeErr != nil {
		if jsonOutput {
			printJSONError(routeErr)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("  %s Failed to scan routes: %v\n", red("Error:"), routeErr)
		}
		os.Exit(1)
	}

	// Sort routes by pattern
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Pattern != routes[j].Pattern {
			return routes[i].Pattern < routes[j].Pattern
		}
		return routes[i].Method < routes[j].Method
	})

	// Scan for pages
	pages, pageErr := scanner.ScanPageInfo()
	if pageErr != nil {
		if jsonOutput {
			printJSONError(pageErr)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("  %s Failed to scan pages: %v\n", red("Error:"), pageErr)
		}
		os.Exit(1)
	}

	// Scan for layouts
	layouts, layoutErr := scanner.ScanLayoutInfo()
	if layoutErr != nil {
		if jsonOutput {
			printJSONError(layoutErr)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("  %s Failed to scan layouts: %v\n", red("Error:"), layoutErr)
		}
		os.Exit(1)
	}

	// JSON output mode
	if jsonOutput {
		output := RoutesOutput{
			Routes:      make([]RouteOutput, 0, len(routes)),
			Pages:       make([]PageOutput, 0, len(pages)),
			TotalRoutes: len(routes),
			TotalPages:  len(pages),
		}

		// Add proxy info
		if proxyErr == nil && proxyInfo != nil && proxyInfo.HasProxy {
			output.Proxy = &ProxyOutput{
				Enabled:  true,
				File:     proxyInfo.FilePath,
				Matchers: proxyInfo.Matchers,
			}
		}

		// Add middleware info
		if mwErr == nil && len(middlewares) > 0 {
			output.Middleware = make([]MiddlewareOutput, 0, len(middlewares))
			for _, mw := range middlewares {
				path := mw.Path
				if path == "" {
					path = "/"
				}
				output.Middleware = append(output.Middleware, MiddlewareOutput{
					Path: path,
					File: mw.FilePath,
				})
			}
		}

		// Add routes
		for _, r := range routes {
			output.Routes = append(output.Routes, RouteOutput{
				Method:   r.Method,
				Pattern:  r.Pattern,
				File:     r.FilePath,
				Priority: r.Priority,
			})
		}

		// Add pages
		for _, p := range pages {
			output.Pages = append(output.Pages, PageOutput{
				Pattern: p.Pattern,
				File:    p.FilePath,
				Title:   p.Title,
				Layout:  findLayoutForPage(p.Pattern, layouts),
			})
		}

		printSuccess(output)
		return
	}

	// Text output mode
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Printf("\n  %s Routes\n\n", cyan("Nexo"))

	// Show proxy info
	if proxyErr != nil {
		fmt.Printf("  %s Failed to scan proxy: %v\n", yellow("Warning:"), proxyErr)
	} else if proxyInfo != nil && proxyInfo.HasProxy {
		fmt.Printf("  %s Proxy enabled\n", magenta("PROXY"))
		if len(proxyInfo.Matchers) > 0 {
			fmt.Printf("        Matchers: %v\n", proxyInfo.Matchers)
		} else {
			fmt.Printf("        Matchers: all paths\n")
		}
		fmt.Printf("        File: %s\n\n", dim(proxyInfo.FilePath))
	}

	// Show middleware info
	if mwErr != nil {
		fmt.Printf("  %s Failed to scan middleware: %v\n", yellow("Warning:"), mwErr)
	} else if len(middlewares) > 0 {
		fmt.Printf("  %s\n", cyan("Middleware:"))
		for _, mw := range middlewares {
			path := mw.Path
			if path == "" {
				path = "/"
			}
			fmt.Printf("        %s  %s\n", fmt.Sprintf("%-30s", path), dim(mw.FilePath))
		}
		fmt.Printf("\n")
	}

	// Method colors
	methodColor := func(method string) string {
		switch method {
		case "GET":
			return green(fmt.Sprintf("%-7s", method))
		case "POST":
			return yellow(fmt.Sprintf("%-7s", method))
		case "PUT":
			return cyan(fmt.Sprintf("%-7s", method))
		case "PATCH":
			return color.MagentaString("%-7s", method)
		case "DELETE":
			return red(fmt.Sprintf("%-7s", method))
		default:
			return fmt.Sprintf("%-7s", method)
		}
	}

	// Print API routes section
	if len(routes) > 0 {
		fmt.Printf("  %s\n\n", cyan("API Routes:"))
		for _, route := range routes {
			fmt.Printf("  %s %s  %s\n",
				methodColor(route.Method),
				fmt.Sprintf("%-30s", route.Pattern),
				dim(route.FilePath),
			)
		}
	}

	// Print pages section (only if pages exist)
	if len(pages) > 0 {
		if len(routes) > 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("  %s\n\n", cyan("Pages:"))
		for _, page := range pages {
			layoutInfo := ""
			if layout := findLayoutForPage(page.Pattern, layouts); layout != "" {
				// Extract just the directory name from the layout path
				layoutDir := filepath.Base(filepath.Dir(layout))
				layoutInfo = dim(fmt.Sprintf(" [layout: %s]", layoutDir))
			}
			fmt.Printf("  %s %s  %s%s\n",
				green("GET    "),
				fmt.Sprintf("%-30s", page.Pattern),
				dim(page.FilePath),
				layoutInfo,
			)
		}
	}

	// Show warning if no routes and no pages
	if len(routes) == 0 && len(pages) == 0 {
		fmt.Printf("  %s No routes or pages found\n\n", yellow("Warning:"))
		fmt.Printf("  Create an API route by adding a route.go file:\n")
		fmt.Printf("    %s/api/health/route.go\n\n", routesAppDir)
		fmt.Printf("  Or create a page by adding a page.templ file:\n")
		fmt.Printf("    %s/page.templ\n\n", routesAppDir)
		return
	}

	fmt.Printf("\n  Total: %d API routes, %d pages\n\n", len(routes), len(pages))
}

// findLayoutForPage returns the layout file path that applies to a page pattern.
// It finds the most specific layout that matches the page path.
func findLayoutForPage(pagePattern string, layouts []nexo.LayoutInfo) string {
	var bestMatch string
	var bestMatchLen int

	for _, layout := range layouts {
		prefix := layout.PathPrefix
		// Check if the page pattern starts with the layout prefix
		// or if the layout is at root level
		if strings.HasPrefix(pagePattern, prefix) || prefix == "/" {
			// Prefer more specific matches (longer prefix)
			if len(prefix) > bestMatchLen {
				bestMatch = layout.FilePath
				bestMatchLen = len(prefix)
			}
		}
	}
	return bestMatch
}
