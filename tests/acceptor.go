package tests

import (
	"bytes"
	"context"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"net"
	"testing"
	"time"
)

func RunAcceptor(port int, t *testing.T, storage session.MessageStorage) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("listen error: %s", err)
	}

	handlerFactory := simplefixgo.NewAcceptorHandlerFactory(fixgen.FieldMsgType, 10)

	server := simplefixgo.NewAcceptor(listener, handlerFactory, func(handler simplefixgo.AcceptorHandler) {
		s := session.NewAcceptorSession(
			context.Background(),
			PseudoGeneratedOpts,
			handler,
			session.LogonSettings{
				HeartBtInt:   30,
				LogonTimeout: time.Second * 30,
				HeartBtLimits: &session.IntLimits{
					Min: 5,
					Max: 60,
				},
			},
			func(request session.LogonSettings) (err error) {
				t.Logf(
					"user '%s' connected with password '%s'",
					request.Username,
					request.Password,
				)

				return nil
			},
		)
		err = s.Run()
		if err != nil {
			t.Fatalf("run s: %s", err)
		}

		s.SetMessageStorage(storage)

		// log messages
		handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
			t.Log("incoming:", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
		})
		handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
			t.Log("outgoing:", string(bytes.Replace(msg, fix.Delimiter, []byte("|"), -1)))
		})
	})

	go server.ListenAndServe()
}
