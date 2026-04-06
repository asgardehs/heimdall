# Heimdall

_Shared configuration library for the Asgard ecosystem._

Heimdall provides a typed, validated configuration store that Odin, Muninn, and
Huginn import directly. Configuration lives in a local SQLite database with
schema validation, change history, and a CLI for manual management.

## Usage as a Go Library

```go
import "github.com/asgardehs/heimdall"

h, err := heimdall.Open(logger)
defer h.Close()

// Read config
entry, err := h.Get("muninn", "vault_path")
fmt.Println(entry.Value)

// Write config (validated against schema)
err = h.Set("ai", "provider", "anthropic")

// List all keys in a namespace
entries, err := h.List("odin")

// Reset to default
err = h.Reset("odin", "theme")
```

## CLI

```bash
heimdall config get muninn vault_path
heimdall config set ai provider anthropic
heimdall config list odin
heimdall config reset odin theme
```

Secrets (like `ai.api_key`) are masked in CLI output.

## Documentation

See the **[Wiki](https://asgardehs.github.io/docs/heimdall/)** for full
documentation:

- [Quick Start](https://asgardehs.github.io/docs/heimdall/quickstart/) — setup
  and first use
- [Go API](https://asgardehs.github.io/docs/heimdall/api/) — library reference
- [CLI Reference](https://asgardehs.github.io/docs/heimdall/cli/) — command
  line usage
- [Configuration Keys](https://asgardehs.github.io/docs/heimdall/keys/) —
  namespaces, types, and defaults

## Building

```bash
# Library — just import it
go get github.com/asgardehs/heimdall

# CLI
go build -o heimdall ./cmd/heimdall
```

No CGO required. Uses pure Go SQLite (`modernc.org/sqlite`).

## License

GPLv3 — see [LICENSE](LICENSE).

## Name

> _In Norse mythology, Heimdall is the watchful guardian of the Bifrost bridge,
> keeper of the horn Gjallarhorn. He sees all, hears all, and ensures the
> boundaries of Asgard are maintained. Here, Heimdall guards the shared
> configuration that binds the ecosystem together._

_Part of the [Asgard EHS Family](https://asgardehs.github.io/)_
