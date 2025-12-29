package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application for production",
	Long: `Build the application as an optimized production binary.

This command:
  1. Runs templ generate (if .templ files exist)
  2. Builds an optimized Go binary with ldflags

Example:
  fuego build
  fuego build --output ./bin/myapp
  fuego build --os linux --arch amd64`,
	Run: runBuild,
}

var (
	buildOutput string
	buildOS     string
	buildArch   string
)

func init() {
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "Output binary path (default: ./bin/<project-name>)")
	buildCmd.Flags().StringVar(&buildOS, "os", "", "Target OS (linux, darwin, windows)")
	buildCmd.Flags().StringVar(&buildArch, "arch", "", "Target architecture (amd64, arm64)")
}

func runBuild(cmd *cobra.Command, args []string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("\n  %s Production Build\n\n", cyan("Fuego"))

	// Check for main.go
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		fmt.Printf("  %s No main.go found in current directory\n", red("Error:"))
		os.Exit(1)
	}

	// Determine output path
	if buildOutput == "" {
		// Use current directory name as binary name
		cwd, _ := os.Getwd()
		projectName := filepath.Base(cwd)
		buildOutput = filepath.Join("bin", projectName)
	}

	// Add .exe extension on Windows
	targetOS := buildOS
	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	if targetOS == "windows" && !strings.HasSuffix(buildOutput, ".exe") {
		buildOutput += ".exe"
	}

	// Create bin directory
	binDir := filepath.Dir(buildOutput)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("  %s Failed to create output directory: %v\n", red("Error:"), err)
		os.Exit(1)
	}

	// Check for templ files and run templ generate
	hasTemplFiles := false
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(path, ".templ") {
			hasTemplFiles = true
			return filepath.SkipAll
		}
		return nil
	})

	if hasTemplFiles {
		fmt.Printf("  %s Running templ generate...\n", yellow("→"))
		templCmd := exec.Command("templ", "generate")
		templCmd.Stdout = os.Stdout
		templCmd.Stderr = os.Stderr
		if err := templCmd.Run(); err != nil {
			fmt.Printf("  %s templ generate failed: %v\n", red("Error:"), err)
			os.Exit(1)
		}
		fmt.Printf("  %s Templates generated\n", green("✓"))
	}

	// Build the binary
	fmt.Printf("  %s Building binary...\n", yellow("→"))

	buildArgs := []string{
		"build",
		"-ldflags", "-s -w", // Strip debug info for smaller binary
		"-o", buildOutput,
		".",
	}

	buildEnv := os.Environ()
	if buildOS != "" {
		buildEnv = append(buildEnv, fmt.Sprintf("GOOS=%s", buildOS))
	}
	if buildArch != "" {
		buildEnv = append(buildEnv, fmt.Sprintf("GOARCH=%s", buildArch))
	}

	goBuild := exec.Command("go", buildArgs...)
	goBuild.Env = buildEnv
	goBuild.Stdout = os.Stdout
	goBuild.Stderr = os.Stderr

	if err := goBuild.Run(); err != nil {
		fmt.Printf("  %s Build failed: %v\n", red("Error:"), err)
		os.Exit(1)
	}

	// Get binary size
	info, err := os.Stat(buildOutput)
	if err != nil {
		fmt.Printf("  %s Failed to stat binary: %v\n", yellow("Warning:"), err)
	}

	size := "unknown"
	if info != nil {
		sizeMB := float64(info.Size()) / 1024 / 1024
		size = fmt.Sprintf("%.2f MB", sizeMB)
	}

	fmt.Printf("  %s Build successful\n\n", green("✓"))
	fmt.Printf("  Output: %s\n", cyan(buildOutput))
	fmt.Printf("  Size:   %s\n", size)

	if buildOS != "" || buildArch != "" {
		fmt.Printf("  Target: %s/%s\n", targetOS, buildArch)
	}

	fmt.Printf("\n  Run with: %s\n\n", cyan("./"+buildOutput))
}
