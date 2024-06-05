package gslog

import (
	"fmt"
	"io"
	"log"
	"os"
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

type WriteFunc func(p []byte) (n int, err error)

type baseWriteSyncer struct {
	recChan   chan *LogEntry
	closeChan chan struct{}
	wg        sync.WaitGroup
	isClosed  atomic.Bool
	writeFunc WriteFunc
}

func newBaseWriteSyncer(writeFunc WriteFunc) *baseWriteSyncer {
	gs := &baseWriteSyncer{
		recChan:   make(chan *LogEntry, _defaultChanSize),
		closeChan: make(chan struct{}),
		writeFunc: writeFunc,
	}

	gs.wg.Add(1)
	go gs.run()

	return gs
}

func (gs *baseWriteSyncer) Write(p []byte) (n int, err error) {
	return gs.writeFunc(p)
}

func (gs *baseWriteSyncer) Sync(entry *LogEntry) {
	if gs.isClosed.Load() {
		return
	}
	gs.recChan <- entry
}

func (gs *baseWriteSyncer) Close() error {
	if gs.isClosed.CompareAndSwap(false, true) {
		gs.closeChan <- struct{}{}
		close(gs.recChan)
		close(gs.closeChan)
		gs.wg.Wait()
	}
	return nil
}

func (gs *baseWriteSyncer) run() {
	defer gs.wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover err: ", err)
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
			_, err := gs.Write(formatEntry)
			if err != nil {
				log.Println("write error: ", err)
				continue
			}
		case <-gs.closeChan:
			// 处理完所有日志条目
			for {
				select {
				case entry, ok := <-gs.recChan:
					if !ok {
						return
					}
					formatEntry := formatLogEntry(entry)
					_, err := gs.Write(formatEntry)
					if err != nil {
						log.Println("write error: ", err)
						continue
					}
				default:
					return
				}
			}
		}
	}
}

type StdWriteSyncer struct {
	*baseWriteSyncer
}

func NewStdWriteSyncer() *StdWriteSyncer {
	gs := &StdWriteSyncer{}

	gs.baseWriteSyncer = newBaseWriteSyncer(gs.Write)

	return gs
}

func (gs *StdWriteSyncer) Write(p []byte) (n int, err error) {
	fmt.Printf(string(p))
	return len(p), nil
}

type FileWriteSyncer struct {
	file *os.File
	*baseWriteSyncer
}

func NewFileWriteSyncer(file *os.File) *FileWriteSyncer {
	gs := &FileWriteSyncer{
		file: file,
	}

	gs.baseWriteSyncer = newBaseWriteSyncer(gs.Write)

	return gs
}

func (gs *FileWriteSyncer) Write(p []byte) (n int, err error) {
	return fmt.Fprintf(gs.file, "%s", string(p))
}

func (gs *FileWriteSyncer) Close() error {
	_ = gs.baseWriteSyncer.Close()
	_ = gs.file.Close()

	return nil
}
