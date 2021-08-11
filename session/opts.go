package session

import (
	"fmt"
	"github.com/b2broker/simplefix-go/session/messages"
)

type MessageBuilders struct {
	HeaderBuilder        messages.HeaderBuilder
	TrailerBuilder       messages.TrailerBuilder
	LogonBuilder         messages.LogonBuilder
	LogoutBuilder        messages.LogoutBuilder
	RejectBuilder        messages.RejectBuilder
	HeartbeatBuilder     messages.HeartbeatBuilder
	TestRequestBuilder   messages.TestRequestBuilder
	ResendRequestBuilder messages.ResendRequestBuilder
}

// Opts is options for Session from generated code
type Opts struct {
	Location                string
	MessageBuilders         MessageBuilders
	Tags                    *messages.Tags
	AllowedEncryptedMethods map[string]struct{} // can contain only None type
	SessionErrorCodes       *messages.SessionErrorCodes
}

type Side int64

const (
	sideAcceptor = iota
	sideInitiator
)

func (opts *Opts) validate() error {
	if opts == nil {
		return ErrMissingSessionOts
	}

	if opts.MessageBuilders.HeaderBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.TrailerBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.HeartbeatBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.ResendRequestBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.TestRequestBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.LogoutBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.LogonBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.RejectBuilder == nil {
		return fmt.Errorf("%w: header", ErrMissingMessageBuilder)
	}

	if opts.Tags == nil {
		return ErrMissingRequiredTag
	}

	if opts.SessionErrorCodes == nil {
		return ErrMissingErrorCodes
	}

	return nil
}
