package fuego

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the configuration for a Fuego application.
type Config struct {
	// Server configuration
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`

	// Directory configuration
	AppDir    string `mapstructure:"app_dir"`
	StaticDir string `mapstructure:"static_dir"`
	StaticURL string `mapstructure:"static_path"`

	// Development configuration
	Dev DevConfig `mapstructure:"dev"`

	// Middleware configuration
	Middleware MiddlewareConfig `mapstructure:"middleware"`
}

// DevConfig holds development-specific configuration.
type DevConfig struct {
	HotReload       bool     `mapstructure:"hot_reload"`
	WatchExtensions []string `mapstructure:"watch_extensions"`
	ExcludeDirs     []string `mapstructure:"exclude_dirs"`
}

// MiddlewareConfig holds middleware-specific configuration.
type MiddlewareConfig struct {
	Logger  bool `mapstructure:"logger"`
	Recover bool `mapstructure:"recover"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Port:      "3000",
		Host:      "0.0.0.0",
		AppDir:    "app",
		StaticDir: "static",
		StaticURL: "/static",
		Dev: DevConfig{
			HotReload:       true,
			WatchExtensions: []string{".go", ".templ"},
			ExcludeDirs:     []string{"node_modules", ".git", "_*"},
		},
		Middleware: MiddlewareConfig{
			Logger:  true,
			Recover: true,
		},
	}
}

// Address returns the full address string for the server.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// ListenAddress returns the address string for listening (typically :port).
func (c *Config) ListenAddress() string {
	return fmt.Sprintf(":%s", c.Port)
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if c.AppDir == "" {
		return fmt.Errorf("app_dir cannot be empty")
	}
	return nil
}

// LoadConfig loads configuration from fuego.yaml if it exists.
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	v := viper.New()
	v.SetConfigName("fuego")
	v.SetConfigType("yaml")

	// Add config path
	if path != "" {
		v.AddConfigPath(path)
	}
	v.AddConfigPath(".")

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is not an error, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// EnsureAppDir checks if the app directory exists.
func (c *Config) EnsureAppDir() error {
	absPath, err := filepath.Abs(c.AppDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoAppDir
		}
		return fmt.Errorf("failed to stat app dir: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", c.AppDir)
	}

	return nil
}
