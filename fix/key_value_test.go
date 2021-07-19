package fix

import (
	"bytes"
	"testing"
)

func TestKeyValues_ToBytes(t *testing.T) {
	keyValues := KeyValues{
		{"8", NewString("FIX.4.4")},
		{"9", NewString("408")},
		{"35", NewString("W")},
		{"49", NewString("LMAXD-MD")},
		{"56", NewString("b2brokerdigmdUATMTF")},
		{"34", NewString("111")},
		{"52", NewString("20190213-17:41:10.200")},
		{"22", NewString("8")},
		{"48", NewString("5005")},
		{"10", NewString("012")},
	}.ToBytes()

	msg := "8=FIX.4.49=40835=W49=LMAXD-MD56=b2brokerdigmdUATMTF34=11152=20190213-17:41:10.20022=848=500510=012"

	if !bytes.Equal(keyValues, []byte(msg)) {
		t.Log(string(keyValues))
		t.Log(msg)
		t.Fatalf("not equal")
	}
}

func TestKeyValue_ToBytes(t *testing.T) {
	keyValue := NewKeyValue("8", NewString("FIX.4.4")).ToBytes()

	msg := "8=FIX.4.4"

	if !bytes.Equal(keyValue, []byte(msg)) {
		t.Log(string(keyValue))
		t.Log(msg)
		t.Fatalf("not equal")
	}
}
