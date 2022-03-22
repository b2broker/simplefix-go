package messages

type Logout interface {
	New() LogoutBuilder
	Build() LogoutBuilder
}

// LogoutBuilder is an interface providing functionality to a builder of auto-generated Logout messages.
type LogoutBuilder interface {
	Logout
	PipelineBuilder
}
