# loglinter

A Go linter that checks log messages for style and security compliance. Compatible with `golangci-lint`.

## Rules

- **lowercase** — log message must start with a lowercase letter
- **english** — log message must be in English only
- **special_chars** — log message must not contain special characters or emoji
- **sensitive_data** — log message must not contain sensitive data keywords

## Supported loggers

- `log/slog`
- `go.uber.org/zap`
- `log`

## Requirements

- Go 1.22+
- golangci-lint v2+

## Installation
```bash
git clone https://github.com/dementievme/golangci-lint-plugin.git
cd golangci-lint-plugin
go mod tidy
```

## Build

Build a custom golangci-lint binary with the plugin:
```bash
golangci-lint custom
```

This will produce a `custom-gcl` binary in the project root.

## Configuration

Create `config/config.yml`:
```yaml
extra_sensitive_keywords:
  - ssn
  - credit_card
  - cvv

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
./custom-gcl run --config ./config/config.yml ./...
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
./custom-gcl run ./...

# standalone
go build -o loglinter ./cmd/loglinter/
./loglinter ./...
```

## Tests
```bash
go test -v ./...
```
