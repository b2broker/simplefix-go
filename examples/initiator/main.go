package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/b2broker/simplefix-go/storages/memory"
	"net"
	"strconv"
	"time"

	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/fix/encoding"
	"github.com/b2broker/simplefix-go/session"
	"github.com/b2broker/simplefix-go/session/messages"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"github.com/b2broker/simplefix-go/utils"
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
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", 9991))
	if err != nil {
		panic(fmt.Errorf("could not dial: %s", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := simplefixgo.NewInitiatorHandler(ctx, fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, handler, 10, time.Second*5)

	handler.OnConnect(func() bool {
		return true
	})

	exampleStorage := memory.NewStorage()

	sess, err := session.NewInitiatorSession(
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
		exampleStorage,
		exampleStorage,
	)
	if err != nil {
		panic(err)
	}

	handler.HandleIncoming(fixgen.MsgTypeLogon, func(msg []byte) bool {
		incomingLogon := fixgen.NewLogon()
		err := encoding.Unmarshal(incomingLogon, msg)
		_, _ = incomingLogon, err
		return true
	})

	handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) bool {
		fmt.Println("incoming", string(bytes.ReplaceAll(msg, fix.Delimiter, []byte("|"))))
		return true
	})
	handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		data, err := msg.ToBytes()
		if err != nil {
			panic(err)
		}
		fmt.Println("outgoing", string(bytes.ReplaceAll(data, fix.Delimiter, []byte("|"))))
		return true
	})

	sess.OnChangeState(utils.EventLogon, func() bool {
		err := sess.Send(fixgen.CreateMarketDataRequest(
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
