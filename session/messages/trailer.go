package messages

import "github.com/b2broker/simplefix-go/fix"

// TrailerBuilder is an interface for Trailer message builder
type TrailerBuilder interface {
	New() TrailerBuilder

	AsComponent() *fix.Component
}
