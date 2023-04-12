package encoding

import (
	"fmt"
	"github.com/b2broker/simplefix-go/session/messages"
)

type DefaultValidator struct{}

func (v DefaultValidator) Do(msg messages.Builder) error {
	return v.checkRequiredFields(msg)
}

func (DefaultValidator) checkRequiredFields(msg messages.Builder) error {
	if msg.BeginString().Value.IsNull() {
		return fmt.Errorf("the required field value is empty: BeginString")
	}
	if msg.BodyLength() == 0 {
		return fmt.Errorf("the required field value is empty: BodyLength")
	}
	if msg.MsgType() == "" {
		return fmt.Errorf("the required field value is empty: MsgType")
	}
	if msg.CheckSum() == "" {
		return fmt.Errorf("the required field value is empty: CheckSum")
	}

	return nil
}
