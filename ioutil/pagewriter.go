package ioutil

import (
	"io"
)

var defaultBufferBytes = 128 * 1024

// PageWriter实现io.writer
type PageWriter struct {
	w                 io.Writer
	pageOffset        int    // 跟踪buffer中page offset
	pageBytes         int    // 每页byte数
	bufferedBytes     int    // buffer中pending的byte数
	buf               []byte // 写buffer
	bufWatermarkBytes int    // 在执行flush之前buffer中byte的数量 小于len(buf),需要slack有空间来完成write操作，并能够保证页面对齐
}

// 创建PageWriter
func NewPageWriter(w io.Writer, pageBytes, pageOffset int) *PageWriter {
	return &PageWriter{
		w:                 w,
		pageOffset:        pageOffset,
		pageBytes:         pageBytes,
		buf:               make([]byte, defaultBufferBytes+pageBytes),
		bufWatermarkBytes: defaultBufferBytes,
	}
}

// page写操作
func (pw *PageWriter) Write(p []byte) (n int, err error) {
	// 当新write的 + 原有buffer总的bytes不大于默认buffer的bytes
	// 首先执行copy，将新write的bytes填到pagewriter的buffer中，置于原有的buffer byte之后
	// 接着变更buffer中的bytes  最后返回本次write的内容和error(未出错返回nil)
	if len(p)+pw.bufferedBytes <= pw.bufWatermarkBytes {
		copy(pw.buf[pw.bufferedBytes:], p)
		pw.bufferedBytes += len(p)
		return len(p), nil
	}

	// 对齐page
	// 根据pagewriter的page bytes减去当前page writer的offset和buffer中已有的byte 剩下的空间即为slack page
	slack := pw.pageBytes - ((pw.pageOffset + pw.bufferedBytes) % pw.pageBytes)
	// 当得到slack空间不等于page大小
	if slack != pw.pageBytes {
		partial := slack > len(p)
		if partial { // 未有足够的数据填充slack page
			slack = len(p)
		}
		copy(pw.buf[pw.bufferedBytes:], p[:slack]) // 将p中byte写入buffer
		pw.bufferedBytes += slack
		n = slack
		p = p[slack:] // 获取p中剩下的byte
		if partial {
			return n, nil
		}
	}
	if err = pw.Flush(); err != nil { // 页面对齐， 刷盘；clean buffer
		return n, nil
	}

	if len(p) > pw.pageBytes { // 直接写
		pages := len(p) / pw.pageBytes
		c, werr := pw.w.Write(p[:pages*pw.pageBytes])
		n += c
		if werr != nil {
			return n, nil
		}
		p = p[pages*pw.pageBytes:]
	}
	c, werr := pw.Write(p)
	n += c
	return n, werr
}

// buffer刷盘
func (pw *PageWriter) Flush() error {
	if pw.bufferedBytes == 0 {
		return nil
	}
	_, err := pw.w.Write(pw.buf[:pw.bufferedBytes])
	pw.pageOffset = (pw.pageOffset + pw.bufferedBytes) % pw.pageBytes // 当前page的offset
	pw.bufferedBytes = 0  // 执行flush之后  clean buffer
	return err
}
