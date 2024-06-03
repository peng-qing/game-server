package gslog

import "io"

type WriteSyncer interface {
	io.Writer

	Sync(entry *LogEntry)
}
