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
)

var mutex sync.Mutex

/* global logger */
var Logger log4go.Logger
var initialized bool = false

type LogConfig struct {
	LogPath  string `yaml:"log_path"`
	LogLevel string `yaml:"log_level"`
	LogSave  int    `yaml:"log_save"`
}

func Init(cfg *LogConfig, hasStdOut bool) error {
	var (
		err    error
		logDir string
	)
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return errors.New("日志已经初始化")
	}
	if nil == cfg {
		return errors.New("日志配置信息为空")
	}
	/* remove the last '/'  */
	logDir = strings.TrimSuffix(cfg.LogPath, "/")
	Logger, err = Create("service", cfg.LogLevel, logDir, hasStdOut, "D", cfg.LogSave)
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
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat("[%L] %D %T [%S] %M")
	logger.AddFilter("log", level, logWriter)

	return logger, nil
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
