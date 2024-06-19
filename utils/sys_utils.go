//go:build !linux

package utils

import "os"

func chown(name string, fileInfo os.FileInfo) error {
	return nil
}
