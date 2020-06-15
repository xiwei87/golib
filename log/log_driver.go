package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/baidu/go-lib/log/log4go"
	slog "github.com/go-eden/slf4go"
	"gitlab.66ifuel.com/golang-tools/golib/config"
)

/* global logger */
var Logger log4go.Logger
var initialized bool = false
var mutex sync.Mutex

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

// logDirCreate checks and creates dir if nonexist
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// filenameGen generates filename
func filenameGen(progName, logDir string, isErrLog bool) string {
	/* remove the last '/'  */
	strings.TrimSuffix(logDir, "/")

	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".wf.log")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
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

// Init initializes log lib
//
// PARAMS:
//   - progName: program name. Name of log file will be progName.log
//   - hasStdOut: whether to have stdout output
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func Init(progName string, hasStdOut bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return errors.New("Initialized Already")
	}

	var logDir string
	/* remove the last '/'  */
	strings.TrimSuffix(config.Cfg.Server.LogPath, "/")
	logDir = config.Cfg.Server.LogPath + "/run/" + logDir

	var err error
	Logger, err = Create(progName, config.Cfg.Server.LogLevel, logDir, hasStdOut, "D", config.Cfg.Server.LogSave)
	if err != nil {
		return err
	}
	/* set log buffer size */
	log4go.SetLogBufferLength(10240)
	/* if blocking, log will be dropped */
	log4go.SetLogWithBlocking(false)
	log4go.SetLogFormat("[%L] %D %t [%S] %M")
	/* registration log */
	slog.SetDriver(&LogDriver{})

	initialized = true
	return nil
}

// Create creates log lib
//
// PARAMS:
//   - progName: program name. Name of log file will be progName.log
//   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
//   - logDir: directory for log. It will be created if noexist
//   - hasStdOut: whether to have stdout output
//   - when:
//       "M", minute
//       "H", hour
//       "D", day
//       "MIDNIGHT", roll over at midnight
//   - backupCount: If backupCount is > 0, when rollover is done, no more than
//       backupCount files are kept - the oldest ones are deleted.
func Create(progName string, levelStr string, logDir string, hasStdOut bool,
	when string, backupCount int) (log4go.Logger, error) {
	/* check when */
	if !log4go.WhenIsValid(when) {
		return nil, fmt.Errorf("invalid value of when: %s", when)
	}
	/* check, and create dir if nonexist    */
	if err := logDirCreate(logDir); err != nil {
		_ = log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return nil, err
	}
	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)
	/* create logger    */
	logger := make(log4go.Logger)
	/* create writer for stdout */
	if hasStdOut {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}
	/* create file writer for all log   */
	fileName := filenameGen("service", logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat("[%L] %D %t [%S] %M")
	logger.AddFilter("log", level, logWriter)

	return logger, nil
}
