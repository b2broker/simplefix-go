package session

import "time"

// todo constructor for acceptor and initiator
type LogonSettings struct {
	TargetCompID  string
	SenderCompID  string
	HeartBtInt    int
	EncryptMethod string
	Password      string
	Username      string
	LogonTimeout  time.Duration // todo
	HeartBtLimits *IntLimits
}
