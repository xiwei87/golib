package log

import (
	"fmt"
	"strings"

	"github.com/baidu/go-lib/log/log4go"
	slog "github.com/go-eden/slf4go"
)

type LogDriver struct {
}

func (d *LogDriver) Name() string {
	return "slf4go-log"
}

func (d *LogDriver) Print(l *slog.Log) {
	var msg string
	if l.Format != nil {
		msg = fmt.Sprintf(*l.Format, l.Args...)
	} else {
		msg = fmt.Sprint(l.Args...)
	}
	var source string
	source = fmt.Sprintf("%s %s:%d", l.Logger, l.Stack.Filename, l.Stack.Line)
	switch l.Level {
	case slog.TraceLevel:
		Logger.Log(log4go.TRACE, source, msg)
	case slog.DebugLevel:
		Logger.Log(log4go.DEBUG, source, msg)
	case slog.InfoLevel:
		Logger.Log(log4go.INFO, source, msg)
	case slog.WarnLevel:
		Logger.Log(log4go.WARNING, source, msg)
	case slog.ErrorLevel:
		Logger.Log(log4go.ERROR, source, msg)
	case slog.PanicLevel:
		Logger.Log(log4go.CRITICAL, source, msg)
	case slog.FataLevel:
		Logger.Log(log4go.CRITICAL, source, msg)
	}
}

func (d *LogDriver) GetLevel(logger string) (sl slog.Level) {
	return slog.TraceLevel
}

// stringToLevel converts level in string to log4go level
func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)

	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}
