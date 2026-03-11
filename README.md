# loglinter

A Go static analysis linter that checks log messages for style and security
compliance. Compatible with **golangci-lint** as a custom plugin.

## Supported loggers

| Package | Methods checked |
|---------|-----------------|
| `log/slog` | `Debug`, `Info`, `Warn`, `Error`, `*Context` variants |
| `go.uber.org/zap` | `Debug*`, `Info*`, `Warn*`, `Error*`, `Fatal*`, `Panic*` |
| `log` (stdlib) | `Print*`, `Fatal*`, `Panic*` |

## Rules

| ID | Rule | Example violation |
|----|------|-------------------|
| `lowercase` | Message must start with a lowercase letter | `"Starting server"` |
| `english` | Message must be in English (ASCII) only | `"запуск сервера"` |
| `special_chars` | No special characters or emoji | `"done! 🚀"` |
| `sensitive_data` | No sensitive keywords (password, token, etc.) | `"token: " + tok` |

---

## Installation

### Prerequisites

- Go 1.22+
- `golangci-lint` v1.57+ (for plugin support)

### Build the standalone binary

```bash
git clone https://github.com/yourusername/loglinter.git
cd loglinter
go build -o loglinter ./cmd/loglinter/
```

Run it directly on a package:

```bash
./loglinter ./...
```

### Build the golangci-lint plugin

```bash
go build -buildmode=plugin -o loglinter.so ./plugin/
```

Then reference the `.so` file in `.golangci.yml`:

```yaml
linters-settings:
  custom:
    loglinter:
      path: ./loglinter.so
      description: "Log message linter"
      original-url: "https://github.com/yourusername/loglinter"

linters:
  enable:
    - loglinter
```

Run:

```bash
golangci-lint run ./...
```

---

## Configuration (bonus feature)

Create a `loglinter.yml` next to your `go.mod`:

```yaml
# loglinter.yml
extra_sensitive_keywords:
  - ssn
  - credit_card
  - cvv

disable_rules:
  - special_chars   # disable the special-characters rule
```

---

## Examples

### ❌ Violations

```go
slog.Info("Starting server on port 8080")  // uppercase start
slog.Error("запуск сервера")               // non-English
slog.Warn("server started! 🚀")            // emoji
slog.Info("token: " + tok)                 // sensitive data
```

### ✅ Correct

```go
slog.Info("starting server on port 8080")
slog.Error("failed to connect to database")
slog.Warn("server started")
slog.Info("token validated")
```

---

## Development

### Run tests

```bash
go test -v -race ./...
```

### Project layout

```
loglinter/
├── analyzer/
│   ├── analyzer.go       # core analysis logic
│   ├── analyzer_test.go  # integration tests (analysistest)
│   └── rules_test.go     # unit tests per rule
├── cmd/loglinter/
│   └── main.go           # standalone CLI entry point
├── plugin/
│   └── main.go           # golangci-lint plugin entry point
├── testdata/src/
│   ├── basic/            # basic rule test cases
│   └── sensitive/        # sensitive data test cases
├── .golangci.yml         # example golangci-lint config
├── .github/workflows/ci.yml
└── go.mod
```

### Auto-fix (SuggestedFixes)

The `lowercase` rule emits an `analysis.SuggestedFix` so tools that support it
(e.g. `gopls`) can apply the fix automatically:

```bash
# apply suggested fixes
go vet -fix ./...
```

---

## Contributing

PRs welcome. Please add test cases in `testdata/src/` for any new rule.

## License

MIT
