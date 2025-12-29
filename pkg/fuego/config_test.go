package fuego

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Port != "3000" {
		t.Errorf("expected default port 3000, got %s", config.Port)
	}
	if config.Host != "0.0.0.0" {
		t.Errorf("expected default host 0.0.0.0, got %s", config.Host)
	}
	if config.AppDir != "app" {
		t.Errorf("expected default app_dir 'app', got %s", config.AppDir)
	}
	if config.StaticDir != "static" {
		t.Errorf("expected default static_dir 'static', got %s", config.StaticDir)
	}
	if config.StaticURL != "/static" {
		t.Errorf("expected default static_url '/static', got %s", config.StaticURL)
	}
	if !config.Dev.HotReload {
		t.Error("expected default hot_reload to be true")
	}
	if !config.Middleware.Logger {
		t.Error("expected default middleware.logger to be true")
	}
	if !config.Middleware.Recover {
		t.Error("expected default middleware.recover to be true")
	}
}

func TestConfig_Address(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		expected string
	}{
		{"default", "0.0.0.0", "3000", "0.0.0.0:3000"},
		{"localhost", "localhost", "8080", "localhost:8080"},
		{"custom port", "127.0.0.1", "9000", "127.0.0.1:9000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Host: tt.host, Port: tt.port}
			if got := config.Address(); got != tt.expected {
				t.Errorf("Address() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestConfig_ListenAddress(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		expected string
	}{
		{"default", "3000", ":3000"},
		{"custom", "8080", ":8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Port: tt.port}
			if got := config.ListenAddress(); got != tt.expected {
				t.Errorf("ListenAddress() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"valid", &Config{Port: "3000", AppDir: "app"}, false},
		{"empty port", &Config{Port: "", AppDir: "app"}, true},
		{"empty app_dir", &Config{Port: "3000", AppDir: ""}, true},
		{"both empty", &Config{Port: "", AppDir: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_EnsureAppDir(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		config := &Config{AppDir: tmpDir}
		if err := config.EnsureAppDir(); err != nil {
			t.Errorf("EnsureAppDir() unexpected error for existing dir: %v", err)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		config := &Config{AppDir: "/nonexistent/path/that/does/not/exist"}
		err := config.EnsureAppDir()
		if err == nil {
			t.Error("EnsureAppDir() expected error for non-existent dir")
		}
		if err != ErrNoAppDir {
			t.Errorf("EnsureAppDir() expected ErrNoAppDir, got %v", err)
		}
	})

	t.Run("file instead of directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "not-a-dir")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		config := &Config{AppDir: filePath}
		err := config.EnsureAppDir()
		if err == nil {
			t.Error("EnsureAppDir() expected error for file path")
		}
	})
}

func TestLoadConfig_DefaultsWhenNoFile(t *testing.T) {
	// Load from a directory without fuego.yaml - should return defaults
	tmpDir := t.TempDir()
	config, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	// Should have default values
	if config.Port != "3000" {
		t.Errorf("expected default port, got %s", config.Port)
	}
}

func TestLoadConfig_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configContent := `
port: "8080"
host: "127.0.0.1"
app_dir: "myapp"
static_dir: "public"
static_path: "/assets"
dev:
  hot_reload: false
middleware:
  logger: false
  recover: false
`
	configPath := filepath.Join(tmpDir, "fuego.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	if config.Port != "8080" {
		t.Errorf("expected port 8080, got %s", config.Port)
	}
	if config.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", config.Host)
	}
	if config.AppDir != "myapp" {
		t.Errorf("expected app_dir 'myapp', got %s", config.AppDir)
	}
	if config.StaticDir != "public" {
		t.Errorf("expected static_dir 'public', got %s", config.StaticDir)
	}
	if config.StaticURL != "/assets" {
		t.Errorf("expected static_url '/assets', got %s", config.StaticURL)
	}
	if config.Dev.HotReload {
		t.Error("expected hot_reload to be false")
	}
	if config.Middleware.Logger {
		t.Error("expected middleware.logger to be false")
	}
	if config.Middleware.Recover {
		t.Error("expected middleware.recover to be false")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidYAML := `
port: "8080"
  invalid_indent: true
`
	configPath := filepath.Join(tmpDir, "fuego.yaml")
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadConfig(tmpDir)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid YAML")
	}
}
