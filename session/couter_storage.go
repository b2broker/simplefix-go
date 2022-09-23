package session

type CounterStorage interface {
	GetNextSeqNum(pk string) (int, error)
	GetCurrSeqNum(pk string) (int, error)
	ResetSeqNum(pk string) error
}
