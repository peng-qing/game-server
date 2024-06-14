package gslog

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	// 压缩后缀
	compressSuffix = ".gz"
	// 默认文件分割大小 100MB
	defaultSize = 1024 * 1024 * 100
	// 追加文件不存在时文件后缀
	defaultNotExistFileSuffix = "-gs_rollover.log"
)

// LogFileRollover 文件日志分割器
// 对File进行包装 实现io.Writer便于集成到log尽量对logger无感
type LogFileRollover struct {
	// 源文件 通过日志分割器的日志会被追加到该文件
	// 如果初始长度超出 MaxSize 会被切割并重命名加上当前时间信息
	// 然后会使用原始文件名创建一个新日志文件
	file *os.File
	// 单位MB
	maxSize int
	// 旧日志保留数量
	maxBackups int
	// 日志保存时间 天(24h)
	maxAge int
	// 是否执行压缩
	// 压缩后会被添加 .gz 后缀
	compress bool

	// 当前文件大小
	size int
	mu   sync.Mutex
}

func NewLogFileRollover(file *os.File, maxSize int, maxBackups int, maxAge int, compress bool) *LogFileRollover {
	return &LogFileRollover{
		file:       file,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

func (gs *LogFileRollover) FileName() string {
	if gs.file != nil {
		return gs.file.Name()
	}
	// 默认路径
	fileName := filepath.Base(os.Args[0]) + defaultNotExistFileSuffix
	return filepath.Join(os.TempDir(), fileName)
}

func (gs *LogFileRollover) MaxSize() int64 {

	return defaultSize
}

func (gs *LogFileRollover) Write(p []byte) (n int, err error) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	//rewriteSize := len(p)
	//if rewriteSize > gs.MaxSize() {
	//	// 写入大小溢出
	//	return 0, fmt.Errorf("write too large (%d > %d)", rewriteSize, gs.MaxSize)
	//}
	//
	//if gs.file == nil {
	//	// 如果追加文件不存在 在当前目录创建尝试创建临时目录
	//	if err = gs.tryOpenOrCreateFile(rewriteSize); err != nil {
	//		return 0, err
	//	}
	//}
	//
	//// 超过 MaxSize
	//if rewriteSize+gs.size > gs.MaxSize() {
	//
	//}

	//
	//	if l.size+writeLen > l.max() {
	//		if err := l.rotate(); err != nil {
	//			return 0, err
	//		}
	//	}
	//
	//	n, err = l.file.Write(p)
	//	l.size += int64(n)
	//
	return n, err
}

///////// internal

func (gs *LogFileRollover) tryOpenOrCreateFile(rewriteSize int) error {
	fileName := gs.FileName()
	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// todo open file
			return nil
		}
		return err
	}

	if info.Size()+int64(rewriteSize) >= gs.MaxSize() {

	}

	return nil
}
