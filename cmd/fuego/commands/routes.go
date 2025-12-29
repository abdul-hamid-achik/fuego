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

Example:
  fuego routes`,
	Run: runRoutes,
}

var (
	routesAppDir string
)

func init() {
	routesCmd.Flags().StringVarP(&routesAppDir, "app-dir", "d", "app", "App directory to scan")
}

func runRoutes(cmd *cobra.Command, args []string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	fmt.Printf("\n  %s Routes\n\n", cyan("Fuego"))

	// Check if app directory exists
	if _, err := os.Stat(routesAppDir); os.IsNotExist(err) {
		fmt.Printf("  %s No app directory found at %s\n\n", yellow("Warning:"), routesAppDir)
		return
	}

	// Scan for routes
	scanner := fuego.NewScanner(routesAppDir)
	routes, err := scanner.ScanRouteInfo()
	if err != nil {
		fmt.Printf("  %s Failed to scan routes: %v\n", red("Error:"), err)
		os.Exit(1)
	}

	if len(routes) == 0 {
		fmt.Printf("  %s No routes found\n\n", yellow("Warning:"))
		fmt.Printf("  Create a route by adding a route.go file:\n")
		fmt.Printf("    %s/api/health/route.go\n\n", routesAppDir)
		return
	}

	// Sort routes by pattern
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Pattern != routes[j].Pattern {
			return routes[i].Pattern < routes[j].Pattern
		}
		return routes[i].Method < routes[j].Method
	})

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
