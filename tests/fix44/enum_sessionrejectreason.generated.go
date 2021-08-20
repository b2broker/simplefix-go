package fix44

import ()

// Enum type EnumSessionRejectReason
const (
	EnumSessionRejectReasonInvalidtagnumber                               string = "0"
	EnumSessionRejectReasonRequiredtagmissing                             string = "1"
	EnumSessionRejectReasonSendingtimeaccuracyproblem                     string = "10"
	EnumSessionRejectReasonInvalidmsgtype                                 string = "11"
	EnumSessionRejectReasonXmlvalidationerror                             string = "12"
	EnumSessionRejectReasonTagappearsmorethanonce                         string = "13"
	EnumSessionRejectReasonTagspecifiedoutofrequiredorder                 string = "14"
	EnumSessionRejectReasonRepeatinggroupfieldsoutoforder                 string = "15"
	EnumSessionRejectReasonIncorrectnumingroupcountforrepeatinggroup      string = "16"
	EnumSessionRejectReasonNondatavalueincludesfielddelimitersohcharacter string = "17"
	EnumSessionRejectReasonTagNotDefinedForThisMessageType                string = "2"
	EnumSessionRejectReasonUndefinedtag                                   string = "3"
	EnumSessionRejectReasonTagspecifiedwithoutavalue                      string = "4"
	EnumSessionRejectReasonValueisincorrectoutofrangeforthistag           string = "5"
	EnumSessionRejectReasonIncorrectdataformatforvalue                    string = "6"
	EnumSessionRejectReasonDecryptionproblem                              string = "7"
	EnumSessionRejectReasonSignatureproblem                               string = "8"
	EnumSessionRejectReasonCompidproblem                                  string = "9"
	EnumSessionRejectReasonOther                                          string = "99"
)
