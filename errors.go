package simplefixgo

import "errors"

var (
	ErrNotEnoughMessages = errors.New("not enough messages in the storage")
	ErrInvalidBoundaries = errors.New("invalid boundaries")
	ErrInvalidSequence   = errors.New("unexpected sequence index")
)
