package memory

import (
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/session"
	"sync"
	"sync/atomic"
)

// Storage is used to store the most recent messages.
type Storage struct {
	counter  int64
	messages map[int]simplefixgo.SendingMessage
	mu       sync.Mutex
}

// NewStorage is a constructor for creation of a new in-memory Storage.
func NewStorage() *Storage {
	return &Storage{
		messages: map[int]simplefixgo.SendingMessage{},
		mu:       sync.Mutex{},
	}
}

func (s *Storage) GetNextSeqNum(pk string) (int, error) {
	return int(atomic.AddInt64(&s.counter, 1)), nil
}

func (s *Storage) GetCurrSeqNum(pk string) (int, error) {
	return int(s.counter), nil
}

func (s *Storage) ResetSeqNum(pk string) error {
	return nil
}

// Save saves a message with seq number to storage
func (s *Storage) Save(pk string, msg simplefixgo.SendingMessage, msgSeqNum int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[msgSeqNum]; ok {
		return fmt.Errorf("the sequence index already exists: %d", msgSeqNum)
	}
	s.messages[msgSeqNum] = msg
	return nil
}

// Messages returns a message list, in a sequential order
// (starting with msgSeqNumFrom and ending with msgSeqNumTo).
func (s *Storage) Messages(pk string, msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msgSeqNumFrom > msgSeqNumTo {
		return nil, session.ErrInvalidBoundaries
	}

	if int64(msgSeqNumTo) > s.counter {
		return nil, session.ErrNotEnoughMessages
	}

	var messages []simplefixgo.SendingMessage
	for i := msgSeqNumFrom; i <= msgSeqNumTo; i++ {
		if _, ok := s.messages[i]; !ok {
			return nil, session.ErrNotEnoughMessages
		}
		messages = append(messages, s.messages[i])
	}

	return messages, nil
}
