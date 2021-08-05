package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session"
	fixgen "github.com/b2broker/simplefix-go/tests/fix44"
	"net"
	"testing"
	"time"
)

func RunNewInitiator(port int, t *testing.T, settings *session.LogonSettings) (s *session.Session, handler *simplefixgo.DefaultHandler) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("could not dial: %s", err)
	}

	handler = simplefixgo.NewInitiatorHandler(context.Background(), fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, handler, 10)

	s, err = session.NewInitiatorSession(
		context.Background(),
		handler,
		&pseudoGeneratedOpts,
		settings,
	)
	if err != nil {
		panic(err)
	}

	// log messages
	handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) {
		fmt.Println("incoming:", string(bytes.ReplaceAll(msg, fix.Delimiter, []byte("|"))))
	})
	handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg []byte) {
		fmt.Println("outgoing:", string(bytes.ReplaceAll(msg, fix.Delimiter, []byte("|"))))
	})

	// todo move
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("resend request after 10 seconds")
		err = s.Send(fixgen.ResendRequest{}.New().SetFieldBeginSeqNo(2).SetFieldEndSeqNo(3))
		if err != nil {
			panic(err)
		}
	}()

	err = s.Run()
	if err != nil {
		t.Fatalf("run session: %s", err)
	}

	go func() {
		err = client.Serve()
		if err != nil && !errors.Is(err, simplefixgo.ErrConnClosed) {
			panic(fmt.Errorf("serve client: %s", err))
		}
	}()

	return s, handler
}
