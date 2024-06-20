package gslog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"GameServer/common/stl"
	"GameServer/utils"
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
	backupTimeFormat = "2006-01-02T15-04-05.000"
)

// LogFileInfo 日志文件元数据
type LogFileInfo struct {
	timestamp time.Time
	fileInfo  os.FileInfo
}

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
	// 锁
	mu sync.Mutex
	// 负责切割文件协程控制
	once      sync.Once       // 确保只创建一个
	ctx       context.Context // 控制协程退出
	ctxCancel context.CancelFunc
	// 通知子协程进行切割&压缩信号
	millChan chan struct{}
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

// Rotate 轮转日志文件
func (gs *LogFileRollover) Rotate() error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	return gs.rotate()
}

/////////// implements

func (gs *LogFileRollover) Write(p []byte) (n int, err error) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	rewriteSize := int64(len(p))
	if rewriteSize > gs.maxFileSize() {
		return 0, fmt.Errorf("write too large (%d>%d)", rewriteSize, gs.maxFileSize())
	}

	if gs.file == nil {
		// 追加文件不存在 尝试创建
		if err = gs.tryOpenOrCreateFile(rewriteSize); err != nil {
			return 0, err
		}
	}

	// 超过MaxSize
	if rewriteSize+gs.size > gs.maxFileSize() {
		// 轮转到新文件
		if err = gs.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = gs.file.Write(p)
	gs.size += int64(n)

	return n, err
}

// Close 调用对外的Close意味着关闭了这个io
// 内部信号channel和子协程一起会被关闭和退出
func (gs *LogFileRollover) Close() error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	return gs.close()
}

///////// internal
///////// 下面接口为内部接口，都不进行加锁，由调用内部接口的对外接口完成加锁和释放

func (gs *LogFileRollover) tryOpenOrCreateFile(rewriteSize int64) error {
	fileName := gs.fileName()
	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return gs.openNew()
		}
		return err
	}

	if info.Size()+rewriteSize >= gs.maxFileSize() {
		// 轮转到新文件
		return gs.rotate()
	}
	// 尝试追加
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 追加失败 创建
		return gs.openNew()
	}

	gs.file = file
	gs.size = info.Size()

	return gs.mill()
}

func (gs *LogFileRollover) openNew() error {
	err := os.MkdirAll(gs.filePath(), 0755)
	if err != nil {
		return err
	}
	name := gs.fileName()
	mode := os.FileMode(0644)
	info, err := os.Stat(name)
	if err == nil {
		// 老文件mode
		mode = info.Mode()
		newName := backupName(gs.fileName())
		if err = os.Rename(name, newName); err != nil {
			return err
		}
	}
	// 截断
	newFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	gs.file = newFile
	gs.size = 0

	return nil
}

func (gs *LogFileRollover) rotate() error {
	var err error
	// 先关闭旧的
	if err = gs.closeFile(); err != nil {
		return err
	}
	// 再开新的
	if err = gs.openNew(); err != nil {
		return err
	}

	return gs.mill()
}

func (gs *LogFileRollover) mill() error {
	gs.once.Do(func() {
		gs.millChan = make(chan struct{}, 1)
		ctx, cancel := context.WithCancel(context.Background())
		gs.ctx = ctx
		gs.ctxCancel = cancel
		go gs.millRun()
	})

	select {
	case gs.millChan <- struct{}{}:
	default:
	}

	return nil
}

func (gs *LogFileRollover) millRun() {
	for {
		select {
		case <-gs.ctx.Done():
			return
		case <-gs.millChan:
			_ = gs.millExec()
		}
	}
}

