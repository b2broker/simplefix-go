package fix

import (
	"bytes"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMessageByteConverter_ConvertToBytes(t *testing.T) {
	conv := NewMessageByteConverter(100)
	msg := NewMessage("8", "9", "10", "35", "FIX4.4", "A").
		SetHeader(NewComponent(NewKeyValue("49", NewString("test")), NewKeyValue("56", NewString("test")))).
		SetBody(NewKeyValue("100", NewString("test"))).
		SetTrailer(NewComponent())

	expected := []byte(`8=FIX4.49=3035=A49=test56=test100=test10=023`)
	b, err := conv.ConvertToBytes(msg)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if !bytes.Equal(expected, b) {
		t.Fatalf("Expected %v, got %v", expected, string(b))
	}
}

func TestMessageByteConverter_Concurrent(t *testing.T) {
	conv := NewMessageByteConverter(100)
	requests := make(chan int, 100)
	wg := sync.WaitGroup{}
	workers := 100

	for workerID := 0; workerID < workers; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			workerKey := strconv.Itoa(workerID)
			for i := range requests {
				time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
				key := strconv.Itoa(i)
				msg := NewMessage("8", "9", "10", "35", "FIX4.4", "A").
					SetHeader(NewComponent(NewKeyValue("worker", NewString(workerKey)), NewKeyValue("messageID", NewString(key)))).
					SetBody(NewKeyValue("100", NewString("test"))).
					SetTrailer(NewComponent())
				expected := []byte(`8=FIX4.49=3035=Aworker=` + workerKey + `messageID=` + key + `100=test10=023`)
				b, err := conv.ConvertToBytes(msg)
				if err != nil {
					t.Errorf("Error: %s", err)
				}
				// here we replace length and checksum to ignore them
				bb := bytes.Split(b, []byte{1})
				expectedb := bytes.Split(expected, []byte{1})
				bb[1] = expectedb[1]
				bb[len(bb)-2] = expectedb[len(expectedb)-2]
				b = bytes.Join(bb, []byte{1})
				if !bytes.Equal(expected, b) {
					t.Errorf("Expected %v,\n got %v", bb, expectedb)
				}
			}
		}(workerID)
	}

	for i := 0; i < workers*100; i++ {
		requests <- i
	}
	close(requests)
	wg.Wait()
}
