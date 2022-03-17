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

func TestGroup_Parse(t *testing.T) {
	var testMsg = []byte("8=FIX.4.49=23635=A34=149=sender56=target52=20210208-15:51:43.000262=1263=1264=20267=2269=0269=1146=355=BTC/USD864=2865=1868=put865=2868=call55=ETH/USD864=2865=1868=put865=2868=call55=KGB/FBI864=2865=1868=put865=2868=call10=051")

	msg := Items{
		NewKeyValue(beginString, NewString("FIX.4.4")),
		NewKeyValue(bodyLength, NewInt(236)),
		NewKeyValue(msgType, NewString("A")),
		NewKeyValue(msgSeqNum, NewInt(1)),
		NewKeyValue(senderCompID, NewString("sender")),
		NewKeyValue(targetCompID, NewString("target")),
		NewKeyValue(sendingTime, NewTime(time.Date(2021, 2, 8, 15, 51, 43, 0, time.UTC))),
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
		NewKeyValue(checksum, NewString("051")),
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
