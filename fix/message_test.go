package fix

import (
	"bytes"
	"testing"
	"time"
)

func TestMessage_ToBytes(t *testing.T) {
	testCases := []struct {
		name    string
		message *Message
		want    []byte
	}{
		{
			name: "Full Message",
			message: NewMessage("8", "9", "10", "35", "FIX.4.4", "V").
				SetHeader(NewComponent(
					NewKeyValue("34", NewString("2")),
					NewKeyValue("49", NewString("ISLD")),
					NewKeyValue("52", NewString("20190313-16:45:21.861")),
					NewKeyValue("56", NewString("TW")),
				)).
				SetBody(
					NewKeyValue("146", NewString("3")),
					NewKeyValue("55", NewString("BTC/USD")),
					NewKeyValue("55", NewString("BTC/USDT_ABCDE")),
					NewKeyValue("55", NewString("BTCABCD/ABCDEFG")),
					NewKeyValue("262", NewString("request_1")),
					NewKeyValue("263", NewString("1")),
					NewKeyValue("264", NewString("5")),
					NewKeyValue("267", NewString("2")),
					NewKeyValue("269", NewString("0")),
					NewKeyValue("269", NewString("1")),
				),
			want: []byte("8=FIX.4.49=14735=V34=249=ISLD52=20190313-16:45:21.86156=TW146=355=BTC/USD55=BTC/USDT_ABCDE55=BTCABCD/ABCDEFG262=request_1263=1264=5267=2269=0269=110=159"),
		},
		{
			name: "Empty Header",
			message: NewMessage("8", "9", "10", "35", "FIX.4.4", "V").
				SetHeader(NewComponent()).
				SetBody(
					NewKeyValue("146", NewString("3")),
					NewKeyValue("55", NewString("BTC/USD")),
					NewKeyValue("55", NewString("BTC/USDT_ABCDE")),
					NewKeyValue("55", NewString("BTCABCD/ABCDEFG")),
					NewKeyValue("262", NewString("request_1")),
					NewKeyValue("263", NewString("1")),
					NewKeyValue("264", NewString("5")),
					NewKeyValue("267", NewString("2")),
					NewKeyValue("269", NewString("0")),
					NewKeyValue("269", NewString("1")),
				),
			want: []byte("8=FIX.4.49=10335=V146=355=BTC/USD55=BTC/USDT_ABCDE55=BTCABCD/ABCDEFG262=request_1263=1264=5267=2269=0269=110=188"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			byteMessage, err := tt.message.ToBytes()
			if err != nil {
				t.Fatalf("could not marshal message: %s", err)
			}

			if !bytes.Equal(byteMessage, tt.want) {
				t.Log(len(byteMessage), string(byteMessage))
				t.Log(len(tt.want), string(tt.want))

				t.Log(len(byteMessage), byteMessage)
				t.Log(len(tt.want), tt.want)

				t.Fatalf("not equal")
			}

			converter := NewMessageByteConverter(200)
			byteMessage, err = converter.ConvertToBytes(tt.message)
			if err != nil {
				t.Fatalf("could not marshal message: %s", err)
			}

			if !bytes.Equal(byteMessage, tt.want) {
				t.Log(len(byteMessage), string(byteMessage))
				t.Log(len(tt.want), string(tt.want))

				t.Log(len(byteMessage), byteMessage)
				t.Log(len(tt.want), tt.want)

				t.Fatalf("not equal")
			}
		})
	}
}

func TestMessage_FromBytes(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=23635=A34=149=sender56=target52=20210208-15:51:43.000262=1263=1264=20267=2269=0269=1146=355=BTC/USD864=2865=1868=put865=2868=call55=ETH/USD864=2865=1868=put865=2868=call55=KGB/FBI864=2865=1868=put865=2868=call10=051")

	msg := NewMessage(beginString, bodyLength, checksum, msgType, "FIX.4.4", "A")

	msg.
		SetHeader(
			NewComponent(
				NewKeyValue(msgSeqNum, NewInt(1)),
				NewKeyValue(senderCompID, NewString("sender")),
				NewKeyValue(targetCompID, NewString("target")),
				NewKeyValue("52", NewTime(time.Date(2021, 2, 8, 15, 51, 43, 0, time.UTC))),
			),
		).
		SetTrailer(NewComponent()).
		SetBody(
			NewKeyValue(mdReqID, NewInt(1)),
			NewKeyValue(subscriptionRequestType, NewString("1")),
			NewKeyValue(marketDepth, NewInt(20)),
			NewGroup(noMDEntryTypes, NewComponent(
				&KeyValue{Key: mdEntryType},
			)).
				AddEntry(NewComponent(
					NewKeyValue(mdEntryType, NewString("0")),
				).Items()).
				AddEntry(NewComponent(
					NewKeyValue(mdEntryType, NewString("1")),
				).Items()),
			NewGroup(noRelatedSym, NewComponent(&KeyValue{Key: symbol},
				NewGroup(noEvents, NewComponent(&KeyValue{Key: eventType}, &KeyValue{Key: eventText}))),
			).
				AddEntry(NewComponent(
					NewKeyValue(symbol, NewString("BTC/USD")),
					NewGroup(noEvents, NewComponent(&KeyValue{Key: eventType}, &KeyValue{Key: eventText})).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("1")),
							NewKeyValue(eventText, NewString("put")),
						).Items()).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("2")),
							NewKeyValue(eventText, NewString("call")),
						).Items()),
				).Items()).
				AddEntry(NewComponent(
					NewKeyValue(symbol, NewString("ETH/USD")),
					NewGroup(noEvents, NewComponent(&KeyValue{Key: eventType}, &KeyValue{Key: eventText})).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("1")),
							NewKeyValue(eventText, NewString("put")),
						).Items()).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("2")),
							NewKeyValue(eventText, NewString("call")),
						).Items()),
				).Items()).
				AddEntry(NewComponent(
					NewKeyValue(symbol, NewString("KGB/FBI")),
					NewGroup(noEvents, NewComponent(&KeyValue{Key: eventType}, &KeyValue{Key: eventText})).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("1")),
							NewKeyValue(eventText, NewString("put")),
						).Items()).
						AddEntry(NewComponent(
							NewKeyValue(eventType, NewString("2")),
							NewKeyValue(eventText, NewString("call")),
						).Items()),
				).Items()),
		)

	byteMessage, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal message: %s", err)
	}

	if !bytes.Equal(byteMessage, testMsg) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(testMsg), string(testMsg))

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(testMsg), testMsg)

		t.Fatalf("not equal")
	}

	converter := NewMessageByteConverter(200)
	byteMessage, err = converter.ConvertToBytes(msg)
	if err != nil {
		t.Fatalf("could not marshal message: %s", err)
	}

	if !bytes.Equal(byteMessage, testMsg) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(testMsg), string(testMsg))

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(testMsg), testMsg)

		t.Fatalf("not equal")
	}
}

