package main

import (
	"bytes"
	"context"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	"github.com/b2broker/simplefix-go/session/storages/memory"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"net"
	"strconv"
	"time"
)

func mustConvToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

// todo move boilerplate to generator
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
		InvalidTagNumber:            mustConvToInt(fixgen.EnumSessionRejectReasonInvalidtagnumber),
		RequiredTagMissing:          mustConvToInt(fixgen.EnumSessionRejectReasonRequiredtagmissing),
		TagNotDefinedForMessageType: mustConvToInt(fixgen.EnumSessionRejectReasonTagnotdefinedforthismessagetype),
		UndefinedTag:                mustConvToInt(fixgen.EnumSessionRejectReasonUndefinedtag),
		TagSpecialWithoutValue:      mustConvToInt(fixgen.EnumSessionRejectReasonTagspecifiedwithoutavalue),
		IncorrectValue:              mustConvToInt(fixgen.EnumSessionRejectReasonValueisincorrectoutofrangeforthistag),
		IncorrectDataFormatValue:    mustConvToInt(fixgen.EnumSessionRejectReasonIncorrectdataformatforvalue),
		DecryptionProblem:           mustConvToInt(fixgen.EnumSessionRejectReasonDecryptionproblem),
		SignatureProblem:            mustConvToInt(fixgen.EnumSessionRejectReasonSignatureproblem),
		CompIDProblem:               mustConvToInt(fixgen.EnumSessionRejectReasonCompidproblem),
		Other:                       mustConvToInt(fixgen.EnumSessionRejectReasonOther),
	},
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 9091))
	if err != nil {
		panic(err)
	}

	handlerFactory := simplefixgo.NewAcceptorHandlerFactory(fixgen.FieldMsgType, 10)

	server := simplefixgo.NewAcceptor(listener, handlerFactory, func(handler simplefixgo.AcceptorHandler) {
		sess, err := session.NewAcceptorSession(
			context.Background(),
			&pseudoGeneratedOpts,
			handler,
			&session.LogonSettings{
				HeartBtInt:   30,
				LogonTimeout: time.Second * 30,
				HeartBtLimits: &session.IntLimits{
					Min: 5,
					Max: 60,
				},
			},
			func(request *session.LogonSettings) (err error) {
				fmt.Printf(
					"free logon for '%s' (%s)\n",
					request.Username,
					request.Password,
				)

				return nil
			},
		)
		if err != nil {
			panic(err)
		}

		_ = sess.Run()
		sess.SetMessageStorage(memory.NewStorage(100, 100))

		handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
			fmt.Println("incoming", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
		})
		handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
			fmt.Println("outgoing", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
		})
	})

	panic(fmt.Errorf("server was stopped: %s", server.ListenAndServe()))
}
