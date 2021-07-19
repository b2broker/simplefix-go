package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

const MsgTypeTestRequest = "1"

type TestRequest struct {
	*fix.Message
}

func makeTestRequest() *TestRequest {
	msg := &TestRequest{
		Message: fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, beginString, MsgTypeTestRequest).
			SetBody(
				fix.NewKeyValue(FieldTestReqID, &fix.String{}),
			),
	}

	msg.SetHeader(makeHeader().AsComponent())
	msg.SetTrailer(makeTrailer().AsComponent())

	return msg
}

func NewTestRequest(header *Header, trailer *Trailer, testReqID string) *TestRequest {
	msg := makeTestRequest().
		SetTestReqID(testReqID)
	msg.SetHeader(header.AsComponent())
	msg.SetTrailer(trailer.AsComponent())
	return msg
}

func ParseTestRequest(data []byte) (*TestRequest, error) {
	msg := fix.NewMessage(FieldBeginString, FieldBodyLength, FieldCheckSum, FieldMsgType, FieldBeginString, beginString).
		SetBody(makeTestRequest().Body()...).
		SetHeader(makeHeader().AsComponent()).
		SetTrailer(makeTrailer().AsComponent())

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return &TestRequest{
		Message: msg,
	}, nil
}

func (testRequest *TestRequest) Header() *Header {
	header := testRequest.Message.Header()

	return &Header{header}
}

func (testRequest *TestRequest) HeaderBuilder() messages.HeaderBuilder {
	return testRequest.Header()
}

func (testRequest *TestRequest) Trailer() *Trailer {
	trailer := testRequest.Message.Trailer()

	return &Trailer{trailer}
}

func (testRequest *TestRequest) TestReqID() string {
	kv := testRequest.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (testRequest *TestRequest) SetTestReqID(testReqID string) *TestRequest {
	kv := testRequest.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(testReqID)
	return testRequest
}

func (TestRequest) New() messages.TestRequestBuilder {
	return makeTestRequest()
}

func (TestRequest) Parse(data []byte) (messages.TestRequestBuilder, error) {
	return ParseTestRequest(data)
}

func (testRequest *TestRequest) SetFieldTestReqID(testReqID string) messages.TestRequestBuilder {
	return testRequest.SetTestReqID(testReqID)
}
