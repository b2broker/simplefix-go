package tests

import (
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
)

var PseudoGeneratedOpts = session.SessionOpts{
	LogonBuilder:         fixgen.Logon{}.New(),
	LogoutBuilder:        fixgen.Logout{}.New(),
	RejectBuilder:        fixgen.Reject{}.New(),
	HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
	TestRequestBuilder:   fixgen.TestRequest{}.New(),
	ResendRequestBuilder: fixgen.ResendRequest{}.New(),
	Tags: messages.Tags{
		MsgType:         mustConvToInt(fixgen.FieldMsgType),
		MsgSeqNum:       mustConvToInt(fixgen.FieldMsgSeqNum),
		HeartBtInt:      mustConvToInt(fixgen.FieldHeartBtInt),
		EncryptedMethod: mustConvToInt(fixgen.FieldEncryptMethod),
	},
	AllowedEncryptedMethods: map[string]struct{}{
		fixgen.EnumEncryptMethodNoneother: {},
	},
	SessionErrorCodes: messages.SessionErrorCodes{
		RequiredTagMissing: 1,
		IncorrectValue:     5,
		Other:              99,
	},
}
