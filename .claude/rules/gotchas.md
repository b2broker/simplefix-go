# Gotchas

## All 8 message builders are mandatory for Session

Session construction requires ALL 8 builders: Header, Trailer, Logon, Logout, Reject, Heartbeat, TestRequest, ResendRequest. Missing any one fails with `ErrMissingMessageBuilder`. See [session/opts.go](../../session/opts.go).

## Conn framing reads until `10=` — messages without checksum hang

`Conn` reads TCP stream until the `10=` (checksum) tag delimiter appears. Any message not ending with a checksum tag causes an indefinite hang. See [conn.go](../../conn.go).

## Tags are string dependencies — mismatches silently break routing

Session requires tag numbers as integers in `messages.Tags` (MsgType=35, MsgSeqNum=34, etc.). A mismatch between actual message tags and the Tags struct silently breaks message routing and sequence tracking. See [session/messages/contstants.go](../../session/messages/contstants.go) (filename typo is in the source).

## `AllMsgTypes` constant is the string "ALL"

Handler subscription uses `"ALL"` as a catch-all. If any FIX message type is literally named "ALL", behavior is undefined. See [handler.go](../../handler.go).

## HandlerPool.Remove() clears entire type list

`Remove()` ignores the handler ID parameter and only deletes the msgType entry when the handlers list is empty, never actually removing individual handlers. Handler IDs are for tracking, not removal. See [handler_func_pool.go](../../handler_func_pool.go).

## Generated code is overwritten on regeneration

`fixgen` writes to the output directory without checking for existing files. Custom edits to generated files in `tests/fix44/` will be lost. Never edit generated files.

## Session state is atomic but transitions are not

Session state uses `atomic.Int64`, but field mutations during state transitions are not atomic with the state change. A reader may see state X with stale field values.

## Message builders are reused — reset state between messages

If builder state is not reset between messages, old field values from the previous message may leak. Visible in test patterns at [tests/acceptor_initiator_test.go](../../tests/acceptor_initiator_test.go).

## Context cancellation race in DefaultHandler

`DefaultHandler` has dual context tracking: its own `ctx` and `handler.Context()`. Simultaneous cancellation of both can cause goroutine leaks. See [initiator.go](../../initiator.go).
