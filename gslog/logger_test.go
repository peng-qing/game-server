package gslog

import (
	"fmt"
	"testing"
)

type std struct{}

func (std) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func TestLogger(t *testing.T) {
	logger := NewLogger(NewTextHandler(std{}, WithLevelEnabler(TraceLevel), WithTextFlag(DefaultBitFlag)))

	logger.Debug("11111", "test", 1, 3333, 111, []int{11111, 23123123})
	logger.Info("22222", "test", 1, "info", 2, "list", []uint64{999, 888, 7777})
	logger.Warn("33333", "test", 1, Int("logLevel", 3))
	logger.Error("44444", "test", 1, Fields("LevelEnabler", Sting("LogLevelStr", "Error"), Int("LogLevel", ErrorLevel)))
	logger.Critical("55555", "test", 2, Fields("LevelEnabler", Sting("LogLevel", "Critical")))

	logger.TraceFields("debug fields", Int("LogLevel", 1), Fields("Fields", Int("Level", 10), Float("Amount", 100.00)))

	loggerJson := NewLogger(NewJsonHandler(std{}, WithLevelEnabler(TraceLevel)))

	loggerJson.Debug("Debug json message", "test", 1111)
}
