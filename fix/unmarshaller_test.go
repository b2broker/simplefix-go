package fix

import (
	"bytes"
	"testing"
)

const visibleDelimiter = "|"

func TestUnmarshalItems(t *testing.T) {
	var exampleRawMas = []byte("8=FIX.4.49=6135=A49=Client56=Server34=152=20210706-19:06:12.838108=510=206")

	testData := []struct {
		raw    []byte
		expect []byte
		items  Items
	}{
		{
			raw: exampleRawMas,
			items: Items{
				&KeyValue{Key: "8", Value: &String{}},
				&KeyValue{Key: "9", Value: &String{}},
				&KeyValue{Key: "35", Value: &String{}},
				&KeyValue{Key: "49", Value: &String{}},
				&KeyValue{Key: "56", Value: &String{}},
				&KeyValue{Key: "34", Value: &String{}},
				&KeyValue{Key: "52", Value: &String{}},
				&KeyValue{Key: "108", Value: &String{}},
				&KeyValue{Key: "10", Value: &String{}},
			},
			expect: exampleRawMas,
		},
		{
			raw: exampleRawMas,
			items: Items{
				&KeyValue{Key: "8", Value: &String{}},
				&KeyValue{Key: "35", Value: &String{}},
				&KeyValue{Key: "9", Value: &String{}},
				&KeyValue{Key: "56", Value: &String{}},
				&KeyValue{Key: "49", Value: &String{}},
				&KeyValue{Key: "34", Value: &String{}},
				&KeyValue{Key: "52", Value: &String{}},
				&KeyValue{Key: "100500", Value: &String{}},
				&KeyValue{Key: "108", Value: &String{}},
				&KeyValue{Key: "10", Value: &String{}},
			},
			expect: []byte("8=FIX.4.435=A9=6156=Server49=Client34=152=20210706-19:06:12.838108=510=206"),
		},
		{
			raw: exampleRawMas,
			items: Items{
				&KeyValue{Key: "8", Value: &String{}},
				&KeyValue{Key: "9", Value: &String{}},
				&KeyValue{Key: "35", Value: &String{}},
				&KeyValue{Key: "56", Value: &String{}},
				&KeyValue{Key: "49", Value: &String{}},
				&KeyValue{Key: "34", Value: &String{}},
				&KeyValue{Key: "108", Value: &String{}},
				&KeyValue{Key: "52", Value: &String{}},
				&KeyValue{Key: "10", Value: &String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server49=Client34=1108=552=20210706-19:06:12.83810=206"),
		},
		{
			raw: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
			items: Items{
				&KeyValue{Key: "8", Value: &String{}},
				&KeyValue{Key: "9", Value: &String{}},
				&KeyValue{Key: "35", Value: &String{}},
				&KeyValue{Key: "56", Value: &String{}},
				NewGroup("161", &KeyValue{Key: "162", Value: &String{}}, &KeyValue{Key: "163", Value: &String{}}),
				&KeyValue{Key: "10", Value: &String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
		},
		{
			raw: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
			items: Items{
				&KeyValue{Key: "8", Value: &String{}},
				&KeyValue{Key: "9", Value: &String{}},
				&KeyValue{Key: "35", Value: &String{}},
				&KeyValue{Key: "56", Value: &String{}},
				NewGroup("161", NewComponent(&KeyValue{Key: "162", Value: &String{}}, &KeyValue{Key: "163", Value: &String{}})),
				&KeyValue{Key: "10", Value: &String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
		},
	}

	for i, item := range testData {
		err := UnmarshalItems(item.raw, item.items, true)
		if err != nil {
			t.Fatal(err)
		}

		res := item.items.ToBytes()
		if !bytes.Equal(item.expect, res) {
			t.Logf("%d. expect %s (%d)", i, showDelimiter(item.expect), len(item.expect))
			t.Logf("%d. result %s (%d)", i, showDelimiter(res), len(res))
			t.Fatal("result doesnt equal expected message")
		}
	}

}

func showDelimiter(in []byte) []byte {
	return bytes.Replace(in, Delimiter, []byte(visibleDelimiter), -1)
}
