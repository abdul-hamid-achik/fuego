package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abdul-hamid-achik/fuego/internal/version"
	"github.com/abdul-hamid-achik/fuego/pkg/tools"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Fuego to the latest version",
	Long: `Check for and install the latest version of Fuego.

By default, prereleases are skipped. Use --prerelease to include them.

Examples:
  fuego upgrade                    Upgrade to latest stable version
  fuego upgrade --check            Check for updates without installing
  fuego upgrade --version v0.5.0   Install a specific version
  fuego upgrade --prerelease       Include prerelease versions
  fuego upgrade --rollback         Restore previous version from backup`,
	Run: runUpgrade,
}

var (
	upgradeCheck      bool
	upgradeVersion    string
	upgradePrerelease bool
	upgradeForce      bool
	upgradeRollback   bool
)

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeCheck, "check", false,
		"Check for updates without installing")
	upgradeCmd.Flags().StringVar(&upgradeVersion, "version", "",
		"Install a specific version (e.g., v0.5.0)")
	upgradeCmd.Flags().BoolVar(&upgradePrerelease, "prerelease", false,
		"Include prerelease versions")
	upgradeCmd.Flags().BoolVar(&upgradeForce, "force", false,
		"Force upgrade even if same version")
	upgradeCmd.Flags().BoolVar(&upgradeRollback, "rollback", false,
		"Restore the previous version from backup")

	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if !jsonOutput {
		fmt.Printf("\n  %s Upgrade\n\n", cyan("Fuego"))
	}

	currentVersion := version.GetVersion()

	// Handle rollback
	if upgradeRollback {
		runRollback(currentVersion)
		return
	}

	updater := tools.NewUpdater()
	updater.IncludePrerelease = upgradePrerelease

	// Get release info
	var release *tools.ReleaseInfo
	var err error
	var hasUpdate bool

	if upgradeVersion != "" {
		// Specific version requested
		if !jsonOutput {
			fmt.Printf("  %s Fetching version %s...\n", yellow("->"), upgradeVersion)
		}
		release, err = updater.GetSpecificRelease(upgradeVersion)
		if err != nil {
			handleUpgradeError(err)
			return
		}
		// For specific version, always consider it an "update" unless force is needed
		hasUpdate = upgradeForce || tools.CompareVersions(currentVersion, release.TagName) != 0
	} else {
		// Latest version
		if !jsonOutput {
			fmt.Printf("  %s Checking for updates...\n", yellow("->"))
		}
		release, hasUpdate, err = updater.CheckForUpdate()
		if err != nil {
			handleUpgradeError(err)
			return
		}
	}

	// Display version info
	if !jsonOutput {
		fmt.Printf("  Current version: %s\n", currentVersion)
		fmt.Printf("  Latest version:  %s", release.TagName)
		if !release.PublishedAt.IsZero() {
			fmt.Printf(" (released %s)", humanizeTime(release.PublishedAt))
		}
		fmt.Println()
		fmt.Println()
	}

	// Check if already up to date
	if !hasUpdate && !upgradeForce {
		if jsonOutput {
			printSuccess(UpgradeOutput{
				CurrentVersion: currentVersion,
				LatestVersion:  release.TagName,
				UpToDate:       true,
			})
		} else {
			fmt.Printf("  %s You're already running the latest version (%s)\n\n",
				green("OK"), currentVersion)
		}
		return
	}

	// Check-only mode
	if upgradeCheck {
		if jsonOutput {
			printSuccess(UpgradeOutput{
				CurrentVersion:  currentVersion,
				LatestVersion:   release.TagName,
				UpdateAvailable: true,
				ReleaseNotes:    release.Body,
				PublishedAt:     release.PublishedAt,
			})
		} else {
			fmt.Printf("  %s Update available!\n", green("OK"))
			fmt.Printf("  Run '%s' to update.\n\n", yellow("fuego upgrade"))

			// Show abbreviated release notes
			if release.Body != "" {
				fmt.Println("  Release notes:")
				printReleaseNotes(release.Body, 5)
				fmt.Println()
			}
		}
		return
	}

	// Find correct asset for this platform
	asset, err := updater.GetAssetForPlatform(release)
	if err != nil {
		handleUpgradeError(err)
		return
	}

	// Download
	if !jsonOutput {
		fmt.Printf("  %s Downloading %s...\n", yellow("->"), asset.Name)
	}

	archivePath, err := updater.Download(asset)
	if err != nil {
		handleUpgradeError(fmt.Errorf("download failed: %w", err))
		return
	}
	defer func() { _ = os.Remove(archivePath) }()

	// Verify checksum
	if !jsonOutput {
		fmt.Printf("  %s Verifying checksum...\n", yellow("->"))
	}

	if err := updater.VerifyChecksum(archivePath, release); err != nil {
		handleUpgradeError(fmt.Errorf("checksum verification failed: %w", err))
		return
	}

	// Extract binary
	if !jsonOutput {
		fmt.Printf("  %s Extracting binary...\n", yellow("->"))
	}

	binaryPath, err := updater.ExtractBinary(archivePath)
	if err != nil {
		handleUpgradeError(fmt.Errorf("extraction failed: %w", err))
		return
	}
	defer func() { _ = os.Remove(binaryPath) }()

	// Install
	if !jsonOutput {
		fmt.Printf("  %s Installing...\n", yellow("->"))
	}

	if err := updater.Install(binaryPath); err != nil {
		handleUpgradeError(fmt.Errorf("installation failed: %w", err))
		return
	}

	// Success!
	if jsonOutput {
		printSuccess(UpgradeOutput{
			CurrentVersion:  currentVersion,
			LatestVersion:   release.TagName,
			UpgradeComplete: true,
			ReleaseNotes:    release.Body,
			BackupPath:      updater.BackupPath(),
		})
	} else {
		fmt.Printf("  %s Upgraded successfully to %s!\n\n",
			green("OK"), release.TagName)

		fmt.Printf("  Backup saved to: %s\n", updater.BackupPath())
		fmt.Printf("  To rollback: %s\n\n", yellow("fuego upgrade --rollback"))

		// Show release notes (abbreviated)
		if release.Body != "" {
			fmt.Println("  Release notes:")
			printReleaseNotes(release.Body, 8)
			fmt.Println()
		}
	}
}

