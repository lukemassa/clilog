# clilog

A simple logger for Go CLI applications. `clilog` is a partial drop-in replacement for the standard `log` package, with support for log levels, formatting, colorization, and timestamps with zero dependencies.

## Goals

- **Human-readable CLI logs**  
  Designed for command-line output, not machines or log aggregators.

- **Formatted log lines**  
  Uses text templates to control layout (`{{ .Level }} {{ .Message }}`), with support for color.

- **Minimal API, partial drop-in for `log`**  
  Functions like `Info`, `Infof`, `Fatalf` behave like standard `log.Print*` but add levels and formatting via simple Set* functions.

- **Single global logger**  
  No logger instances or hierarchies — one logger, one output, meant for short-lived CLI tools.

- **No dependencies**  
  Just the Go standard library (`fmt`, `text/template`, etc.).

---

## Non-Goals

- **Structured logging**  
  No JSON, key-value pairs, or slog-style handlers. This is for humans, not log processors.

- **Logging to multiple outputs or files**  
  Output is a single `io.Writer`, currently unconfigurably `os.Stderr`. No rotation, no sinks.

- **Runtime log level reconfiguration**  
  You set the log level once — there's no dynamic or context-based level switching.

- **Log contexts or tracing metadata**  
  No `WithField`, no `WithContext`, no span propagation. Keep it simple.

- **Performance**  
  Template rendering is flexible but not fast. This logger prioritizes readability and configurability over speed — fine for a CLI, but not for high-throughput services.

- **Framework integration**  
  This logger is intentionally agnostic — it doesn’t integrate with Viper, Cobra, or slog adapters.


## Installation

```bash
go get github.com/lukemassa/clilog
```

## Example

```go
package main

import log "github.com/lukemassa/clilog"

func main() {
	log.SetLogLevel(log.LevelDebug)

	log.Debug("initializing subsystems")
	log.Infof("Hello %s", "World")

	log.Warn("deprecated feature in use")
	log.Error("could not read config file")
}
```

## Example Output

```
D 2025/06/22 22:50:37.130 initializing subsystems
I 2025/06/22 22:50:37.130 Hello World
W 2025/06/22 22:50:37.130 deprecated feature in use
E 2025/06/22 22:50:37.130 could not read config fil
```

By default, the timestamp will be colorized based on severity.

## Log Levels

```go
log.LevelDebug // D
log.LevelInfo  // I
log.LevelWarn  // W
log.LevelError // E
log.LevelFatal // F (calls os.Exit(1))
```

Set the minimum level to show with:

```go
log.SetLogLevel(log.LevelWarn) // Only warn and above will print
```

## Formatting

We use Go’s `text/template` syntax to customize the layout of each log line.

### Available fields

- `{{ .Time }}` — the raw `time.Time` value (you’ll usually format this with `timef`)
- `{{ .Level }}` — the log level enum (e.g. `LevelInfo`, `LevelError`)
- `{{ .Message }}` — the actual log message content

### Template functions

- `timef "<layout>"` — formats a `time.Time` using Go’s `Time.Format` layout string
- `color <level>` — adds ANSI color codes based on the level
- `abbrev` — converts a `Level` to a single-letter abbreviation (e.g. `D`, `I`, `W`, `E`, `F`)

### Example

```go
log.SetFormat(`{{ .Level | abbrev }} {{ .Time | timef "2006/01/02 15:04:05.000" | color .Level }} {{ .Message }}`)
```

This starts with the abbreviation of the level (i.e. `D` for Debug) followed by colorized millisecond precision timestamp, followed by the message.

(Note: this is currently the default format).


