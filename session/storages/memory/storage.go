package memory

import (
	"fmt"
	simplefixgo "github.com/b2broker/simplefix-go"
	"github.com/b2broker/simplefix-go/session"
	"sync"
)

// Storage is used to store the most recent messages.
type Storage struct {
	maxSize    int
	bufferSize int
	start      int
	messages   map[int]simplefixgo.SendingMessage
	mu         sync.Mutex
}

// NewStorage is a constructor for creation of a new in-memory Storage.
// bufferSize determines the size of additional messages that can be put to the storage
// (the larger the buffer size, the less frequently the storage is flushed).
// maxSize specifies the maximum number of messages the storage can contain.
func NewStorage(maxSize int, bufferSize int) *Storage {
	return &Storage{
		maxSize:    maxSize,
		bufferSize: bufferSize,
		start:      1,
		messages:   make(map[int]simplefixgo.SendingMessage, maxSize+bufferSize),
		mu:         sync.Mutex{},
	}
}

func (s *Storage) flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.messages) < s.maxSize+s.bufferSize {
		return nil
	}

	cutEnd := s.start + s.bufferSize
	for i := s.start; i < cutEnd; i++ {
		if _, ok := s.messages[i]; !ok {
			return fmt.Errorf("The sequence index is not found: %d", i)
		}
		delete(s.messages, i)
	}
	s.start = cutEnd

	return nil
}

func (s *Storage) Save(msg simplefixgo.SendingMessage, msgSeqNum int) error {
	if err := s.flush(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	expectedID := s.start + len(s.messages)
	if expectedID != msgSeqNum {
		return fmt.Errorf("%w: %d, expected value: %d", session.ErrInvalidSequence, msgSeqNum, expectedID)
	}

	if _, ok := s.messages[msgSeqNum]; ok {
		return fmt.Errorf("The sequence index already exists: %d", msgSeqNum)
	}

	s.messages[msgSeqNum] = msg

	return nil
}

// Messages returns a message list, in a sequential order
// (starting with msgSeqNumFrom and ending with msgSeqNumTo).
func (s *Storage) Messages(msgSeqNumFrom, msgSeqNumTo int) ([]simplefixgo.SendingMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msgSeqNumFrom > msgSeqNumTo {
		return nil, session.ErrInvalidBoundaries
	}

	if s.start > msgSeqNumFrom || s.start+len(s.messages) < msgSeqNumTo {
		return nil, session.ErrNotEnoughMessages
	}

	var messages []simplefixgo.SendingMessage
	for i := msgSeqNumFrom; i <= msgSeqNumTo; i++ {
		messages = append(messages, s.messages[i])
	}

	return messages, nil
}