func runRollback(currentVersion string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	updater := tools.NewUpdater()

	if !updater.HasBackup() {
		if jsonOutput {
			printJSONError(fmt.Errorf("no backup found"))
		} else {
			fmt.Printf("  %s No backup found\n", yellow("Warning:"))
			fmt.Printf("  Backup location: %s\n\n", updater.BackupPath())
		}
		os.Exit(1)
	}

	if !jsonOutput {
		fmt.Printf("  Current version: %s\n", currentVersion)
		fmt.Printf("  %s Restoring from backup...\n", yellow("->"))
	}

	if err := updater.Rollback(); err != nil {
		handleUpgradeError(fmt.Errorf("rollback failed: %w", err))
		return
	}

	if jsonOutput {
		printSuccess(UpgradeOutput{
			CurrentVersion:  currentVersion,
			UpgradeComplete: true,
			BackupPath:      updater.BackupPath(),
		})
	} else {
		fmt.Printf("  %s Rollback successful!\n\n", green("OK"))
		fmt.Printf("  Run '%s' to verify the restored version.\n\n",
			cyan("fuego --version"))
	}
}

func handleUpgradeError(err error) {
	if jsonOutput {
		printJSONError(err)
	} else {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("  %s %v\n\n", red("Error:"), err)
	}
	os.Exit(1)
}

func humanizeTime(t time.Time) string {
	diff := time.Since(t)
	if diff < time.Hour*24 {
		return "today"
	} else if diff < time.Hour*24*2 {
		return "yesterday"
	} else if diff < time.Hour*24*7 {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	} else if diff < time.Hour*24*30 {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < time.Hour*24*365 {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
	return t.Format("Jan 2, 2006")
}

func printReleaseNotes(body string, maxLines int) {
	lines := strings.Split(body, "\n")
	yellow := color.New(color.FgYellow).SprintFunc()

	for i, line := range lines {
		if i >= maxLines {
			fmt.Printf("    %s\n", yellow("..."))
			break
		}
		// Indent each line
		if strings.TrimSpace(line) != "" {
			fmt.Printf("    %s\n", line)
		}
	}
}

// CheckForUpdateInBackground checks for updates without blocking
// This is called from the dev command
func CheckForUpdateInBackground() {
	updater := tools.NewUpdater()

	// Rate limit: only check once per 24 hours
	if !updater.ShouldCheckForUpdate() {
		return
	}

	release, hasUpdate, err := updater.CheckForUpdate()
	if err != nil || !hasUpdate {
		// Save check time even on error to avoid hammering the API
		_ = updater.SaveLastCheckTime()
		return
	}

	// Save check time
	_ = updater.SaveLastCheckTime()

	// Print notification (after a small delay to not interfere with startup)
	time.Sleep(500 * time.Millisecond)

	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\n  %s A new version of Fuego is available: %s -> %s\n",
		yellow("Update:"),
		version.GetVersion(),
		cyan(release.TagName))
	fmt.Printf("  Run '%s' to upgrade.\n\n", yellow("fuego upgrade"))
}
