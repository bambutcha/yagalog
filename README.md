<div align="center">

# YagaLog

Simple, colorful Go logger. Dual output (console + file). Thread‑safe. Zero setup.

[![Go Reference](https://pkg.go.dev/badge/github.com/Chelaran/yagalog.svg)](https://pkg.go.dev/github.com/Chelaran/yagalog)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.20+%20(1.25.3%20tested)-00ADD8?logo=go)]()

</div>

---

## Install last version

```bash
go get github.com/Chelaran/yagalog@v0.3.0
```

## TL;DR

```go
import (
    logger "github.com/Chelaran/yagalog"
)

// console only
l, _ := logger.NewLogger()
defer l.Close()

// or with file
l2, err := logger.NewLogger(logger.WithFilePath("app.log"))
if err == nil {
    defer l2.Close()
}

l.Info("Ready!")
```

## Quick start

```go
import logger "github.com/Chelaran/yagalog"

func main() {
    // console only
    l, err := logger.NewLogger()
    if err != nil { panic(err) }
    defer l.Close()

    // or with file
    lf, err := logger.NewLogger(logger.WithFilePath("app.log"))
    if err == nil {
        defer lf.Close()
    }

    l.Info("Hello from YagaLog %d", 1)
    l.Warning("Be careful")
    l.Error("Oops: %s", "something went wrong")
}
```

## Features
- Colorized console levels: DEBUG, INFO, WARNING, ERROR, FATAL
- File logging with date/time
- Dual output: console + file (simultaneously)
- Thread‑safe file writes
- Lightweight dependency footprint

## Levels

| Level   | Purpose                  |
|-------- |--------------------------|
| DEBUG   | подробные события в DEV  |
| INFO    | обычный поток событий    |
| WARNING | потенциальные проблемы   |
| ERROR   | ошибки, но сервис жив    |
| FATAL   | критично, завершение     |

## API
```go
// Create new logger with functional options
func NewLogger(opts ...Option) (*Logger, error)

// Logging methods (printf‑style formatting supported)
func (l *Logger) Debug(msg string, v ...interface{})
func (l *Logger) Info(msg string, v ...interface{})
func (l *Logger) Warning(msg string, v ...interface{})
func (l *Logger) Error(msg string, v ...interface{})
func (l *Logger) Fatal(msg string, v ...interface{}) // exits with code 1

// Gracefully close file handle
func (l *Logger) Close() error
```

## When to use
- You want a tiny, readable logger with colorized console and file output
- You don’t need JSON/structured logs, hooks, or advanced sinks (yet)

## Compatibility
- Go 1.20+ (tested with 1.25.3)
- OS: Linux/macOS/Windows

## Options (functional options API)

```go
// Examples using options at construction time or via helpers

// Construction-time options
lf, _ := logger.NewLogger(
        logger.WithFilePath("app.log"),
        logger.WithLevel(logger.INFO),
        logger.WithJSON(),
)
defer lf.Close()

// Runtime helpers remain available
l.SetLevel(logger.INFO)
l.WithColors(true) // toggles color output for console
// or use WithWriter to direct console logs to a custom io.Writer
buf := &bytes.Buffer{}
_ = logger.NewLogger(logger.WithWriter(buf))

// You can still DisableFile()/EnableFile(path) at runtime
// to switch file sinks dynamically.
```

## Changelog

- v0.3.0 — 2025-11-20
    - Introduced functional-options API for `NewLogger(opts ...Option)`.
    - New options: `WithFilePath`, `WithWriter`, `WithLevel`, `WithJSON`, `WithTimeFormat`, `WithCaller`, `WithColors`.
    - Backwards compatibility: constructor signature changed — replace `NewLogger(path)` with `NewLogger(WithFilePath(path))` or call `NewLogger()` for console-only.

- v0.2.0 — previous

## Examples
Run any example (more in `examples/`):
```bash
go run ./examples/basic    # INFO level, RFC3339Nano time
go run ./examples/caller   # adds file:line
go run ./examples/json     # writes JSON lines to app.log
```

## Best practices

DEV preset
```go
l.SetLevel(logger.DEBUG)
l.WithCaller(true)
// colors auto (TTY), human‑readable console + file
```

PROD preset
```go
l.SetLevel(logger.INFO)
l.WithCaller(false)
l.WithJSON()           // JSON lines to file for aggregators
// optionally: l.DisableFile() if logs collected from stdout only
```

## Roadmap (short)
- Optional JSON/structured output
- Custom formatters and minimal hook interface
- Rolling file options

## Contributing
Issues and PRs are welcome. Keep code clear, small, and well‑tested.

## Security
Found a vulnerability? Please open a private report via GitHub Security Advisories.

## License
MIT © 2025 Chelaran. Created by bambutcha.
