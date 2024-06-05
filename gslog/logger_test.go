package gslog

import (
	"os"
	"testing"
	"time"
)

func TestLogLevelUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		args string
		want LogLevel
	}{
		{
			name: "trace low case",
			args: "trace",
			want: TraceLevel,
		},
		{
			name: "debug low case",
			args: "debug",
			want: DebugLevel,
		},
		{
			name: "info low case",
			args: "info",
			want: InfoLevel,
		},
		{
			name: "warn low case",
			args: "warn",
			want: WarnLevel,
		},
		{
			name: "error low case",
			args: "error",
			want: ErrorLevel,
		},
		{
			name: "critical low case",
			args: "critical",
			want: CriticalLevel,
		},
		{
			name: "trace up case",
			args: "TRACE",
			want: TraceLevel,
		},
		{
			name: "debug up case",
			args: "DEBUG",
			want: DebugLevel,
		},
		{
			name: "info up case",
			args: "INFO",
			want: InfoLevel,
		},
		{
			name: "warn up case",
			args: "WARN",
			want: WarnLevel,
		},
		{
			name: "error up case",
			args: "ERROR",
			want: ErrorLevel,
		},
		{
			name: "critical up case",
			args: "CRITICAL",
			want: CriticalLevel,
		},
		{
			name: "trace capital case",
			args: "Trace",
			want: TraceLevel,
		},
		{
			name: "debug capital case",
			args: "Debug",
			want: DebugLevel,
		},
		{
			name: "info capital case",
			args: "Info",
			want: InfoLevel,
		},
		{
			name: "warn capital case",
			args: "Warn",
			want: WarnLevel,
		},
		{
			name: "error capital case",
			args: "Error",
			want: ErrorLevel,
		},
		{
			name: "critical capital case",
			args: "Critical",
			want: CriticalLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if lv, err := ParseLogLevel(tt.args); (err != nil) || tt.want != lv {
				t.Errorf("ParseLogLevel error = %v, wantE = %v, get = %v", err, tt.want, lv)
			}
		})
	}
}

func TestStdLogger(t *testing.T) {
	logger := NewLogger(WithWriteSyncer("console", NewStdWriteSyncer()), WithLevelEnabler(TraceLevel))
	defer logger.Close()

	for i := 0; i < 20; i++ {
		logger.Debug("test %d", i)
		time.Sleep(time.Millisecond * 100)
	}
}

func TestFileLogger(t *testing.T) {
	file, err := os.OpenFile("./example.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return
	}
	logger := NewLogger(WithWriteSyncer("file", NewFileWriteSyncer(file)), WithLevelEnabler(TraceLevel))
	defer logger.Close()

	for i := 0; i < 20; i++ {
		logger.Debug("test %d", i)
		time.Sleep(time.Millisecond * 100)
	}
}
