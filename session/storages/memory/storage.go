package memory

import (
	"fmt"
	"github.com/b2broker/simplefix-go/session"
	"sync"
)

// Storage stores last messages
type Storage struct {
	maxSize    int
	bufferSize int
	start      int
	messages   map[int][]byte
	mu         sync.Mutex
}

// NewStorage is a constructor of in-memory Storage
// bufferSize is a size of additional messages at storage
// the larger the buffer, the less frequent of the flush
// maxSize is a count of messages in the storage
func NewStorage(maxSize int, bufferSize int) *Storage {
	return &Storage{
		maxSize:    maxSize,
		bufferSize: bufferSize,
		start:      1,
		messages:   make(map[int][]byte, maxSize+bufferSize),
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
			return fmt.Errorf("sequence index does not exit: %d", i)
		}
		delete(s.messages, i)
	}
	s.start = cutEnd

	return nil
}

func (s *Storage) Save(msg []byte, msgSeqNum int) error {
	if err := s.flush(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	expectedID := s.start + len(s.messages)
	if expectedID != msgSeqNum {
		return fmt.Errorf("%w: %d, expect: %d", session.ErrInvalidSequence, msgSeqNum, expectedID)
	}

	if _, ok := s.messages[msgSeqNum]; ok {
		return fmt.Errorf("sequence index alreasy exists: %d", msgSeqNum)
	}

	s.messages[msgSeqNum] = msg

	return nil
}

// Messages returns message list in sequential order from msgSeqNumFrom to msgSeqNumTo
func (s *Storage) Messages(msgSeqNumFrom, msgSeqNumTo int) ([][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msgSeqNumFrom > msgSeqNumTo {
		return nil, session.ErrInvalidBoundaries
	}

	if s.start > msgSeqNumFrom || s.start+len(s.messages) < msgSeqNumTo {
		return nil, session.ErrNotEnoughMessages
	}

	var messages [][]byte
	for i := msgSeqNumFrom; i <= msgSeqNumTo; i++ {
		messages = append(messages, s.messages[i])
	}

	return messages, nil
}
