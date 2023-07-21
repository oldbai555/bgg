package _const

import "strings"

var (
	LogWithHint = strings.ToUpper("hint")

	HeaderLBTraceId   = strings.ToUpper("X-LB-Trace-Id")
	HeaderLBDeviceId  = strings.ToUpper("X-LB-Device-Id")
	HeaderLBSid       = strings.ToUpper("X-LB-Sid")
	HeaderLBUid       = strings.ToUpper("X-LB-Uid")
	HeaderLBApiMethod = strings.ToUpper("X-LB-Api-Method")
	HeaderLBAuthType  = strings.ToUpper("X-LB-Auth-Type")
	HeaderLBCallFrom  = strings.ToUpper("X-LB-Call-From")
)
