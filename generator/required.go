package generator

var ExcludedFields = map[string]bool{
	"BeginString": true,
	"BodyLength":  true,
	"MsgType":     true,
	"CheckSum":    true,
}

var RequiredHeaderFields = map[string]bool{
	// Always unencrypted, must be first field in message
	"BeginString": true,

	// Always unencrypted, must be second field in message
	"BodyLength": true,

	// Always unencrypted, must be third field in message
	"MsgType": true,

	// Assigned value used to identify firm sending message.
	"SenderCompID": true,

	// Assigned value used to identify receiving firm.
	"TargetCompID": true,

	// Integer message sequence number.
	"MsgSeqNum": true,

	// Time of message transmission (always expressed in UTC (Universal Time Coordinated, also known as "GMT")
	"SendingTime": true,
}

var RequiredTrailerFields = map[string]bool{
	// Always unencrypted, always last field in message
	"CheckSum": true,
}

var DefaultFlowFields = map[string][]string{
	"Logon":         {"HeartBtInt", "EncryptMethod", "Password", "Username"},
	"Logout":        nil,
	"Heartbeat":     {"TestReqID"},
	"TestRequest":   {"TestReqID"},
	"ResendRequest": {"BeginSeqNo", "EndSeqNo"},
	"SequenceReset": {"NewSeqNo"},
	"Reject":        {"SessionRejectReason", "RefSeqNum", "RefTagID"},
}
