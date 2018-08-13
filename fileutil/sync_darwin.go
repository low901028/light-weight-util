package fileutil

import (
	"os"
	"syscall"
)

// Fsync在HFS/OSX上flush数据到物理驱动，但物理驱动在某些时候不会立马完成刷盘到物理存储介质，并且将以无序序列写入
// 使用F_FULLFSYNC来确保物理驱动buffer数据能够刷新到物理存储介质中
func Fsync(f *os.File) error {
	_, _, err := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), uintptr(syscall.F_FULLFSYNC), uintptr(0))
	if err != nil {
		return nil
	}
	return err
}

// Fdatasync在darwin平台上invoke F_FULLFYSNC来真实持久化到物理存储介质
func Fdatasync(f *os.File) error {
	return Fsync(f)
}
