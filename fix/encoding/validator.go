package encoding

import (
	"fmt"
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

type DefaultValidator struct{}

func (v DefaultValidator) Do(msg messages.Builder) error {
	if err := v.checkRequiredFields(msg); err != nil {
		return err
	}

	bodyLength := msg.CalcBodyLength()
	if bodyLength != msg.BodyLength() {
		return fmt.Errorf("an invalid body length; specified: %d, required: %d",
			msg.BodyLength(),
			bodyLength,
		)
	}

	checkSum := fix.CalcCheckSum(msg.BytesWithoutChecksum())

	if string(checkSum) != msg.CheckSum() {
		return fmt.Errorf("an invalid checksum; specified: %s, required: %s", msg.CheckSum(), checkSum)
	}

	return nil
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
