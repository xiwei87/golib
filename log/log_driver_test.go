package log

import (
	"testing"
	"time"

	slog "github.com/go-eden/slf4go"
)

func TestGlobalLogger(t *testing.T) {
	config.Cfg.Server.LogSave = 72
	config.Cfg.Server.LogLevel = "INFO"
	config.Cfg.Server.LogPath = "../test/run"

	_ = Init("test", true)
	slog.Info("global logger")
	slog.Warnf("global logger, warnning: %v", "surrender")
	time.Sleep(1000 * time.Millisecond)
}
