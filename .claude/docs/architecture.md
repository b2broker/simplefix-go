# Architecture

## Layer Diagram

```
┌─────────────────────────────────────────┐
│  Session (state machine, heartbeats)    │
├─────────────────────────────────────────┤
│  Handler/Router (DefaultHandler, Pools) │
├─────────────────────────────────────────┤
│  Encoding (Unmarshaller, validation)    │
├─────────────────────────────────────────┤
│  FIX Message (Message, KV, Group, etc)  │
├─────────────────────────────────────────┤
│  Transport (Acceptor, Initiator, Conn)  │
└─────────────────────────────────────────┘
```

## Transport Layer

- **Acceptor** ([acceptor.go](../../acceptor.go)) — TCP server. Listens on port, spawns one `AcceptorHandler` per connection via `HandlerFactory`.
- **Initiator** ([initiator.go](../../initiator.go)) — TCP client. Connects to server, uses `InitiatorHandler` for message routing.
- **Conn** ([conn.go](../../conn.go)) — wraps `net.Conn`. Splits incoming byte stream on `10=` (checksum delimiter) for message framing. Manages read/write channels.

## Handler/Router Layer

- **DefaultHandler** ([handler.go](../../handler.go)) — central message router with incoming and outgoing `HandlerPool`s.
- **HandlerPool** ([handler_func_pool.go](../../handler_func_pool.go)) — map of message type → handler functions. Supports `AllMsgTypes` ("ALL") for catch-all subscriptions.

## FIX Message Model

- **Message** ([fix/message.go](../../fix/message.go)) — composite: beginString, bodyLength, msgType, header, body, trailer, checkSum
- **KeyValue** ([fix/key_value.go](../../fix/key_value.go)) — fundamental Tag=Value pair
- **Group** ([fix/group.go](../../fix/group.go)) — repeating groups with template and instances
- **Component** — reusable message components (Header, Trailer, etc.)
- **Value interface** ([fix/types.go](../../fix/types.go)) — implementations: String, Int, Float, Raw for serialization/deserialization

## Session State Machine

States defined in [session/session.go](../../session/session.go):

```
WaitingLogon → SuccessfulLogged ←→ WaitingTestReqAnswer
     ↓              ↓
WaitingLogonAnswer  WaitingLogoutAnswer → ReceivedLogoutAnswer → Disconnect
```

Session validates incoming messages, enforces heartbeat intervals, handles sequence numbering, and manages logon/logout lifecycle. Requires 8 message builders injected via options.

## Encoding

**DefaultUnmarshaller** ([fix/encoding/unmarshaler.go](../../fix/encoding/unmarshaler.go)) — decodes byte arrays into message structures. Supports strict mode for schema validation.

## Storage

**memory.Storage** ([storages/memory/storage.go](../../storages/memory/storage.go)) — in-memory implementation for CounterStorage (sequence numbers) and MessageStorage (message history for resend requests).
