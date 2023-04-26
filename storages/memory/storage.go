package memory

import (
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/fix"
	"sync"
	"sync/atomic"
)

// Storage is used to store the most recent messages.
type Storage struct {
	counterIncoming int64
	counterOutgoing int64
	messages        map[int]simplefixgo.SendingMessage
	mu              sync.Mutex
}

// NewStorage is a constructor for creation of a new in-memory Storage.
func NewStorage() *Storage {
	return &Storage{
		messages: map[int]simplefixgo.SendingMessage{},
		mu:       sync.Mutex{},
	}
}

func (s *Storage) GetNextSeqNum(storageID fix.StorageID) (int, error) {
	if storageID.Side == fix.Incoming {
		return int(atomic.AddInt64(&s.counterIncoming, 1)), nil
	} else {
		return int(atomic.AddInt64(&s.counterOutgoing, 1)), nil
	}
}

func (s *Storage) GetCurrSeqNum(storageID fix.StorageID) (int, error) {
	if storageID.Side == fix.Incoming {
		return int(s.counterIncoming), nil
	} else {
		return int(s.counterOutgoing), nil
	}
}

func (s *Storage) ResetSeqNum(storageID fix.StorageID) error {
	if storageID.Side == fix.Incoming {
		s.counterIncoming = 0
	} else {
		s.counterOutgoing = 0
	}
	return nil
}

func (s *Storage) SetSeqNum(storageID fix.StorageID, seqNum int) error {
	if storageID.Side == fix.Incoming {
		s.counterIncoming = int64(seqNum)
	} else {
		s.counterOutgoing = int64(seqNum)
	}
	return nil
}

// Save saves a message with seq number to storage
func (s *Storage) Save(_ fix.StorageID, msg simplefixgo.SendingMessage, msgSeqNum int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[msgSeqNum] = msg
	return nil
}

// Messages returns a message list, in a sequential order
// (starting with msgSeqNumFrom and ending with msgSeqNumTo).
func (s *Storage) Messages(_ fix.StorageID, msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msgSeqNumFrom > msgSeqNumTo {
		return nil, simplefixgo.ErrInvalidBoundaries
	}

	if int64(msgSeqNumTo) > s.counterOutgoing {
		return nil, simplefixgo.ErrNotEnoughMessages
	}

	var sendingMessages []simplefixgo.SendingMessage
	for i := msgSeqNumFrom; i <= msgSeqNumTo; i++ {
		if _, ok := s.messages[i]; !ok {
			return nil, simplefixgo.ErrNotEnoughMessages
		}
		sendingMessages = append(sendingMessages, s.messages[i])
	}

	return sendingMessages, nil
}
