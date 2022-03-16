package session

import (
	"github.com/b2broker/simplefix-go/fix/encoding"
)

type Unmarshaller interface {
	Unmarshal(msg encoding.MessageBuilder, d []byte) error
}
