package messages

type Heartbeat interface {
	New() HeartbeatBuilder
	TestReqID() string
	SetFieldTestReqID(string) HeartbeatBuilder
	HeaderBuilder() HeaderBuilder
}

// HeartbeatBuilder is an interface providing functionality to a builder of auto-generated Heartbeat messages.
type HeartbeatBuilder interface {
	Heartbeat
	PipelineBuilder
}
