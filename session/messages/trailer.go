package messages

import "github.com/b2broker/simplefix-go/fix"

// TrailerBuilder is an interface providing functionality to a builder of Trailer messages.
type TrailerBuilder interface {
	New() TrailerBuilder

	AsComponent() *fix.Component
}
