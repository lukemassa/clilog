# clilog

A simple logger for Go CLI applications. `clilog` is a partial drop-in replacement for the standard `log` package, with support for log levels, formatting, colorization, and timestamps — all with zero dependencies.

## Goals

- **Human-readable CLI logs**  
  Designed for command-line output, not machines or log aggregators.

- **Formatted log lines**  
  Uses text templates to control layout (`timestamp level message`), with support for color and millisecond timestamps.

- **Minimal API, partial drop-in for `log`**  
  Functions like `Info`, `Infof`, `Fatalf` behave like standard `log.Print*` but add levels and formatting.

- **Single global logger**  
  No logger instances or hierarchies — one logger, one output, meant for short-lived CLI tools.

- **No dependencies**  
  Just the Go standard library (`fmt`, `text/template`, etc.).

---

## Non-Goals

- **Structured logging**  
  No JSON, key-value pairs, or slog-style handlers. This is for humans, not log processors.

- **Logging to multiple outputs or files**  
  Output is a single `io.Writer`, typically `os.Stderr`. No rotation, no sinks.

- **Runtime log level reconfiguration**  
  You set the log level once — there's no dynamic or context-based level switching.

- **Log contexts or tracing metadata**  
  No `WithField`, no `WithContext`, no span propagation. Keep it simple.

- **Performance**
  Using go's templates makes it flexible and easy to configure, but not performant. It's good enough for a CLI however.

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
	log.SetDisableColor(false)

	log.Infof("starting up")
	log.Debug("initializing subsystems")

	// Optional custom format
	log.SetFormat(`{{ .Time }} {{ .LevelName }} {{ .Message }}`)

	log.Warn("deprecated feature in use")
	log.Errorf("could not read config file")
}
```

## Example Output

```
2025/06/22 16:08:19.103 D starting up
2025/06/22 16:08:19.104 D initializing subsystems
2025/06/22 16:08:19.105 W deprecated feature in use
2025/06/22 16:08:19.106 E could not read config file
```

If color is enabled, the timestamp and log level will be colorized based on severity.

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

Use Go’s `text/template` syntax to control log layout.

Available fields:
- `{{ .Time }}` — formatted timestamp (millisecond precision)
- `{{ .LevelCode }}` — short level code like `D`, `W`, `F`
- `{{ .LevelName }}` — padded full name: `DEBUG `, `ERROR `, etc.
- `{{ .Message }}` — your actual log content

Example format:
```go
log.SetFormat(`{{ .Time }} [{{ .LevelCode }}] {{ .Message }}`)
```

Invalid template variables will return an error at `SetFormat` time.

## Other Configuration

```go
log.SetLogLevel(log.LevelInfo)                      // Minimum level to print
log.SetDisableColor(true)                           // Disable color entirely
log.SetTimestampFormat("2006-01-02 15:04:05.000")   // Custom time layout
log.SetOutput(os.Stdout)                            // Redirect output
```

## License

MIT

