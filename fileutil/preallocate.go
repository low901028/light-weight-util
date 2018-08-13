package fileutil

import (
	"io"
	"os"
)

// Preallocate尝试给指定的文件分配空间
// 该操作仅支持Linux少数的文件系统: btrfs/ext4等
// 若是该操作不被支持，将没有任何error返回
// 否则，将返回碰到的错误
func Preallocate(f *os.File, sizeInBytes int64, extendFile bool) error {
	if sizeInBytes == 0 {
		return nil
	}
	if extendFile {
		return preallocExtend(f, sizeInBytes)
	}
	return preallocFixed(f, sizeInBytes)
}

//
func preallocExtendTrunc(f *os.File, sizeInBytes int64) error {
	curOff, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	size, err := f.Seek(sizeInBytes, io.SeekEnd)
	if err != nil {
		return err
	}

	if _, err = f.Seek(curOff, io.SeekStart); err != nil {
		return err
	}

	if sizeInBytes > size {
		return nil
	}
	return f.Truncate(sizeInBytes)
}
