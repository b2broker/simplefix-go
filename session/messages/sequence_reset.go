package messages

type SequenceReset interface {
	New() SequenceResetBuilder
	Build() SequenceResetBuilder
	NewSeqNo() int
	SetFieldNewSeqNo(newSeqNo int) SequenceResetBuilder
}

// SequenceResetBuilder is an interface providing functionality to a builder of auto-generated SequenceReset messages.
type SequenceResetBuilder interface {
	SequenceReset
	PipelineBuilder
}
