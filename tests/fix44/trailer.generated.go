package fix44

import (
	"github.com/b2broker/simplefix-go/fix"
	"github.com/b2broker/simplefix-go/session/messages"
)

type Trailer struct {
	*fix.Component
}

func makeTrailer() *Trailer {
	return &Trailer{fix.NewComponent(
		fix.NewKeyValue(FieldSignatureLength, &fix.Int{}),
		fix.NewKeyValue(FieldSignature, &fix.String{}),
	)}
}

func NewTrailer() *Trailer {
	return makeTrailer()
}

func (trailer *Trailer) SignatureLength() int {
	kv := trailer.Get(0)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(int)
}

func (trailer *Trailer) SetSignatureLength(signatureLength int) *Trailer {
	kv := trailer.Get(0).(*fix.KeyValue)
	_ = kv.Load().Set(signatureLength)
	return trailer
}

func (trailer *Trailer) Signature() string {
	kv := trailer.Get(1)
	v := kv.(*fix.KeyValue).Load().Value()
	return v.(string)
}

func (trailer *Trailer) SetSignature(signature string) *Trailer {
	kv := trailer.Get(1).(*fix.KeyValue)
	_ = kv.Load().Set(signature)
	return trailer
}

func (Trailer) New() messages.TrailerBuilder {
	return makeTrailer()
}
