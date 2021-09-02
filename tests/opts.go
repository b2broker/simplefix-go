package tests

import (
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
)

var pseudoGeneratedOpts = session.Opts{
	MessageBuilders: session.MessageBuilders{
		HeaderBuilder:        fixgen.Header{}.New(),
		TrailerBuilder:       fixgen.Trailer{}.New(),
		LogonBuilder:         fixgen.Logon{}.New(),
		LogoutBuilder:        fixgen.Logout{}.New(),
		RejectBuilder:        fixgen.Reject{}.New(),
		HeartbeatBuilder:     fixgen.Heartbeat{}.New(),
		TestRequestBuilder:   fixgen.TestRequest{}.New(),
		ResendRequestBuilder: fixgen.ResendRequest{}.New(),
	},
	Tags: &messages.Tags{
		MsgType:         mustConvToInt(fixgen.FieldMsgType),
		MsgSeqNum:       mustConvToInt(fixgen.FieldMsgSeqNum),
		HeartBtInt:      mustConvToInt(fixgen.FieldHeartBtInt),
		EncryptedMethod: mustConvToInt(fixgen.FieldEncryptMethod),
	},
	AllowedEncryptedMethods: map[string]struct{}{
		fixgen.EnumEncryptMethodNoneother: {},
	},
	SessionErrorCodes: &messages.SessionErrorCodes{
		InvalidTagNumber:   mustConvToInt(fixgen.EnumSessionRejectReasonInvalidtagnumber),
		RequiredTagMissing: mustConvToInt(fixgen.EnumSessionRejectReasonRequiredtagmissing),
		//TagNotDefinedForMessageType: mustConvToInt(fixgen.EnumSessionRejectReasonTagNotDefinedForThisMessageType),
		UndefinedTag:             mustConvToInt(fixgen.EnumSessionRejectReasonUndefinedtag),
		TagSpecialWithoutValue:   mustConvToInt(fixgen.EnumSessionRejectReasonTagspecifiedwithoutavalue),
		IncorrectValue:           mustConvToInt(fixgen.EnumSessionRejectReasonValueisincorrectoutofrangeforthistag),
		IncorrectDataFormatValue: mustConvToInt(fixgen.EnumSessionRejectReasonIncorrectdataformatforvalue),
		DecryptionProblem:        mustConvToInt(fixgen.EnumSessionRejectReasonDecryptionproblem),
		SignatureProblem:         mustConvToInt(fixgen.EnumSessionRejectReasonSignatureproblem),
		CompIDProblem:            mustConvToInt(fixgen.EnumSessionRejectReasonCompidproblem),
		Other:                    mustConvToInt(fixgen.EnumSessionRejectReasonOther),
	},
}
