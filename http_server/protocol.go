package http_server

import (
	"time"

	"github.com/xiwei87/golib/common"
)

type HttpRequest struct {
	httpMethod  string
	requestId   string
	startTime   time.Time
	remoteAddr  string
	requestTime string
	userAgent   string
}

type HttpResponse struct {
	status  int
	errCode common.ErrorCodeType
}
