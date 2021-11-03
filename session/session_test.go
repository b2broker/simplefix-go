package session

import (
	"context"
	"errors"
	"testing"
	"time"

	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/session/messages"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
)

var (
	validLocations     = []string{"UTC", "Moscow/Europe"}
	validLogonSettings = LogonSettings{
		TargetCompID:  "turret",
		SenderCompID:  "sander",
		HeartBtInt:    30,
		EncryptMethod: "post",
		Password:      "sword",
		Username:      "user",
		LogonTimeout:  time.Millisecond,
		HeartBtLimits: &IntLimits{
			Min: 1,
			Max: 30,
		},
		CloseTimeout: time.Millisecond,
	}
	validTags = &messages.Tags{
		MsgType:         35,
		MsgSeqNum:       123,
		HeartBtInt:      321,
		EncryptedMethod: 555,
	}
	validMessageBuilders = MessageBuilders{
		HeaderBuilder:        &fixgen.Header{},
		TrailerBuilder:       &fixgen.Trailer{},
		LogonBuilder:         &fixgen.Logon{},
		LogoutBuilder:        &fixgen.Logout{},
		RejectBuilder:        &fixgen.Reject{},
		HeartbeatBuilder:     &fixgen.Heartbeat{},
		TestRequestBuilder:   &fixgen.TestRequest{},
		ResendRequestBuilder: &fixgen.ResendRequest{},
	}
	validSessionErrorCodes     = &messages.SessionErrorCodes{}
	validEncryptedMethod       = map[string]struct{}{"test": {}}
	invalidMessageBuildersList = []MessageBuilders{
		{
			TrailerBuilder:       &fixgen.Trailer{},
			LogonBuilder:         &fixgen.Logon{},
			LogoutBuilder:        &fixgen.Logout{},
			RejectBuilder:        &fixgen.Reject{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			LogonBuilder:         &fixgen.Logon{},
			LogoutBuilder:        &fixgen.Logout{},
			RejectBuilder:        &fixgen.Reject{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			TrailerBuilder:       &fixgen.Trailer{},
			LogoutBuilder:        &fixgen.Logout{},
			RejectBuilder:        &fixgen.Reject{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			TrailerBuilder:       &fixgen.Trailer{},
			LogonBuilder:         &fixgen.Logon{},
			RejectBuilder:        &fixgen.Reject{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			TrailerBuilder:       &fixgen.Trailer{},
			LogonBuilder:         &fixgen.Logon{},
			LogoutBuilder:        &fixgen.Logout{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			TrailerBuilder:       &fixgen.Trailer{},
			LogonBuilder:         &fixgen.Logon{},
			LogoutBuilder:        &fixgen.Logout{},
			RejectBuilder:        &fixgen.Reject{},
			TestRequestBuilder:   &fixgen.TestRequest{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:        &fixgen.Header{},
			TrailerBuilder:       &fixgen.Trailer{},
			LogonBuilder:         &fixgen.Logon{},
			LogoutBuilder:        &fixgen.Logout{},
			RejectBuilder:        &fixgen.Reject{},
			HeartbeatBuilder:     &fixgen.Heartbeat{},
			ResendRequestBuilder: &fixgen.ResendRequest{},
		},
		{
			HeaderBuilder:      &fixgen.Header{},
			TrailerBuilder:     &fixgen.Trailer{},
			LogonBuilder:       &fixgen.Logon{},
			LogoutBuilder:      &fixgen.Logout{},
			RejectBuilder:      &fixgen.Reject{},
			HeartbeatBuilder:   &fixgen.Heartbeat{},
			TestRequestBuilder: &fixgen.TestRequest{},
		},
	}
)

func TestNewAcceptorSession(t *testing.T) {
	ctx := context.Background()

	session, err := NewAcceptorSession(&Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, simplefixgo.NewAcceptorHandler(ctx, "35", 100),
		&validLogonSettings,
		func(request *LogonSettings) (err error) { return nil },
	)
	if err != nil {
		t.Fatalf("unexpected behaviour, got error: %v", err)
	}

	if err := session.Stop(); err != nil {
		t.Fatalf("unexpected behaviour, got error: %v", err)
	}

	session.changeState(WaitingLogoutAnswer)
	session.changeState(ReceivedLogoutAnswer)
}

func TestNewInitiatorSession(t *testing.T) {
	ctx := context.Background()

	session, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, &validLogonSettings)
	if err != nil {
		t.Fatalf("unexpected behaviour, got error: %v", err)
	}

	if err := session.Stop(); err != nil {
		t.Fatalf("unexpected behaviour, got error: %v", err)
	}

	session.changeState(WaitingLogoutAnswer)
	session.changeState(ReceivedLogoutAnswer)
}

func TestNewInitiatorSessionOpts(t *testing.T) {
	ctx := context.Background()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), nil,
		&validLogonSettings)
	if !errors.Is(err, ErrMissingSessionOts) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingSessionOts, err)
	}
}

func TestNewAcceptorSessionMessageBuilders(t *testing.T) {
	ctx := context.Background()

	for _, invalidMessageBuilders := range invalidMessageBuildersList {
		_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
			Location:                validLocations[0],
			MessageBuilders:         invalidMessageBuilders,
			Tags:                    validTags,
			AllowedEncryptedMethods: validEncryptedMethod,
			SessionErrorCodes:       validSessionErrorCodes,
		}, &validLogonSettings)
		if !errors.Is(err, ErrMissingMessageBuilder) {
			t.Fatalf("unexpected behaveour, expect error: %s, got: %v", ErrMissingMessageBuilder, err)
		}
	}
}

