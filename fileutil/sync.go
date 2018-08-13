package fileutil

import "os"

// Fsync是一个file.sync的包装类
// 在darwin上需要特殊实现
func Fsync(f *os.File) error {
	return f.Sync()
}

// Fdatasync是file.sync()的包装类
// 在linux平台上需要特殊的处理
func Fdatasync(f *os.File) error {
	return f.Sync()
}
