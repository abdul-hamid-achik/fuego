package commands

import (
	"fmt"

	"github.com/abdul-hamid-achik/nexo/pkg/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generateMiddlewareCmd = &cobra.Command{
	Use:   "middleware <name>",
	Short: "Generate middleware",
	Long: `Generate a middleware file with common patterns.

Available templates:
  blank   - Empty middleware (default)
  auth    - Authentication checking
  logging - Request/response logging
  timing  - Response time headers
  cors    - CORS headers

Examples:
  nexo generate middleware auth --path api/protected
  nexo generate middleware logging --path api --template logging
  nexo generate middleware cors --template cors`,
	Args: cobra.ExactArgs(1),
	Run:  runGenerateMiddleware,
}

var (
	middlewarePath     string
	middlewareTemplate string
	middlewareAppDir   string
)

func init() {
	generateMiddlewareCmd.Flags().StringVarP(&middlewarePath, "path", "p", "", "Path prefix (e.g., api/protected)")
	generateMiddlewareCmd.Flags().StringVarP(&middlewareTemplate, "template", "t", "blank", "Template: blank, auth, logging, timing, cors")
	generateMiddlewareCmd.Flags().StringVarP(&middlewareAppDir, "app-dir", "d", "app", "App directory")
	generateCmd.AddCommand(generateMiddlewareCmd)
}

func runGenerateMiddleware(cmd *cobra.Command, args []string) {
	name := args[0]

	result, err := generator.GenerateMiddleware(generator.MiddlewareConfig{
		Name:     name,
		Path:     middlewarePath,
		Template: middlewareTemplate,
		AppDir:   middlewareAppDir,
	})

	if err != nil {
		if jsonOutput {
			printJSONError(err)
		} else {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("  %s %v\n", red("Error:"), err)
		}
		return
	}

	if jsonOutput {
		printSuccess(GenerateOutput{
			Command: "generate middleware",
			Path:    middlewarePath,
			Files:   result.Files,
		})
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\n  %s Generated middleware\n\n", green("✓"))
	for _, f := range result.Files {
		fmt.Printf("    Created: %s\n", cyan(f))
	}
	if middlewarePath != "" {
		fmt.Printf("    Applies to: /%s/*\n", middlewarePath)
	} else {
		fmt.Printf("    Applies to: all routes\n")
	}
	fmt.Printf("    Template: %s\n\n", middlewareTemplate)
}
