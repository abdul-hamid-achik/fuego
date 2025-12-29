package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List all registered routes",
	Long: `Display all routes discovered in the app directory.

This command scans the app/ directory and displays all route.go files
with their HTTP methods and patterns.

Examples:
  fuego routes
  fuego routes --json
  fuego routes --app-dir custom/app`,
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
				Routes: []RouteOutput{},
				Total:  0,
			})
		} else {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("\n  %s No app directory found at %s\n\n", yellow("Warning:"), routesAppDir)
		}
		return
	}

	// Scan for routes
	scanner := fuego.NewScanner(routesAppDir)

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

	// JSON output mode
	if jsonOutput {
		output := RoutesOutput{
			Routes: make([]RouteOutput, 0, len(routes)),
			Total:  len(routes),
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

	fmt.Printf("\n  %s Routes\n\n", cyan("Fuego"))

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

	if len(routes) == 0 {
		fmt.Printf("  %s No routes found\n\n", yellow("Warning:"))
		fmt.Printf("  Create a route by adding a route.go file:\n")
		fmt.Printf("    %s/api/health/route.go\n\n", routesAppDir)
		return
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

	// Print routes
	for _, route := range routes {
		fmt.Printf("  %s %s  %s\n",
			methodColor(route.Method),
			fmt.Sprintf("%-30s", route.Pattern),
			dim(route.FilePath),
		)
	}

	fmt.Printf("\n  Total: %d routes\n\n", len(routes))
}
