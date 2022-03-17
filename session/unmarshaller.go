package session

import (
	"github.com/b2broker/simplefix-go/session/messages"
)

type Unmarshaller interface {
	Unmarshal(msg messages.Builder, d []byte) error
}
