---
path_scope: "**/*_test.go"
---

# Testing

## Integration test helpers

- `RunAcceptor(port, t)` — creates server, starts listening, returns Acceptor + address. See [tests/acceptor.go](../../tests/acceptor.go).
- `RunNewInitiator(addr, t, settings)` — connects client, performs logon, returns Session + DefaultHandler. See [tests/initiator.go](../../tests/initiator.go).

## Full Acceptor/Initiator integration tests

Tests in [tests/acceptor_initiator_test.go](../../tests/acceptor_initiator_test.go):
- `TestHeartbeat` — verifies heartbeat exchange between server and client
- `TestGroup` — verifies repeating group encoding/decoding with nested fields
- Uses `utils.TimedWaitGroup` for timeout-aware synchronization

## MockMessage for isolated handler tests

[session/messages/mock_message.go](../../session/messages/mock_message.go) — implements `SendingMessage` with configurable Type, Data, Err fields. Use for handler tests without full message infrastructure.

## Storage: always use memory.NewStorage()

Tests use in-memory implementation from [storages/memory/](../../storages/memory/) for both `CounterStorage` and `MessageStorage`. No external dependencies needed.

## Unit test locations

- `fix/message_test.go` — message assembly and serialization
- `fix/types_test.go` — Value type serialization (9 tests)
- `fix/convertor_test.go` — type conversions
- `fix/key_value_test.go` — KeyValue operations
- `fix/group_test.go` — repeating group assembly
- `fix/encoding/unmarshaller_test.go` — binary deserialization with various field orderings
- `session/session_test.go` — initialization validation (missing builders, invalid limits)
- `utils/timed_wg_test.go` — TimedWaitGroup
- `utils/timer_test.go` — timer utilities
- `generator/generator_test.go` — code generation pipeline