func TestMessage_FromBytes_Coincidence(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=5935=A56=client115=server122=20210305-15:16:58.263108=3010=191")

	msg := NewMessage(beginString, bodyLength, checksum, msgType, "FIX.4.4", "A")

	msg.
		SetHeader(
			NewComponent(
				NewKeyValue("48", NewString("")),
				NewKeyValue("56", NewString("client")),
				NewKeyValue("115", NewString("server")),
				NewKeyValue("122", NewTime(time.Date(2021, 3, 05, 15, 16, 58, 263000000, time.UTC))),
				NewKeyValue("91", NewString("")),
				NewKeyValue("108", NewString("30")),
			),
		).
		SetTrailer(NewComponent())

	byteMessage, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal message: %s", err)
	}

	if !bytes.Equal(byteMessage, testMsg) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(testMsg), string(testMsg))

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(testMsg), testMsg)

		t.Fatalf("not equal")
	}

	converter := NewMessageByteConverter(200)
	byteMessage, err = converter.ConvertToBytes(msg)
	if err != nil {
		t.Fatalf("could not marshal message: %s", err)
	}

	if !bytes.Equal(byteMessage, testMsg) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(testMsg), string(testMsg))

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(testMsg), testMsg)

		t.Fatalf("not equal")
	}
}

func TestMessage_CalcBodyLength(t *testing.T) {
	testCases := []struct {
		name    string
		message *Message
		want    int
	}{
		{
			name: "base",
			message: NewMessage("8", "9", "10", "35", "FIX4.4", "A").
				SetHeader(NewComponent()).
				SetTrailer(NewComponent()),
			want: 5,
		},
		{
			name: "headers only",
			message: NewMessage("8", "9", "10", "35", "FIX4.4", "A").
				SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
				SetTrailer(NewComponent()),
			want: 21,
		},
		{
			name: "body only",
			message: NewMessage("8", "9", "10", "35", "FIX4.4", "A").
				SetHeader(NewComponent()).
				SetBody(NewKeyValue("100", NewString("test"))).
				SetTrailer(NewComponent()),
			want: 14,
		},
		{
			name: "headers and body",
			message: NewMessage("8", "9", "10", "35", "FIX4.4", "A").
				SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
				SetBody(NewKeyValue("100", NewString("test"))).
				SetTrailer(NewComponent()),
			want: 30,
		},
		{
			name: "headers, body and checksum",
			message: NewMessage("8", "9", "10", "35", "FIX4.4", "A").
				SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
				SetBody(NewKeyValue("100", NewString("test"))).
				SetTrailer(NewComponent(NewKeyValue("10", NewString("000")))),
			want: 30,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			converter := NewMessageByteConverter(200)
			_, _ = converter.ConvertToBytes(tt.message)
			actual := tt.message.bodyLength.Value.Value().(int)
			if tt.want != actual {
				t.Logf("body length: expected %d, got %d", tt.want, actual)
				t.FailNow()
			}
		})
	}
}

// BenchmarkMessage_Prepare-24    	 1263648	       982.7 ns/op
func BenchmarkMessage_Prepare(b *testing.B) {
	msg := NewMessage("8", "9", "10", "35", "FIX4.4", "A").
		SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
		SetBody(NewKeyValue("100", NewString("test"))).
		SetTrailer(NewComponent())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = msg.Prepare()
	}
}

var testMsg = NewMessage("8", "9", "10", "35", "FIX4.4", "A").
	SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
	SetBody(NewKeyValue("100", NewString("test"))).
	SetTrailer(NewComponent())
var testConverter = NewMessageByteConverter(500)

// BenchmarkMessage_PrepareBuffered-24    	 5583252	       211.5 ns/op
func BenchmarkMessage_PrepareBuffered(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = testConverter.ConvertToBytes(testMsg)
	}
}
