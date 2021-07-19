package messages

import "github.com/b2broker/simplefix-go/fix"

type TrailerBuilder interface {
	New() TrailerBuilder

	AsComponent() *fix.Component
}
