package session

import (
	"github.com/b2broker/simplefix-go/session/messages"
)

type Opts struct {
	Location string

	HeaderBuilder  messages.HeaderBuilder
	TrailerBuilder messages.TrailerBuilder

	LogonBuilder         messages.LogonBuilder
	LogoutBuilder        messages.LogoutBuilder
	RejectBuilder        messages.RejectBuilder
	HeartbeatBuilder     messages.HeartbeatBuilder
	TestRequestBuilder   messages.TestRequestBuilder
	ResendRequestBuilder messages.ResendRequestBuilder

	Tags                    messages.Tags
	AllowedEncryptedMethods map[string]struct{} // can be only None
	SessionErrorCodes       messages.SessionErrorCodes
}

type Side int64

const (
	sideAcceptor = iota
	sideInitiator
)
