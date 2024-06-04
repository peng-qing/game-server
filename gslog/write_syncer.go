package gslog

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type WriteSyncer interface {
	io.Writer

	Sync(entry *LogEntry)
	Close() error
}

const (
	_defaultChanSize = 128
)

type StdWriteSyncer struct {
	recChan   chan *LogEntry
	closeChan chan struct{}
	wg        sync.WaitGroup
	isClosed  atomic.Bool
}

func NewStdWriteSyncer() *StdWriteSyncer {
	gs := &StdWriteSyncer{
		recChan:   make(chan *LogEntry, _defaultChanSize),
		closeChan: make(chan struct{}),
	}

	gs.wg.Add(1)
	go gs.run()

	return gs
}

func (gs *StdWriteSyncer) Write(p []byte) (n int, err error) {

}

func (gs *StdWriteSyncer) Sync(entry *LogEntry) {
	if gs.isClosed.Load() {
		return
	}
	gs.recChan <- entry
}

func (gs *StdWriteSyncer) Close() error {
	if gs.isClosed.CompareAndSwap(false, true) {
		gs.closeChan <- struct{}{}
		close(gs.closeChan)
		gs.wg.Wait()
		close(gs.recChan)
	}
	return nil
}

func (gs *StdWriteSyncer) run() {
	defer gs.wg.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover err: ", err)
			return
		}
	}()
	for {
		select {
		case entry, ok := <-gs.recChan:
			if !ok {
				return
			}
			formatEntry := formatLogEntry(entry)
			gs.Write(formatEntry)
		case <-gs.closeChan:

		}
	}
}
