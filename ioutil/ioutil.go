package ioutil

import (
	"fileutil"
	"io"
	"os"
)

// WriteAndSyncFile与标准库中ioutil.WriteFile功能类似：WriteAndSyncFile在关闭文件前调用同步
// 在未出现error时，则完成data同步
func WriteAndSyncFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if err == nil {
		fileutil.Fsync(f)
	}

	if err1 := f.Close(); err != nil {
		err = err1
	}

	return err
}
