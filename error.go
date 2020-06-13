package golib

type ErrorCodeType int32

const (
	CODE_OK             ErrorCodeType = 0
	CODE_INTERNAL_ERROR               = 1001
	CODE_NO_PERMISSION                = 1002
	CODE_NO_METHOD                    = 1003
)