func TestNewInitiatorSessionHandler(t *testing.T) {
	_, err := NewInitiatorSession(nil, &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, &validLogonSettings)
	if !errors.Is(err, ErrMissingHandler) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingHandler, err)
	}
}

func TestNewInitiatorSessionTags(t *testing.T) {
	ctx := context.Background()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    nil,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, &validLogonSettings)
	if !errors.Is(err, ErrMissingRequiredTag) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingRequiredTag, err)
	}
}

func TestNewInitiatorSessionEncryptedMethod(t *testing.T) {
	ctx := context.Background()

	_, err := NewAcceptorSession(&Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: nil,
		SessionErrorCodes:       validSessionErrorCodes,
	}, simplefixgo.NewInitiatorHandler(ctx, "35", 100), &validLogonSettings,
		func(request *LogonSettings) (err error) {
			return nil
		},
	)
	if !errors.Is(err, ErrMissingEncryptedMethods) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingEncryptedMethods, err)
	}
}

func TestNewInitiatorSessionErrorCodes(t *testing.T) {
	ctx := context.Background()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       nil,
	}, &validLogonSettings)
	if !errors.Is(err, ErrMissingErrorCodes) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingErrorCodes, err)
	}
}

func TestNewInitiatorSessionLogonSettings(t *testing.T) {
	ctx := context.Background()

	invalid1 := validLogonSettings
	invalid1.EncryptMethod = ""

	invalid2 := validLogonSettings
	invalid2.HeartBtInt = 0

	cases := map[string]struct {
		settings LogonSettings
		err      error
	}{
		"nil encrypted method": {settings: invalid1, err: ErrMissingEncryptMethod},
		"zero heartbeat int":   {settings: invalid2, err: ErrInvalidHeartBtInt},
	}

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, nil)
	if !errors.Is(err, ErrMissingLogonSettings) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingLogonSettings, err)
	}

	for name, c := range cases {
		_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
			Location:                validLocations[0],
			MessageBuilders:         validMessageBuilders,
			Tags:                    validTags,
			AllowedEncryptedMethods: validEncryptedMethod,
			SessionErrorCodes:       validSessionErrorCodes,
		}, &c.settings)
		if !errors.Is(err, c.err) {
			t.Fatalf("unexpected behavior in case '%s', expect: %s, got: %v", name, c.err, err)
		}
	}
}

func TestNewAcceptorSessionLogonSettings(t *testing.T) {
	ctx := context.Background()

	invalid1 := validLogonSettings
	invalid1.HeartBtLimits = &IntLimits{
		Min: 0,
		Max: 0,
	}

	invalid2 := validLogonSettings
	invalid2.LogonTimeout = 0

	cases := map[string]struct {
		settings LogonSettings
		err      error
	}{
		"invalid heartbeat limits": {settings: invalid1, err: ErrInvalidHeartBtLimits},
		"logon request timeout":    {settings: invalid2, err: ErrInvalidLogonTimeout},
	}

	_, err := NewAcceptorSession(&Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, simplefixgo.NewInitiatorHandler(ctx, "35", 100), nil, func(request *LogonSettings) (err error) {
		return nil
	})
	if !errors.Is(err, ErrMissingLogonSettings) {
		t.Fatalf("unexpected error, expect: %s, got: %v", ErrMissingLogonSettings, err)
	}

	for name, c := range cases {
		_, err := NewAcceptorSession(&Opts{
			Location:                validLocations[0],
			MessageBuilders:         validMessageBuilders,
			Tags:                    validTags,
			AllowedEncryptedMethods: validEncryptedMethod,
			SessionErrorCodes:       validSessionErrorCodes,
		}, simplefixgo.NewInitiatorHandler(ctx, "35", 100), &c.settings, func(request *LogonSettings) (err error) {
			return nil
		})
		if !errors.Is(err, c.err) {
			t.Fatalf("unexpected behavior in case '%s', expect: %s, got: %v", name, c.err, err)
		}
	}
}
