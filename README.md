# Heimdall

_Shared configuration library for the Asgard EHS ecosystem._

![Status](https://img.shields.io/badge/status-active-2A6E3F?style=for-the-badge)

## Overview

Every tool in the Asgard EHS ecosystem needs configuration: vault paths,
database locations, API keys, UI preferences. Heimdall is the shared library
that keeps that configuration in one validated place, so Odin can read the
same vault path Muninn wrote, the CLI can inspect what any tool has stored,
and invalid values fail at write time rather than at runtime. It is a library
first — Odin, Muninn, and Huginn import it directly — with a CLI for manual
management.

## When Not to Use Heimdall

Heimdall is a local, per-machine config store. It is:

- **Not a secrets manager.** Secrets are stored in plaintext SQLite and masked
  only in CLI output. If you need encryption at rest or OS keychain
  integration, use a dedicated secrets tool.
- **Not a runtime config service.** Each process reads configuration on
  startup. Cross-process change notification is not provided.
- **Not a distributed store.** The SQLite database lives on local disk and is
  not synchronized across machines.

For those use cases, reach for something designed for them.

## Quick Example

### As a library

```go
import "github.com/asgardehs/heimdall"

h, err := heimdall.Open(logger)
if err != nil {
    return err
}
defer h.Close()

// Read config: returns a ConfigEntry with value, type, and source.
entry, err := h.Get("muninn", "vault_path")
fmt.Printf("%s = %s (type: %s, source: %s)\n",
    entry.Key, entry.Value, entry.Type, entry.Source)

// Write config: validated against the schema.
err = h.Set("ai", "provider", "anthropic")

// List all keys in a namespace.
entries, err := h.List("odin")

// Reset to the schema default.
err = h.Reset("odin", "theme")
```

### As a CLI

```bash
heimdall config get muninn vault_path
heimdall config set ai provider anthropic
heimdall config list odin
heimdall config reset odin theme
```

Secrets (like `ai.api_key`) are masked in CLI output.

## Installation

**As a library:**

```bash
go get github.com/asgardehs/heimdall
```

**As a CLI:**

```bash
go install github.com/asgardehs/heimdall/cmd/heimdall@latest
```

Requires Go 1.26 or later. No CGO — pure Go SQLite via `modernc.org/sqlite`,
so cross-compilation is trivial.

## Building from Source

```bash
git clone https://github.com/asgardehs/heimdall.git
cd heimdall
go build -o heimdall ./cmd/heimdall
```

## Documentation

Full documentation lives on the
[Asgard EHS docs site](https://asgardehs.github.io/docs/heimdall/):

- [Quick Start](https://asgardehs.github.io/docs/heimdall/quickstart/) — setup
  and first use
- [Go API](https://asgardehs.github.io/docs/heimdall/api/) — library reference
- [CLI Reference](https://asgardehs.github.io/docs/heimdall/cli/) — command
  line usage
- [Configuration Keys](https://asgardehs.github.io/docs/heimdall/keys/) —
  namespaces, types, and defaults

## Project

- **License:** GPL-3.0 — see [LICENSE](LICENSE)
- **Code of Conduct:** see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- **Contributing:** see [CONTRIBUTING.md](CONTRIBUTING.md)
- **Security:** report vulnerabilities to
  [muninn.developer@protonmail.com](mailto:muninn.developer@protonmail.com)

## Name

> _In Norse mythology, Heimdall is the watchful guardian of the Bifrost
> bridge, keeper of the horn Gjallarhorn. He sees all, hears all, and ensures
> the boundaries of Asgard are maintained. Here, Heimdall guards the shared
> configuration that binds the ecosystem together._

_Part of the [Asgard EHS family](https://asgardehs.github.io/)._
