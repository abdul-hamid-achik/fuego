package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with hot reload",
	Long: `Start the development server with automatic hot reloading.

The server will automatically rebuild and restart when Go or templ files change.

Example:
  fuego dev
  fuego dev --port 8080`,
	Run: runDev,
}

var (
	devPort string
	devHost string
)

func init() {
	devCmd.Flags().StringVarP(&devPort, "port", "p", "3000", "Port to run the server on")
	devCmd.Flags().StringVarP(&devHost, "host", "H", "0.0.0.0", "Host to bind to")
}

func runDev(cmd *cobra.Command, args []string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("\n  %s Development Server\n\n", cyan("Fuego"))

	// Check for main.go or app directory
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		fmt.Printf("  %s No main.go found in current directory\n", red("Error:"))
		fmt.Printf("  Run this command from your project root\n\n")
		os.Exit(1)
	}

	// Check for templ files and run templ generate if needed
	hasTemplFiles := false
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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
			fmt.Printf("  %s templ generate failed (is templ installed?): %v\n", yellow("Warning:"), err)
			fmt.Printf("  Install with: go install github.com/a-h/templ/cmd/templ@latest\n\n")
		}
	}

	// Start the server
	var serverProcess *exec.Cmd
	serverProcess = startDevServer(devPort)

	// Set up file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("  %s Failed to create file watcher: %v\n", red("Error:"), err)
		os.Exit(1)
	}
	defer watcher.Close()

	// Watch directories recursively
	watchDirs := []string{"."}
	if _, err := os.Stat("app"); err == nil {
		watchDirs = append(watchDirs, "app")
	}

	for _, dir := range watchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			// Skip hidden directories and common non-source directories
			if info.IsDir() {
				name := info.Name()
				if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "tmp" {
					return filepath.SkipDir
				}
				watcher.Add(path)
			}
			return nil
		})
	}

	fmt.Printf("  %s Watching for changes...\n", green("✓"))
	fmt.Printf("\n  ➜ Local:   %s\n", cyan(fmt.Sprintf("http://localhost:%s", devPort)))
	fmt.Printf("  ➜ Network: %s\n\n", cyan(fmt.Sprintf("http://%s:%s", devHost, devPort)))

	// Debounce channel
	var debounceTimer *time.Timer
	debounceDuration := 100 * time.Millisecond

	// Signal handling
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only react to write, create, and remove events
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
				continue
			}

			// Check file extension
			ext := filepath.Ext(event.Name)
			if ext != ".go" && ext != ".templ" {
				continue
			}

			// Skip generated templ files
			if strings.HasSuffix(event.Name, "_templ.go") {
				continue
			}

			// Debounce
			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			debounceTimer = time.AfterFunc(debounceDuration, func() {
				timestamp := time.Now().Format("15:04:05")

				// Run templ generate if it's a templ file
				if ext == ".templ" {
					fmt.Printf("  [%s] %s Regenerating templates...\n", timestamp, yellow("→"))
					templCmd := exec.Command("templ", "generate")
					if err := templCmd.Run(); err != nil {
						fmt.Printf("  [%s] %s templ generate failed: %v\n", timestamp, red("✗"), err)
						return
					}
				}

				fmt.Printf("  [%s] %s Rebuilding...\n", timestamp, yellow("→"))

				// Stop old server
				if serverProcess != nil && serverProcess.Process != nil {
					serverProcess.Process.Signal(syscall.SIGTERM)
					serverProcess.Wait()
				}

				// Start new server
				serverProcess = startDevServer(devPort)

				fmt.Printf("  [%s] %s Ready\n", timestamp, green("✓"))
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("  %s Watcher error: %v\n", yellow("Warning:"), err)

		case <-signals:
			fmt.Println("\n  Shutting down...")
			if serverProcess != nil && serverProcess.Process != nil {
				serverProcess.Process.Signal(syscall.SIGTERM)
				serverProcess.Wait()
			}
			os.Exit(0)
		}
	}
}

func startDevServer(port string) *exec.Cmd {
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", port))

	if err := cmd.Start(); err != nil {
		fmt.Printf("  %s Failed to start server: %v\n", color.RedString("Error:"), err)
		return nil
	}

	return cmd
}
