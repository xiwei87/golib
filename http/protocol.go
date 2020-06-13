package http

import (
	"golib"
	"time"
)

type HttpRequest struct {
	httpMethod  string
	requestId   uint32
	startTime   time.Time
	remoteAddr  string
	requestTime string
	userAgent   string
}

type HttpResponse struct {
	status  int
	errCode golib.ErrorCodeType
}