func (gs *LogFileRollover) millExec() error {
	if gs.maxAge == 0 && gs.maxBackups == 0 && !gs.compress {
		return nil
	}
	var err error
	// 加载出日志文件列表
	logFiles, err := gs.loadFileList()
	if err != nil {
		return errors.Join(err)
	}

	removed := make([]*LogFileInfo, 0)
	// 移除超出备份数量的日志
	if gs.maxBackups > 0 && len(logFiles) >= gs.maxBackups {
		savedSet := stl.NewSet[string]()
		remaining := make([]*LogFileInfo, 0)
		for _, logFile := range logFiles {
			filename := logFile.fileInfo.Name()
			if strings.HasSuffix(filename, compressSuffix) {
				filename = strings.TrimSuffix(filename, compressSuffix)
			}
			// 标记为保留
			savedSet.Insert(filename)
			if savedSet.Size() > gs.maxBackups {
				removed = append(removed, logFile)
				continue
			}
			remaining = append(remaining, logFile)
		}
		// 剩余文件
		logFiles = remaining
	}
	// 移除到期
	if gs.maxAge > 0 {
		remaining := make([]*LogFileInfo, 0)
		// 截至时间
		diff := time.Duration(int64(24*gs.maxAge) * int64(time.Hour))
		cutOffTm := time.Now().Add(-1 * diff)
		for _, logFile := range logFiles {
			// 早于截至时间删除
			if logFile.timestamp.Before(cutOffTm) {
				removed = append(removed, logFile)
				continue
			}
			remaining = append(remaining, logFile)
		}
		logFiles = remaining
	}
	// 删除
	for _, logFile := range removed {
		errRemove := os.Remove(filepath.Join(gs.filePath(), logFile.fileInfo.Name()))
		if errRemove != nil {
			err = errors.Join(err, errRemove)
			continue
		}
	}
	// 压缩
	if gs.compress {
		for _, logFile := range logFiles {
			if strings.HasSuffix(logFile.fileInfo.Name(), compressSuffix) {
				continue
			}
			// 压缩文件
			filename := filepath.Join(gs.filePath(), logFile.fileInfo.Name())
			errCompress := utils.CompressFileByGzip(filename, filename+compressSuffix)
			if errCompress != nil {
				err = errors.Join(err, errCompress)
				continue
			}
		}
	}

	return err
}

func (gs *LogFileRollover) closeFile() (err error) {
	if gs.file != nil {
		err = gs.file.Close()
	}

	return err
}

func (gs *LogFileRollover) close() error {
	err := gs.closeFile()
	if err != nil {
		return err
	}

	if gs.ctxCancel != nil {
		gs.ctxCancel()
		close(gs.millChan)
	}

	return nil
}

func (gs *LogFileRollover) fileName() string {
	if gs.file != nil {
		return gs.file.Name()
	}
	// 默认路径
	fileName := filepath.Base(os.Args[0]) + defaultNotExistFileSuffix
	return filepath.Join(os.TempDir(), fileName)
}

func (gs *LogFileRollover) maxFileSize() int64 {
	if gs.maxSize != 0 {
		return int64(gs.maxSize * megaByte)
	}
	return defaultSize
}

func (gs *LogFileRollover) filePath() string {
	return filepath.Dir(gs.fileName())
}

func (gs *LogFileRollover) loadFileList() ([]*LogFileInfo, error) {
	entries, err := os.ReadDir(gs.filePath())
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(gs.fileName())
	ext := filepath.Ext(filename)
	prefix := strings.TrimSuffix(filename, ext) + "_"

	logFiles := make([]*LogFileInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}
		// 解析日志文件名
		if ts, err := timeFromFileName(fileInfo.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, &LogFileInfo{
				timestamp: ts,
				fileInfo:  fileInfo,
			})
			continue
		}
		if ts, err := timeFromFileName(entry.Name(), prefix, ext+compressSuffix); err == nil {
			logFiles = append(logFiles, &LogFileInfo{
				timestamp: ts,
				fileInfo:  fileInfo,
			})
			continue
		}
		// 不满足设置的文件格式 当其他文件不进行处理
	}
	// 排序
	sort.Sort(sortableLogFiles(logFiles))

	return logFiles, nil
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

func timeFromFileName(name, prefix, ext string) (time.Time, error) {
	// 前缀
	if !strings.HasPrefix(name, prefix) {
		return time.Time{}, errors.New("invalid file name not prefix")
	}
	// 后缀
	if !strings.HasSuffix(name, ext) {
		return time.Time{}, errors.New("invalid file name not suffix")
	}

	ts := name[len(prefix) : len(name)-len(ext)]
	return time.Parse(backupTimeFormat, ts)
}

// 对LogFileInfos 进行排序 按照时间距离当前时间点 近->远
type sortableLogFiles []*LogFileInfo

func (gs sortableLogFiles) Len() int {
	return len(gs)
}

func (gs sortableLogFiles) Less(i, j int) bool {
	return gs[i].timestamp.After(gs[j].timestamp)
}

func (gs sortableLogFiles) Swap(i, j int) {
	gs[i], gs[j] = gs[j], gs[i]
}
