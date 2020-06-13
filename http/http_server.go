package http

import (
	"fmt"
	"golib"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
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

func (s *HttpServer) run() error {
	addr := ":" + strconv.Itoa(golib.Cfg.Http.ListenPort)
	s.Server = &http.Server{
		Addr:           addr,
		Handler:        s.Dispatch,
		ReadTimeout:    time.Duration(golib.Cfg.Http.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(golib.Cfg.Http.WriteTimeout) * time.Second,
		MaxHeaderBytes: golib.Cfg.Http.MaxHeaderSize,
	}
	return s.Server.ListenAndServe()
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
		s.Response.errCode = golib.CODE_OK
		c.Header("Server", "66ifuel/1.0.0")
		c.Header("X-Request-Id", fmt.Sprintf("%d", s.Request.requestId))
		// Process request
		c.Next()
		// Print Access Log
		requestTime := (time.Now().UnixNano() - s.Request.startTime.UnixNano()) / 1e6
		log.Info("errno[%d] ip[%s] logId[%d] uri[%s] cost[%d] status[%d] ua[%s] request done",
			s.Response.errCode, s.Request.remoteAddr, s.Request.requestId, c.Request.URL.Path, requestTime,
			s.Response.status, s.Request.userAgent)
	}
}

func (s *HttpServer) NoMethodHandler(ctx *gin.Context) {
	format := `{"request_id":%d,"err_code":%d,"err_msg":"%s"}`
	resStr := fmt.Sprintf(format, s.Request.requestId, golib.CODE_NO_METHOD, "method not allowed")
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.Write([]byte(resStr))
}
