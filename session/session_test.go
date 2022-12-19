package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/b2broker/simplefix-go/storages/memory"

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
		HeaderBuilder:        fixgen.Header{}.New(),
		TrailerBuilder:       fixgen.Trailer{}.New(),
		LogonBuilder:         fixgen.Logon{}.New(),
		LogoutBuilder:        fixgen.Logout{}.New(),
		RejectBuilder:        fixgen.Reject{}.New(),
		HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
		TestRequestBuilder:   fixgen.TestRequest{}.New(),
		ResendRequestBuilder: fixgen.ResendRequest{}.New(),
	}
	validSessionErrorCodes     = &messages.SessionErrorCodes{}
	validEncryptedMethod       = map[string]struct{}{"test": {}}
	invalidMessageBuildersList = []MessageBuilders{
		{
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			TestRequestBuilder:   fixgen.TestRequest{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:        fixgen.Header{}.New(),
			TrailerBuilder:       fixgen.Trailer{}.New(),
			LogonBuilder:         fixgen.Logon{}.New(),
			LogoutBuilder:        fixgen.Logout{}.New(),
			RejectBuilder:        fixgen.Reject{}.New(),
			HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
			ResendRequestBuilder: fixgen.ResendRequest{}.New(),
		},
		{
			HeaderBuilder:      fixgen.Header{}.New(),
			TrailerBuilder:     fixgen.Trailer{}.New(),
			LogonBuilder:       fixgen.Logon{}.New(),
			LogoutBuilder:      fixgen.Logout{}.New(),
			RejectBuilder:      fixgen.Reject{}.New(),
			HeartbeatBuilder:   fixgen.Heartbeat{}.New(),
			TestRequestBuilder: fixgen.TestRequest{}.New(),
		},
	}
)

func TestNewAcceptorSession(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	session, err := NewAcceptorSession(&Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, simplefixgo.NewAcceptorHandler(ctx, "35", 100),
		&validLogonSettings,
		func(request *LogonSettings) (err error) { return nil },
		testStorage,
		testStorage,
	)
	if err != nil {
		t.Fatalf("unexpected behavior, returned error: %v", err)
	}

	if err := session.Stop(); err != nil {
		t.Fatalf("unexpected behavior, returned error: %v", err)
	}

	session.changeState(WaitingLogoutAnswer, true)
	session.changeState(ReceivedLogoutAnswer, true)
}

func TestNewInitiatorSession(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	session, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	},
		&validLogonSettings,
		testStorage,
		testStorage,
	)
	if err != nil {
		t.Fatalf("unexpected behavior, returned error: %v", err)
	}

	if err := session.Stop(); err != nil {
		t.Fatalf("unexpected behavior, returned error: %v", err)
	}

	session.changeState(WaitingLogoutAnswer, true)
	session.changeState(ReceivedLogoutAnswer, true)
}

func TestNewInitiatorSessionOpts(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), nil,
		&validLogonSettings, testStorage, testStorage)
	if !errors.Is(err, ErrMissingSessionOts) {
		t.Fatalf("unexpected error, expected: %s, returned: %v", ErrMissingSessionOts, err)
	}
}

func TestNewAcceptorSessionMessageBuilders(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	for _, invalidMessageBuilders := range invalidMessageBuildersList {
		_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
			Location:                validLocations[0],
			MessageBuilders:         invalidMessageBuilders,
			Tags:                    validTags,
			AllowedEncryptedMethods: validEncryptedMethod,
			SessionErrorCodes:       validSessionErrorCodes,
		},
			&validLogonSettings,
			testStorage,
			testStorage,
		)
		if !errors.Is(err, ErrMissingMessageBuilder) {
			t.Fatalf("unexpected behavior, expected error: %s, received: %v", ErrMissingMessageBuilder, err)
		}
	}
}

func TestNewInitiatorSessionHandler(t *testing.T) {
	testStorage := memory.NewStorage()

	_, err := NewInitiatorSession(nil, &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	},
		&validLogonSettings,
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingHandler) {
		t.Fatalf("unexpected error, expected: %s, received: %v", ErrMissingHandler, err)
	}
}

func TestNewInitiatorSessionTags(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    nil,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	},
		&validLogonSettings,
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingRequiredTag) {
		t.Fatalf("unexpected error, expected: %s, received: %v", ErrMissingRequiredTag, err)
	}
}

func TestNewInitiatorSessionEncryptedMethod(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

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
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingEncryptedMethods) {
		t.Fatalf("unexpected error, expected: %s, received: %v", ErrMissingEncryptedMethods, err)
	}
}

func TestNewInitiatorSessionErrorCodes(t *testing.T) {
	ctx := context.Background()

	testStorage := memory.NewStorage()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       nil,
	},
		&validLogonSettings,
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingErrorCodes) {
		t.Fatalf("unexpected error, expected: %s, received: %v", ErrMissingErrorCodes, err)
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
		"encryption method is nil":   {settings: invalid1, err: ErrMissingEncryptMethod},
		"heartbeat interval is zero": {settings: invalid2, err: ErrInvalidHeartBtInt},
	}

	testStorage := memory.NewStorage()

	_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	},
		nil,
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingLogonSettings) {
		t.Fatalf("unexpected error, expected: %s, returned: %v", ErrMissingLogonSettings, err)
	}

	for name, c := range cases {
		_, err := NewInitiatorSession(simplefixgo.NewInitiatorHandler(ctx, "35", 100), &Opts{
			Location:                validLocations[0],
			MessageBuilders:         validMessageBuilders,
			Tags:                    validTags,
			AllowedEncryptedMethods: validEncryptedMethod,
			SessionErrorCodes:       validSessionErrorCodes,
		},
			&c.settings,
			testStorage,
			testStorage,
		)
		if !errors.Is(err, c.err) {
			t.Fatalf("unexpected behavior in case '%s', expected: %s, returned: %v", name, c.err, err)
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

	testStorage := memory.NewStorage()

	_, err := NewAcceptorSession(&Opts{
		Location:                validLocations[0],
		MessageBuilders:         validMessageBuilders,
		Tags:                    validTags,
		AllowedEncryptedMethods: validEncryptedMethod,
		SessionErrorCodes:       validSessionErrorCodes,
	}, simplefixgo.NewInitiatorHandler(ctx, "35", 100), nil, func(request *LogonSettings) (err error) {
		return nil
	},
		testStorage,
		testStorage,
	)
	if !errors.Is(err, ErrMissingLogonSettings) {
		t.Fatalf("unexpected error, expected: %s, returned: %v", ErrMissingLogonSettings, err)
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
		},
			testStorage,
			testStorage,
		)
		if !errors.Is(err, c.err) {
			t.Fatalf("unexpected behavior in case '%s', expected: %s, returned: %v", name, c.err, err)
		}
	}
}
