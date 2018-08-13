package fileutil

import "os"

func Fsync(f *os.File) error {
	return f.Sync()
}

// Fdatasync类似于fsync()，但不会刷新修改后的元数据
// 除非需要元数据才能正确处理后续数据检索
func Fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}
