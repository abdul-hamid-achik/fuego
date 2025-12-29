// Package version provides version information for the Fuego CLI.
package version

// Version is set via ldflags during build.
var Version = "dev"

// GetVersion returns the current version string.
func GetVersion() string {
	return Version
}
