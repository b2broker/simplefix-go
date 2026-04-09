# simplefix-go

Pure Go implementation of the FIX (Financial Information eXchange) protocol. Provides building blocks for FIX clients (Initiators) and servers (Acceptors) with session management, message validation, custom FIX dialect support, and code generation from XML schemas.

- **Module:** `github.com/b2broker/simplefix-go`
- **Go version:** 1.22 (go.mod minimum); CI tests with Go 1.24
- **Default branch:** `master`
- **Dependencies:** `golang.org/x/sync` only

## Packages

| Package | Description | Entry |
|---------|-------------|-------|
| root | Transport layer — Acceptor, Initiator, Conn, DefaultHandler | [acceptor.go](../acceptor.go), [initiator.go](../initiator.go) |
| `fix` | Message primitives — Message, KeyValue, Group, Component, Value types | [fix/message.go](../fix/message.go) |
| `fix/encoding` | Message unmarshaling from bytes with optional validation | [fix/encoding/unmarshaler.go](../fix/encoding/unmarshaler.go) |
| `fix/buffer` | Message byte buffer management | [fix/buffer/buffer.go](../fix/buffer/buffer.go) |
| `session` | FIX session state machine — logon/logout, heartbeats, sequence management | [session/session.go](../session/session.go) |
| `session/messages` | Message builder interfaces and standard FIX control message defs | [session/messages/header.go](../session/messages/header.go) |
| `generator` | Code generator from XML schemas → Go types, builders, enums | [generator/generator.go](../generator/generator.go) |
| `cmd/fixgen` | CLI tool for code generation | [cmd/fixgen/main.go](../cmd/fixgen/main.go) |
| `storages/memory` | In-memory storage for message sequences and counters | [storages/memory/storage.go](../storages/memory/storage.go) |
| `utils` | Event handling, timers, utilities | [utils/](../utils/) |

## Architecture

**Layered (bottom to top):**

1. **Transport** — `Acceptor` (server), `Initiator` (client), `Conn` (TCP wrapper with FIX framing on `10=` delimiter)
2. **Handler/Router** — `DefaultHandler` with `HandlerPool` for incoming/outgoing message routing by type
3. **Message** — `Message` composite (header + body + trailer + checksum), `KeyValue` pairs, `Group` repeating groups
4. **Encoding** — `DefaultUnmarshaller` for byte→message deserialization with optional strict validation
5. **Session** — State machine: WaitingLogon → SuccessfulLogged → (heartbeat/test/reject cycle) → Disconnect
6. **Generator** — `fixgen` CLI reads XML schemas, outputs typed message builders, field constants, enums

## Commands

```bash
go test -race ./...                                                    # run all tests
go vet ./...                                                           # vet
golangci-lint run                                                      # lint (uses defaults, no .golangci.yml)
go run ./cmd/fixgen -o=./fix44 -s=./source/fix44.xml -t=./source/types.xml  # generate FIX types from XML
```

No Makefile. CI runs via GitHub Actions ([.github/workflows/go.yml](../.github/workflows/go.yml)): lint → test → semantic-release on master.

## Key Interfaces

| Interface | Purpose | Location |
|-----------|---------|----------|
| `SendingMessage` | Any outgoing message (MsgType + Marshal) | [handler.go](../handler.go) |
| `Sender` | Send(SendingMessage) error | [acceptor.go](../acceptor.go) |
| `AcceptorHandler` | Server message handler with lifecycle | [acceptor.go](../acceptor.go) |
| `InitiatorHandler` | Client message handler | [initiator.go](../initiator.go) |
| `Value` | Field value serialization contract | [fix/types.go](../fix/types.go) |
| `HeaderBuilder`, `LogonBuilder`, etc. | Message-specific builders (8 required) | [session/messages/](../session/messages/) |

## Reference Docs

- [docs/architecture.md](docs/architecture.md) — detailed layer descriptions, state machine, message model
- [docs/generator.md](docs/generator.md) — fixgen XML schema format, generation pipeline
