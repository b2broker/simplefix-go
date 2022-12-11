package generator

// ExcludedFields specifies the tags that will be omitted in generated messages
// because they are already included into the base message structure.
var ExcludedFields = map[string]bool{
	"BeginString": true,
	"BodyLength":  true,
	"MsgType":     true,
	"CheckSum":    true,
}

// RequiredHeaderFields indicates the required fields that must be contained in the header.
// A FIX message is not considered properly structured unless it contains these fields in its header.
var RequiredHeaderFields = map[string]bool{
	// Always unencrypted, must be the first field in a message.
	"BeginString": true,

	// Always unencrypted, must be the second field in a message.
	"BodyLength": true,

	// Always unencrypted, must be the third field in a message.
	"MsgType": true,

	// The assigned value is used to identify the party sending a message.
	"SenderCompID": true,

	// The assigned value is used to identify the party receiving a message.
	"TargetCompID": true,

	// An integer value, indicating the message sequence number.
	"MsgSeqNum": true,

	// The date and time of message transmission, in UTC time.
	"SendingTime": true,
}

// RequiredTrailerFields indicates the required field(s) that must be contained in the trailer.
// A FIX message is not considered properly structured unless it contains these fields in its trailer.
var RequiredTrailerFields = map[string]bool{
	// Always unencrypted, must be the last field in a message.
	"CheckSum": true,
}

// DefaultFlowFields indicates the required tags for each message type
// that must be contained in the trailer.
// A FIX session pipeline will not operate properly if any of these tags are missing for the specified messages.
var DefaultFlowFields = map[string][]string{
	"Logon":         {"HeartBtInt", "EncryptMethod", "Password", "Username"},
	"Logout":        nil,
	"Heartbeat":     {"TestReqID"},
	"TestRequest":   {"TestReqID"},
	"ResendRequest": {"BeginSeqNo", "EndSeqNo"},
	"SequenceReset": {"NewSeqNo"},
	"Reject":        {"SessionRejectReason", "RefSeqNum", "RefTagID"},
}
