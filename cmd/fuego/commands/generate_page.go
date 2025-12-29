package commands

import (
	"fmt"

	"github.com/abdul-hamid-achik/fuego/pkg/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generatePageCmd = &cobra.Command{
	Use:   "page <path>",
	Short: "Generate a page template",
	Long: `Generate a page.templ file for rendering HTML pages.

Pages use the templ templating language and inherit from layouts.
Use --with-layout to also generate a layout.templ for that section.

Examples:
  fuego generate page dashboard
  fuego generate page admin/settings
  fuego generate page blog/posts --with-layout`,
	Args: cobra.ExactArgs(1),
	Run:  runGeneratePage,
}

var (
	pageWithLayout bool
	pageAppDir     string
)

func init() {
	generatePageCmd.Flags().BoolVar(&pageWithLayout, "with-layout", false, "Also generate a layout.templ for this section")
	generatePageCmd.Flags().StringVarP(&pageAppDir, "app-dir", "d", "app", "App directory")
	generateCmd.AddCommand(generatePageCmd)
}

func runGeneratePage(cmd *cobra.Command, args []string) {
	path := args[0]

	result, err := generator.GeneratePage(generator.PageConfig{
		Path:       path,
		AppDir:     pageAppDir,
		WithLayout: pageWithLayout,
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
			Command: "generate page",
			Path:    path,
			Files:   result.Files,
			Pattern: result.Pattern,
		})
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\n  %s Generated page\n\n", green("✓"))
	for _, f := range result.Files {
		fmt.Printf("    Created: %s\n", cyan(f))
	}
	fmt.Printf("    URL: %s\n\n", result.Pattern)

	if pageWithLayout {
		fmt.Printf("    Note: Layout created. Pages in this directory will use it.\n\n")
	}
}
