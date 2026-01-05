# Contributing to Nexo

Thank you for your interest in contributing to Nexo! This document provides guidelines and information for contributors.

## Code of Conduct

Please be respectful and constructive in all interactions. We welcome contributors of all experience levels.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nexo.git
   cd nexo
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/abdul-hamid-achik/nexo.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Run tests to verify your setup:
   ```bash
   go test ./...
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes with clear, focused commits

3. Ensure all tests pass:
   ```bash
   go test ./...
   ```

4. Format your code:
   ```bash
   go fmt ./...
   ```

5. Run the linter (if installed):
   ```bash
   golangci-lint run
   ```

### Commit Messages

We follow conventional commit messages:

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions or changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

Examples:
```
feat: add support for WebSocket routes
fix: correct middleware execution order
docs: update README with new examples
```

### Pull Requests

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Open a Pull Request against the `main` branch

3. Fill out the PR template with:
   - Description of changes
   - Related issues (if any)
   - Testing done

4. Wait for CI to pass and address any review feedback

## Code Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Document all exported types and functions
- Keep functions focused and reasonably sized

### Testing

- Write table-driven tests when appropriate
- Aim for good coverage of new code
- Include both positive and negative test cases
- Test edge cases

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"normal case", "input", "expected"},
        {"edge case", "", ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Feature(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### Package Organization

- `cmd/nexo/` - CLI commands
- `pkg/nexo/` - Public API (the framework)
- `internal/` - Internal packages not for public use
- `examples/` - Example projects

### Documentation

- Update README.md for user-facing changes
- Add GoDoc comments for new public APIs
- Update relevant docs in `docs/` directory
- Include examples in documentation

## Project Structure

```
nexo/
├── cmd/nexo/           # CLI entry point
│   ├── main.go
│   └── commands/        # CLI commands
├── pkg/nexo/           # Main framework package
│   ├── app.go           # App struct
│   ├── context.go       # Request context
│   ├── router.go        # Route tree
│   ├── scanner.go       # File system scanner
│   ├── middleware.go    # Built-in middleware
│   └── renderer.go      # Templ rendering
├── internal/
│   ├── templates/       # Project scaffolding templates
│   └── version/         # Version information
├── examples/            # Example projects
└── docs/                # Documentation
```

## Reporting Issues

### Bug Reports

Please include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant code snippets

### Feature Requests

Please include:
- Use case description
- Proposed solution (if any)
- Alternatives considered

## Questions?

- Open an issue for questions
- Check existing issues for answers
- Review the documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Nexo!
