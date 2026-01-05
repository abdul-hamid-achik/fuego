package nexo

// Option is a functional option for configuring the App.
type Option func(*App)

// WithPort sets the server port.
func WithPort(port string) Option {
	return func(a *App) {
		a.config.Port = port
	}
}

// WithHost sets the server host.
func WithHost(host string) Option {
	return func(a *App) {
		a.config.Host = host
	}
}

// WithAddress sets both host and port from an address string (e.g., ":3000" or "0.0.0.0:8080").
func WithAddress(addr string) Option {
	return func(a *App) {
		// Parse address - if it starts with ":", it's just a port
		if len(addr) > 0 && addr[0] == ':' {
			a.config.Port = addr[1:]
		} else {
			// Try to split host:port
			for i := len(addr) - 1; i >= 0; i-- {
				if addr[i] == ':' {
					a.config.Host = addr[:i]
					a.config.Port = addr[i+1:]
					return
				}
			}
			// No colon found, assume it's just a port
			a.config.Port = addr
		}
	}
}

// WithAppDir sets the app directory.
func WithAppDir(dir string) Option {
	return func(a *App) {
		a.config.AppDir = dir
	}
}

// WithStaticDir sets the static files directory.
func WithStaticDir(dir string) Option {
	return func(a *App) {
		a.config.StaticDir = dir
	}
}

// WithStaticURL sets the URL path for serving static files.
func WithStaticURL(url string) Option {
	return func(a *App) {
		a.config.StaticURL = url
	}
}

// WithConfig sets the entire configuration.
func WithConfig(config *Config) Option {
	return func(a *App) {
		if config != nil {
			a.config = config
		}
	}
}

// WithLogger enables or disables the logger middleware.
func WithLogger(enabled bool) Option {
	return func(a *App) {
		a.config.Middleware.Logger = enabled
	}
}

// WithRecover enables or disables the recover middleware.
func WithRecover(enabled bool) Option {
	return func(a *App) {
		a.config.Middleware.Recover = enabled
	}
}

// WithHotReload enables or disables hot reload in development.
func WithHotReload(enabled bool) Option {
	return func(a *App) {
		a.config.Dev.HotReload = enabled
	}
}
