package session

import (
	"fmt"
	"github.com/b2broker/simplefix-go/session/messages"
)

type MessageBuilders struct {
	HeaderBuilder              messages.HeaderBuilder
	TrailerBuilder             messages.TrailerBuilder
	LogonBuilder               messages.LogonBuilder
	LogoutBuilder              messages.LogoutBuilder
	RejectBuilder              messages.RejectBuilder
	HeartbeatBuilder           messages.HeartbeatBuilder
	TestRequestBuilder         messages.TestRequestBuilder
	ResendRequestBuilder       messages.ResendRequestBuilder
	SequenceResetBuilder       messages.SequenceResetBuilder
	ExecutionReportBuilder     messages.ExecutionReportBuilder
	NewOrderSingleBuilder      messages.NewOrderSingleBuilder
	MarketDataRequestBuilder   messages.MarketDataRequestBuilder
	OrderCancelRequestBuilder  messages.OrderCancelRequestBuilder
	SecurityListRequestBuilder messages.SecurityListRequestBuilder
}

// Opts is a structure providing auto-generated Session options.
type Opts struct {
	Location                string
	MessageBuilders         MessageBuilders
	Tags                    *messages.Tags
	AllowedEncryptedMethods map[string]struct{} // Can only be of the "None" type.
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
		return fmt.Errorf("%w: trailer", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.HeartbeatBuilder == nil {
		return fmt.Errorf("%w: heartbeat", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.ResendRequestBuilder == nil {
		return fmt.Errorf("%w: resend request", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.TestRequestBuilder == nil {
		return fmt.Errorf("%w: test request", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.LogoutBuilder == nil {
		return fmt.Errorf("%w: logout", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.LogonBuilder == nil {
		return fmt.Errorf("%w: logon", ErrMissingMessageBuilder)
	}

	if opts.MessageBuilders.RejectBuilder == nil {
		return fmt.Errorf("%w: reject", ErrMissingMessageBuilder)
	}

	if opts.Tags == nil {
		return ErrMissingRequiredTag
	}

	if opts.SessionErrorCodes == nil {
		return ErrMissingErrorCodes
	}

	return nil
}
