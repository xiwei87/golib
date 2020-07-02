package http_server

import (
	"time"

	"github.com/xiwei87/golib/common"
)

type HttpRequest struct {
	HttpMethod  string
	RequestId   string
	StartTime   time.Time
	RemoteAddr  string
	RequestTime string
	UserAgent   string
}

type HttpResponse struct {
	Status  int
	ErrCode common.ErrorCodeType
}
