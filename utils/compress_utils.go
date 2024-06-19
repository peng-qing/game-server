package utils

import (
	"compress/gzip"
	"io"
	"os"
)

// CompressFileByGzip 压缩文件为gzip
func CompressFileByGzip(src, dst string) (err error) {
	defer func() {
		if err != nil {
			// 如果中途出错 删除目标文件
			os.Remove(dst)
		}
	}()

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	// 调整输出文件所有权和当前文件一致
	if err := chown(dst, fileInfo); err != nil {
		return err
	}

	gzFile, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gz := gzip.NewWriter(gzFile)
	// 拷贝到gz写入器
	if _, err = io.Copy(gz, file); err != nil {
		return err
	}
	if err = gz.Flush(); err != nil {
		return err
	}
	if err = gz.Close(); err != nil {
		return err
	}
	if err = gzFile.Close(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	if err = os.Remove(src); err != nil {
		return err
	}

	return nil
}
