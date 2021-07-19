package fix

import (
	"bytes"
	"testing"
)

func TestMessageToBytes(t *testing.T) {
	message := NewMessage("8", "9", "10", "35", "FIX.4.4", "V").
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
		)

	byteMessage, err := message.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal msg: %s", err)
	}
	origin := "8=FIX.4.49=14735=V34=249=ISLD52=20190313-16:45:21.86156=TW146=355=BTC/USD55=BTC/USDT_ABCDE55=BTCABCD/ABCDEFG262=request_1263=1264=5267=2269=0269=110=159"

	if !bytes.Equal(byteMessage, []byte(origin)) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(origin), origin)

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(origin), []byte(origin))

		t.Fatalf("not equal")
	}
}

func TestMessage_FromBytes(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=23635=A34=149=sender56=target52=20210208-15:51:43.000262=1263=1264=20267=2269=0269=1146=355=BTC/USD864=2865=1868=put865=2868=call55=ETH/USD864=2865=1868=put865=2868=call55=KGB/FBI864=2865=1868=put865=2868=call10=051")

	msg := NewMessageFromBytes(beginString, bodyLength, checksum, msgType)

	msg.
		SetHeader(
			NewComponent(
				NewKeyValue(msgSeqNum, &Int{}),
				NewKeyValue(senderCompID, &String{}),
				NewKeyValue(targetCompID, &String{}),
				NewKeyValue(sendingTime, &Time{}),
			),
		).
		SetTrailer(NewComponent()).
		SetBody(
			NewKeyValue(mdReqID, &Int{}),
			NewKeyValue(subscriptionRequestType, &String{}),
			NewKeyValue(marketDepth, &Int{}),
			NewGroup(noMDEntryTypes,
				NewComponent(
					NewKeyValue(mdEntryType, &String{}),
				),
			),
			NewGroup(noRelatedSym,
				NewComponent(
					NewKeyValue(symbol, &String{}),
					NewGroup(noEvents,
						NewKeyValue(eventType, &String{}),
						NewKeyValue(eventText, &String{}),
						NewKeyValue(eventDate, &Time{}),
					),
				),
			),
		)

	err := msg.Unmarshal(testMsg)
	if err != nil {
		t.Fatal(err)
	}

	byteMessage, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal msg: %s", err)
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

	msg := NewMessageFromBytes(beginString, bodyLength, checksum, msgType)

	msg.
		SetHeader(
			NewComponent(
				NewKeyValue("48", &String{}),
				NewKeyValue("56", &String{}),
				NewKeyValue("115", &String{}),
				NewKeyValue("122", &Time{}),
				NewKeyValue("91", &String{}),
				NewKeyValue("108", &String{}),
			),
		).
		SetTrailer(NewComponent())

	err := msg.Unmarshal(testMsg)
	if err != nil {
		t.Fatal(err)
	}

	byteMessage, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal msg: %s", err)
	}

	if !bytes.Equal(byteMessage, testMsg) {
		t.Log(len(byteMessage), string(byteMessage))
		t.Log(len(testMsg), string(testMsg))

		t.Log(len(byteMessage), byteMessage)
		t.Log(len(testMsg), testMsg)

		t.Fatalf("not equal")
	}
}
