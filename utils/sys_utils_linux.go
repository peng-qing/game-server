//go:build linux
// +build linux

package utils

import "os"

func chown(name string, fileInfo os.FileInfo) error {
	// 确保文件存在 不存在创建
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return err
	}
	dstFile.Close()

	stat := fileInfo.Sys().(*syscall.Stat_t)
	return os.Chown(name, int(stat.Uid), int(stat.Gid))
}
