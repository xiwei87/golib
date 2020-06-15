package http

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.66ifuel.com/golang-tools/golib/common"
	"gitlab.66ifuel.com/golang-tools/golib/config"
	"gitlab.66ifuel.com/golang-tools/golib/log"
)

type HttpServer struct {
	Dispatch *gin.Engine
	Server   *http.Server
	Request  *HttpRequest
	Response *HttpResponse
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func getNewLogid() uint32 {
	pid := os.Getpid()
	return uint32((pid&0xfff)<<20) | rand.Uint32()
}

func NewHttpServer() *HttpServer {
	s := &HttpServer{
		Dispatch: gin.New(),
		Request:  &HttpRequest{},
		Response: &HttpResponse{},
	}
	/* init route */
	s.Dispatch.Use(s.printAccessLog())
	s.Dispatch.NoMethod(s.NoMethodHandler)
	s.Dispatch.NoRoute(s.NoMethodHandler)
	/* for monitor */
	s.Dispatch.GET("/opmon", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("STATUS OK"))
	})
	return s
}

func (s *HttpServer) Run() error {
	addr := ":" + strconv.Itoa(config.Cfg.Http.ListenPort)
	s.Server = &http.Server{
		Addr:           addr,
		Handler:        s.Dispatch,
		ReadTimeout:    time.Duration(config.Cfg.Http.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Cfg.Http.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.Cfg.Http.MaxHeaderSize,
	}
	return s.Server.ListenAndServe()
}

func (s *HttpServer) Close() error {
	return s.Server.Close()
}

func (s *HttpServer) SetKeepAlivesEnabled(v bool) {
	s.Server.SetKeepAlivesEnabled(v)
}

func (s *HttpServer) printAccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init logid, start_time, request_id, method, remoteAddr */
		s.Request.startTime = time.Now()
		s.Request.httpMethod = c.Request.Method
		s.Request.requestId = getNewLogid()
		s.Request.remoteAddr = c.ClientIP()
		if userAgent := c.Request.Header.Get("User-Agent"); userAgent != "" {
			s.Request.userAgent = userAgent
		} else {
			s.Request.userAgent = "unknown"
		}
		/* init response */
		s.Response.status = http.StatusOK
		s.Response.errCode = common.CODE_OK
		c.Header("Server", "66ifuel/1.0.0")
		c.Header("X-Request-Id", fmt.Sprintf("%d", s.Request.requestId))
		// Process request
		c.Next()
		// Print Access Log
		requestTime := (time.Now().UnixNano() - s.Request.startTime.UnixNano()) / 1e6
		format := "errno[%d] ip[%s] logId[%d] uri[%s] cost[%d] status[%d] ua[%s] request done"
		log.Logger.Info(format, s.Response.errCode, s.Request.remoteAddr,
			s.Request.requestId, c.Request.URL.Path, requestTime,
			s.Response.status, s.Request.userAgent)
	}
}

func (s *HttpServer) NoMethodHandler(ctx *gin.Context) {
	format := `{"request_id":%d,"err_code":%d,"err_msg":"%s"}`
	resStr := fmt.Sprintf(format, s.Request.requestId, common.CODE_NO_METHOD, "method not allowed")
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.Write([]byte(resStr))
}
