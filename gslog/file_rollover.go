package gslog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 检查实现 io.Closer 接口
var _ io.Closer = (*LogFileRollover)(nil)

const (
	// 压缩后缀
	compressSuffix = ".gz"
	// 默认单位 MB
	megaByte = 1024 * 1024
	// 默认文件分割大小 100MB
	defaultSize = megaByte * 100
	// 追加文件不存在时文件后缀
	defaultNotExistFileSuffix = "-gs_rollover.log"
	// 备份文件名格式化
	backupTimeFormat = "2006-01-02T15:04:05.000"
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
	size int64
	mu   sync.Mutex
}

/////////// constructs

func NewLogFileRollover(file *os.File, maxSize int, maxBackups int, maxAge int, compress bool) *LogFileRollover {
	return &LogFileRollover{
		file:       file,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
	}
}

/////////// Accessors

func (gs *LogFileRollover) FileName() string {
	if gs.file != nil {
		return gs.file.Name()
	}
	// 默认路径
	fileName := filepath.Base(os.Args[0]) + defaultNotExistFileSuffix
	return filepath.Join(os.TempDir(), fileName)
}

func (gs *LogFileRollover) MaxSize() int64 {
	if gs.maxSize != 0 {
		return int64(gs.maxSize * megaByte)
	}
	return defaultSize
}

func (gs *LogFileRollover) FilePath() string {
	return filepath.Dir(gs.FileName())
}

/////////// implements

func (gs *LogFileRollover) Write(p []byte) (n int, err error) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	rewriteSize := int64(len(p))
	if rewriteSize > gs.MaxSize() {
		return 0, fmt.Errorf("write too large (%d>%d)", rewriteSize, gs.MaxSize())
	}

	if gs.file == nil {
		// 追加文件不存在 尝试创建
		if err = gs.tryOpenOrCreateFile(rewriteSize); err != nil {
			return 0, err
		}
	}

	// 超过MaxSize
	if rewriteSize+gs.size > gs.MaxSize() {
		//TODO 轮转到新文件
	}

	n, err = gs.file.Write(p)
	gs.size += int64(n)

	return n, err
}

func (gs *LogFileRollover) Close() error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	var err error
	if gs.file != nil {
		err = gs.file.Close()
		gs.file = nil
	}

	return err
}

///////// internal
///////// 下面接口为内部接口，都不进行加锁，由调用内部接口得外部接口完成加锁和释放

func (gs *LogFileRollover) tryOpenOrCreateFile(rewriteSize int64) error {
	fileName := gs.FileName()
	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return gs.openNew()
		}
		return err
	}

	if info.Size()+rewriteSize >= gs.MaxSize() {
		// TODO 轮转到新文件

	}
	// 尝试追加
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 追加失败 创建
		return gs.openNew()
	}

	gs.file = file
	gs.size = info.Size()

	return nil
}

func (gs *LogFileRollover) openNew() error {
	err := os.MkdirAll(gs.FilePath(), 0755)
	if err != nil {
		return err
	}
	name := gs.FileName()
	mode := os.FileMode(0644)
	info, err := os.Stat(name)
	if err == nil {
		// 老文件mode
		mode = info.Mode()
		newName := backupName(gs.FileName())
		if err = os.Rename(name, newName); err != nil {
			return err
		}
	}

	newFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	gs.file = newFile
	gs.size = 0

	return nil
}

/////////// util functions

func backupName(name string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := strings.TrimSuffix(filename, ext)
	nowTm := time.Now()

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", prefix, nowTm.Format(backupTimeFormat), ext))
}
