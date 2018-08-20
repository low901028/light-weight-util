package logutil

import (
	"time"
	"light-weight-util/logutil/capnslog"
	"fmt"
	"sync"
)

var (
	defaultMergePeriod  = time.Second
	defaultTimrOutputScale = 10 * time.Millisecond

	outputInterval = time.Second
)

type line struct {
	level capnslog.LogLevel
	str string
}

func (l line) append(s string) line{
	return line{
		level: l.level,
		str: l.str + " " + s,
	}
}


// 日志合并line的状态
type status struct {
	period time.Duration

	start time.Time

	count int
}

func (s *status) isInMergePeriod(now time.Time) bool{
	return s.period == 0 || s.start.Add(s.period).After(now)
}

func (s *status) isEmpty() bool{return s.count == 0}

func (s *status) summary(now time.Time) string{
	ts := s.start.Round(defaultTimrOutputScale)
	took := now.Round(defaultTimrOutputScale).Sub(ts)

	return fmt.Sprintf("[merged %d repeated lines in %s]", s.count, took)
}

func (s *status) reset(now time.Time){
	s.start = now
	s.count = 0
}


// MergeLogger支持日志合并，合并重复行并打印日志行概要
type MergeLogger struct {
	*capnslog.PackageLogger

	mu sync.Mutex
	statusm map[line]*status
}

func NewMergeLogger(logger *capnslog.PackageLogger) *MergeLogger{
	l := &MergeLogger{
		PackageLogger: logger,
		statusm: make(map[line]*status),
	}
	go l.outputLoop()
	return l
}

func (l *MergeLogger) MergeInfo(entries ...interface{}){
	l.merge(line{
		level: capnslog.INFO,
		str: fmt.Sprint(entries...),
	})
}

func (l *MergeLogger) MergeInfof(format string, args ...interface{}){
	l.merge(line{
		level: capnslog.INFO,
		str: fmt.Sprintf(format, args...),
	})
}

func (l *MergeLogger) MergeNotice(entries ...interface{}){
	l.merge(line{
		level: capnslog.INFO,
		str: fmt.Sprint(entries...),
	})
}

func (l *MergeLogger) MergeNoticef(format string, args ...interface{}){
	l.merge(line{
		level: capnslog.INFO,
		str: fmt.Sprintf(format, args...),
	})
}


func (l *MergeLogger) MergeWarning(entries ...interface{}) {
	l.merge(line{
		level: capnslog.WARNING,
		str:   fmt.Sprint(entries...),
	})
}

func (l *MergeLogger) MergeWarningf(format string, args ...interface{}) {
	l.merge(line{
		level: capnslog.WARNING,
		str:   fmt.Sprintf(format, args...),
	})
}

func (l *MergeLogger) MergeError(entries ...interface{}) {
	l.merge(line{
		level: capnslog.ERROR,
		str:   fmt.Sprint(entries...),
	})
}

func (l *MergeLogger) MergeErrorf(format string, args ...interface{}) {
	l.merge(line{
		level: capnslog.ERROR,
		str:   fmt.Sprintf(format, args...),
	})
}

func (l *MergeLogger) merge(ln line){
	l.mu.Lock()
	if status, ok := l.statusm[ln]; ok{
		status.count++
		l.mu.Unlock()
		return
	}

	l.statusm[ln] = &status{
		period: defaultMergePeriod,
		start: time.Now(),
	}

	l.mu.Unlock()
	l.PackageLogger.Logf(ln.level, ln.str)
}

func (l *MergeLogger) outputLoop(){
	for now := range time.Tick(outputInterval){
		var outputs []line

		l.mu.Lock()
		for ln, status := range l.statusm{
			if status.isInMergePeriod(now){
				continue
			}
			if status.isEmpty(){
				delete(l.statusm, ln)
				continue
			}
			outputs = append(outputs, ln.append(status.summary(now)))
			status.reset(now)
		}
		l.mu.Unlock()

		for _,o := range outputs{
			l.PackageLogger.Logf(o.level, o.str)
		}
	}
}