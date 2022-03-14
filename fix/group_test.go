package fix

import (
	"bytes"
	"testing"
	"time"
)

const (
	beginString  = "8"
	bodyLength   = "9"
	checksum     = "10"
	msgType      = "35"
	msgSeqNum    = "34"
	senderCompID = "49"
	targetCompID = "56"
	sendingTime  = "52"

	mdReqID                 = "262"
	subscriptionRequestType = "263"
	marketDepth             = "264"

	noMDEntryTypes = "267"
	mdEntryType    = "269"

	noRelatedSym = "146"
	symbol       = "55"

	noEvents  = "864"
	eventType = "865"
	eventDate = "866"
	eventText = "868"
)

func newHeader(msgSeqNumV int, senderCompIDV, targetCompIDV string, sendingTimeV time.Time) *Component {
	return NewComponent(
		NewKeyValue(msgSeqNum, NewInt(msgSeqNumV)),
		NewKeyValue(senderCompID, NewString(senderCompIDV)),
		NewKeyValue(targetCompID, NewString(targetCompIDV)),
		NewKeyValue(sendingTime, NewTime(sendingTimeV)),
	)
}

func TestGroup_AddItem(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=23635=A34=149=sender56=target52=20210208-12:51:43.000262=1263=1264=20267=2269=0269=1146=355=BTC/USD864=2865=1868=put865=2868=call55=ETH/USD864=2865=1868=put865=2868=call55=KGB/FBI864=2865=1868=put865=2868=call10=048")

	var (
		beginStringValue = "FIX.4.4"
		msgLogon         = "A"

		sender = "sender"
		target = "target"
	)

	entryTypes := NewGroup(noMDEntryTypes,
		NewKeyValue(mdEntryType, &String{}),
	)

	relatedSym := NewGroup(noRelatedSym,
		NewKeyValue(symbol, &String{}),
	)

	msg := NewMessage(beginString, bodyLength, checksum, msgType, beginStringValue, msgLogon).
		SetBody(
			NewKeyValue(mdReqID, NewString("1")),
			NewKeyValue(subscriptionRequestType, NewString("1")),
			NewKeyValue(marketDepth, NewString("20")),
			entryTypes,
			relatedSym,
		).
		SetHeader(newHeader(1, sender, target, time.Unix(1612788703, 0).UTC()))

	entryTypes.AddEntry(Items{
		NewKeyValue(mdEntryType, NewString("0")),
	})
	entryTypes.AddEntry(Items{
		NewKeyValue(mdEntryType, NewString("1")),
	})

	makeInstrumentComponent := func(sym string) *Component {
		events := NewGroup(noEvents,
			NewKeyValue(eventType, &String{}),
			NewKeyValue(eventText, &String{}),
		)
		events.AddEntry(Items{
			NewKeyValue(eventType, NewString("1")),
			NewKeyValue(eventText, NewString("put")),
		})

		events.AddEntry(Items{
			NewKeyValue(eventType, NewString("2")),
			NewKeyValue(eventText, NewString("call")),
		})

		instrument := NewComponent(
			NewKeyValue(symbol, NewString(sym)),
			events,
		)

		return instrument
	}

	relatedSym.AddEntry(Items{
		makeInstrumentComponent("BTC/USD"),
	})
	relatedSym.AddEntry(Items{
		makeInstrumentComponent("ETH/USD"),
	})
	relatedSym.AddEntry(Items{
		makeInstrumentComponent("KGB/FBI"),
	})

	res, err := msg.ToBytes()
	if err != nil {
		t.Fatalf("could not marshal message: %s", err)
	}
	if !bytes.Equal(testMsg, res) {
		t.Log(len(testMsg), string(testMsg))
		t.Log(len(res), string(res))
		t.Fatalf("message length is not equal")
	}
}

type TestGroup struct {
	*Group
}

func makeTestGroup() TestGroup {
	return TestGroup{
		Group: NewGroup(noRelatedSym,
			NewComponent(
				NewKeyValue(symbol, &String{}),
				NewGroup(noEvents,
					NewKeyValue(eventType, &String{}),
					NewKeyValue(eventText, &String{}),
					NewKeyValue(eventDate, &Time{}),
				),
			),
		),
	}
}

func TestGroup_Parse(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=23635=A34=149=sender56=target52=20210208-15:51:43.000262=1263=1264=20267=2269=0269=1146=355=BTC/USD864=2865=1868=put865=2868=call55=ETH/USD864=2865=1868=put865=2868=call55=KGB/FBI864=2865=1868=put865=2868=call10=051")

	msg := Items{
		NewKeyValue(beginString, &String{}),
		NewKeyValue(bodyLength, &Int{}),
		NewKeyValue(msgType, &String{}),
		NewKeyValue(msgSeqNum, &Int{}),
		NewKeyValue(senderCompID, &String{}),
		NewKeyValue(targetCompID, &String{}),
		NewKeyValue(sendingTime, &Time{}),
		NewKeyValue(mdReqID, &Int{}),
		NewKeyValue(subscriptionRequestType, &String{}),
		NewKeyValue(marketDepth, &Int{}),
		NewGroup(noMDEntryTypes,
			NewComponent(
				NewKeyValue(mdEntryType, &String{}),
			),
		),
		makeTestGroup().Group,
		NewKeyValue("10", &String{}),
	}

	err := UnmarshalItems(testMsg, msg, true)
	if err != nil {
		panic(err)
	}

	res := msg.ToBytes()
	if !bytes.Equal(testMsg, res) {
		if !bytes.Equal(testMsg, res) {
			t.Log(len(testMsg), string(testMsg))
			t.Log(len(res), string(res))
			t.Fatalf("message length is not equal")
		}
	}
}
