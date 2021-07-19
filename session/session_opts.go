package session

import (
	"github.com/b2broker/simplefix-go/session/messages"
)

type SessionOpts struct {
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
