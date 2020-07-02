package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	slog "github.com/go-eden/slf4go"
	"github.com/xiwei87/golib/common"
	"github.com/xiwei87/golib/utils"
)

var server *http.Server

type HttpServer struct {
	Dispatch *gin.Engine
	Server   *http.Server
	Request  *HttpRequest
	Response *HttpResponse
}

type HttpCfg struct {
	ListenPort       int  `yaml:"listen_port"`
	ReadTimeout      uint `yaml:"read_timeout"`
	ReadIdleTimeout  uint `yaml:"read_idle_timeout"`
	WriteTimeout     uint `yaml:"write_timeout"`
	WriteIdleTimeout uint `yaml:"write_idle_timeout"`
	MaxHeaderSize    int  `yaml:"max_header_size"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
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
		_, _ = ctx.Writer.Write([]byte("STATUS OK"))
	})
	return s
}

func (s *HttpServer) StartServer(cfg *HttpCfg) error {
	if cfg == nil {
		return errors.New("配置文件为空")
	}
	addr := ":" + strconv.Itoa(cfg.ListenPort)
	server = &http.Server{
		Addr:           addr,
		Handler:        s.Dispatch,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderSize,
	}
	return server.ListenAndServe()
}

func (s *HttpServer) StopServer() error {
	if nil == server {
		return errors.New("Http Server Not Run")
	}
	return server.Shutdown(context.TODO())
}

func (s *HttpServer) printAccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* init logid, start_time, request_id, method, remoteAddr */
		s.Request.startTime = time.Now()
		s.Request.httpMethod = c.Request.Method
		s.Request.requestId = utils.NewRequestId()
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
		c.Header("X-Request-Id", s.Request.requestId)
		// Process request
		c.Next()
		// Print Access Log
		requestTime := (time.Now().UnixNano() - s.Request.startTime.UnixNano()) / 1e6
		slog.Infof("errno[%d] ip[%s] logId[%s] uri[%s] cost[%d] status[%d] ua[%s] request done",
			s.Response.errCode, s.Request.remoteAddr,
			s.Request.requestId, c.Request.URL.Path, requestTime,
			s.Response.status, s.Request.userAgent)
	}
}

func (s *HttpServer) NoMethodHandler(ctx *gin.Context) {
	format := `{"request_id":"%s","err_code":%d,"err_msg":"%s"}`
	resStr := fmt.Sprintf(format, s.Request.requestId, common.CODE_NO_METHOD, "method not allowed")
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.Write([]byte(resStr))
}
