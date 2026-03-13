# loglinter

A Go linter that checks log messages for style and security compliance. Compatible with `golangci-lint`.

## Rules

- **lowercase** – log message must start with a lowercase letter
- **english** – log message must be in English only
- **special_chars** – log message must not contain special characters or emoji
- **sensitive_data** – log message must not contain sensitive data keywords

## Supported loggers

- `log/slog`
- `go.uber.org/zap`
- `log`

## Requirements

- Go 1.22+
- golangci-lint v2+
- lefthook

## Installation
```bash
git clone https://github.com/dementievme/golangci-lint-plugin.git
cd golangci-lint-plugin
go mod tidy
```

## Build

Build a custom golangci-lint binary with the plugin:
```bash
cd plugin
golangci-lint custom
```

This will produce a `custom-gcl` binary in the `plugin` directory.

## Configuration

Create `config/config.yml`:
```yaml
extra_sensitive_keywords:
  - token
  - token_secret
  - password

disable_rules: []

loggers:
  log/slog:
    Debug: true
    Info: true
    Warn: true
    Error: true
    DebugContext: true
    InfoContext: true
    WarnContext: true
    ErrorContext: true
  go.uber.org/zap:
    Debug: true
    Info: true
    Warn: true
    Error: true
    Debugf: true
    Infof: true
    Warnf: true
    Errorf: true
    Debugw: true
    Infow: true
    Warnw: true
    Errorw: true
  log:
    Print: true
    Printf: true
    Println: true
    Fatal: true
    Panic: true
```

Set the config path via `.env` in the project root:
```bash
CONFIG_PATH=./config/config.yml
```

Or pass it as a flag:
```bash
./plugin/custom-gcl run --config ./config/config.yml ./...
```

To disable specific rules:
```yaml
disable_rules:
  - lowercase
  - special_chars
```

To add custom sensitive keywords:
```yaml
extra_sensitive_keywords:
  - ssn
  - credit_card
```

## Usage
```bash
# via custom-gcl
./plugin/custom-gcl run ./...

# standalone
go build -o loglinter ./cmd/loglinter/
./loglinter ./...
```

## Git hooks

Install lefthook:
```bash
go install github.com/evilmartians/lefthook@latest
lefthook install
```

Runs before every commit automatically:
```
golangci-lint run ./...
go test -v ./...
```

## CI/CD

CI runs automatically on push to `development` and on pull requests to `main`.

Jobs:
- **Test** – runs tests with race detector
- **Lint** – builds custom golangci-lint and runs linter
- **Build** – builds standalone binary

## Tests
```bash
go test -v ./...
```
