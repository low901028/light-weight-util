package capnslog

import (
	"os"
	"io"
	"syscall"
)

func init() {
	initHijack()

	// Go `log` pacakge uses os.Stderr.
	SetFormatter(NewDefaultFormatter(os.Stderr))
	SetGlobalLogLevel(INFO)
}

func NewDefaultFormatter(out io.Writer) Formatter {
	if syscall.Getppid() == 1 {
		f, err := NewJournaldFormatter()
		if err == nil {
			return f
		}
	}
	return NewPrettyFormatter(out, false)
}