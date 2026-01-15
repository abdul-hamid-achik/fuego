package commands

import (
	"fmt"

	"github.com/abdul-hamid-achik/nexo/pkg/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generateProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Generate a proxy file",
	Long: `Generate a proxy.go file for request interception.

Available templates:
  blank        - Empty proxy (default)
  auth-check   - Authentication checking before routing
  rate-limit   - Simple IP-based rate limiting
  maintenance  - Maintenance mode with allowed IPs
  redirect-www - WWW/non-WWW redirect handling

The proxy runs before route matching and can:
  - Rewrite URLs (A/B testing, feature flags)
  - Redirect requests
  - Return early responses (auth, rate limiting)
  - Add request headers

Examples:
  nexo generate proxy --template auth-check
  nexo generate proxy --template rate-limit
  nexo generate proxy --template maintenance`,
	Run: runGenerateProxy,
}

var (
	proxyTemplate string
	proxyAppDir   string
)

func init() {
	generateProxyCmd.Flags().StringVarP(&proxyTemplate, "template", "t", "blank", "Template: blank, auth-check, rate-limit, maintenance, redirect-www")
	generateProxyCmd.Flags().StringVarP(&proxyAppDir, "app-dir", "d", "app", "App directory")
	generateCmd.AddCommand(generateProxyCmd)
}

func runGenerateProxy(cmd *cobra.Command, args []string) {
	result, err := generator.GenerateProxy(generator.ProxyConfig{
		Template: proxyTemplate,
		AppDir:   proxyAppDir,
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
			Command: "generate proxy",
			Files:   result.Files,
		})
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\n  %s Generated proxy\n\n", green("✓"))
	for _, f := range result.Files {
		fmt.Printf("    Created: %s\n", cyan(f))
	}
	fmt.Printf("    Template: %s\n\n", proxyTemplate)
}
