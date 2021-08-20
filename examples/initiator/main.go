package main

import (
	"bytes"
	"context"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"github.com/b2broker/simplefix-go/utils"
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
		TagNotDefinedForMessageType: mustConvToInt(fixgen.EnumSessionRejectReasonTagNotDefinedForThisMessageType),
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
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", 9091))
	if err != nil {
		panic(fmt.Errorf("could not dial: %s", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := simplefixgo.NewInitiatorHandler(ctx, fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, handler, 10)

	handler.OnConnect(func() bool {
		return true
	})

	sess, err := session.NewInitiatorSession(
		context.Background(),
		handler,
		&pseudoGeneratedOpts,
		&session.LogonSettings{
			TargetCompID:  "Server",
			SenderCompID:  "Client",
			HeartBtInt:    5,
			EncryptMethod: fixgen.EnumEncryptMethodNoneother,
			Password:      "password",
			Username:      "login",
		},
	)
	if err != nil {
		panic(err)
	}

	handler.HandleIncoming(fixgen.MsgTypeLogon, func(msg []byte) {
		incomingLogon, err := fixgen.ParseLogon(msg)
		_, _ = incomingLogon, err
	})

	handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
		fmt.Println("incoming", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
	})
	handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
		fmt.Println("outgoing", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
	})

	sess.OnChangeState(utils.EventLogon, func() bool {
		err := sess.Send(fixgen.NewMarketDataRequest(
			"test",
			fixgen.EnumSubscriptionRequestTypeSnapshot,
			20,
			fixgen.NewMDEntryTypesGrp(),
			fixgen.NewRelatedSymGrp().
				AddEntry(fixgen.NewRelatedSymEntry().SetInstrument(fixgen.NewInstrument().SetSymbol("BTC/USDT"))).
				AddEntry(fixgen.NewRelatedSymEntry().SetInstrument(fixgen.NewInstrument().SetSymbol("ETH/USDT"))),
		))
		if err != nil {
			panic(err)
		}

		return true
	})

	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("resend request after 10 seconds")
		_ = sess.Send(fixgen.ResendRequest{}.New().SetFieldBeginSeqNo(2).SetFieldEndSeqNo(3))
	}()

	_ = sess.Run()

	panic(client.Serve())
}
