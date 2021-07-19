package fix

import (
	"bytes"
	"fmt"
	"strconv"
)

const TimeLayout = "20060102-15:04:05.000"

const CountOfSOHSymbols = 3

var Delimiter = []byte{1}

func joinBody(values ...[]byte) []byte {
	return bytes.Join(values, Delimiter)
}

func makeTagValue(tag string, value []byte) []byte {
	return bytes.Join([][]byte{
		[]byte(tag), value,
	}, []byte{61})
}

// nolint
func makeGroup(entries []map[string][]byte, tags []string) []byte {
	var groupItems [][]byte
	for _, entry := range entries {
		for _, tag := range tags {
			groupItems = append(groupItems, bytes.Join([][]byte{[]byte(tag), entry[tag]}, []byte{61}))
		}
	}

	return bytes.Join(groupItems, Delimiter)
}

func calcCheckSum(body []byte) []byte {
	var sum int
	for _, b := range body {
		sum += int(b)
	}
	sum += int(byte(1))

	return []byte(fmt.Sprintf("%03s", strconv.Itoa(sum%256)))
}
