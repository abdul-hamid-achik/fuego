package nexo

import (
	"testing"
)

func TestWithPort(t *testing.T) {
	app := New()
	WithPort("8080")(app)
	if app.config.Port != "8080" {
		t.Errorf("expected port 8080, got %s", app.config.Port)
	}
}

func TestWithHost(t *testing.T) {
	app := New()
	WithHost("localhost")(app)
	if app.config.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", app.config.Host)
	}
}

func TestWithAddress(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		expectedHost string
		expectedPort string
	}{
		{"port only with colon", ":8080", "0.0.0.0", "8080"},
		{"host and port", "127.0.0.1:9000", "127.0.0.1", "9000"},
		{"localhost and port", "localhost:3000", "localhost", "3000"},
		{"port only without colon", "8080", "0.0.0.0", "8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			WithAddress(tt.address)(app)
			if app.config.Host != tt.expectedHost {
				t.Errorf("expected host %s, got %s", tt.expectedHost, app.config.Host)
			}
			if app.config.Port != tt.expectedPort {
				t.Errorf("expected port %s, got %s", tt.expectedPort, app.config.Port)
			}
		})
	}
}

func TestWithAppDir(t *testing.T) {
	app := New()
	WithAppDir("myapp")(app)
	if app.config.AppDir != "myapp" {
		t.Errorf("expected app_dir myapp, got %s", app.config.AppDir)
	}
}

func TestWithStaticDir(t *testing.T) {
	app := New()
	WithStaticDir("public")(app)
	if app.config.StaticDir != "public" {
		t.Errorf("expected static_dir public, got %s", app.config.StaticDir)
	}
}

func TestWithStaticURL(t *testing.T) {
	app := New()
	WithStaticURL("/assets")(app)
	if app.config.StaticURL != "/assets" {
		t.Errorf("expected static_url /assets, got %s", app.config.StaticURL)
	}
}

func TestWithConfig(t *testing.T) {
	t.Run("with valid config", func(t *testing.T) {
		customConfig := &Config{
			Port:   "9999",
			Host:   "custom.host",
			AppDir: "custom_app",
		}

		app := New()
		WithConfig(customConfig)(app)

		if app.config.Port != "9999" {
			t.Errorf("expected port 9999, got %s", app.config.Port)
		}
		if app.config.Host != "custom.host" {
			t.Errorf("expected host custom.host, got %s", app.config.Host)
		}
	})

	t.Run("with nil config", func(t *testing.T) {
		app := New()
		originalPort := app.config.Port

		WithConfig(nil)(app)

		// Should not change anything when nil
		if app.config.Port != originalPort {
			t.Error("config should not change when WithConfig(nil) is called")
		}
	})
}

func TestWithLogger(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			WithLogger(tt.enabled)(app)
			if app.config.Middleware.Logger != tt.expected {
				t.Errorf("expected logger %v, got %v", tt.expected, app.config.Middleware.Logger)
			}
		})
	}
}

func TestWithRecover(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			WithRecover(tt.enabled)(app)
			if app.config.Middleware.Recover != tt.expected {
				t.Errorf("expected recover %v, got %v", tt.expected, app.config.Middleware.Recover)
			}
		})
	}
}

func TestWithHotReload(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled", true, true},
		{"disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New()
			WithHotReload(tt.enabled)(app)
			if app.config.Dev.HotReload != tt.expected {
				t.Errorf("expected hot_reload %v, got %v", tt.expected, app.config.Dev.HotReload)
			}
		})
	}
}

func TestNewWithMultipleOptions(t *testing.T) {
	app := New(
		WithPort("8080"),
		WithHost("127.0.0.1"),
		WithAppDir("myapp"),
		WithStaticDir("public"),
		WithStaticURL("/assets"),
		WithLogger(false),
		WithRecover(false),
		WithHotReload(false),
	)

	if app.config.Port != "8080" {
		t.Errorf("expected port 8080, got %s", app.config.Port)
	}
	if app.config.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", app.config.Host)
	}
	if app.config.AppDir != "myapp" {
		t.Errorf("expected app_dir myapp, got %s", app.config.AppDir)
	}
	if app.config.StaticDir != "public" {
		t.Errorf("expected static_dir public, got %s", app.config.StaticDir)
	}
	if app.config.StaticURL != "/assets" {
		t.Errorf("expected static_url /assets, got %s", app.config.StaticURL)
	}
	if app.config.Middleware.Logger {
		t.Error("expected logger to be false")
	}
	if app.config.Middleware.Recover {
		t.Error("expected recover to be false")
	}
	if app.config.Dev.HotReload {
		t.Error("expected hot_reload to be false")
	}
}
