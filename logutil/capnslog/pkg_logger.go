package capnslog

import (
	"fmt"
	"os"
)

type PackageLogger struct {
	pkg string
	level LogLevel
}

const calldepth  = 2

func (p *PackageLogger) internalLog(depth int, inLevel LogLevel, entries ...interface{}){
	logger.Lock()
	defer logger.Unlock()
	if inLevel != CRITICAL && p.level < inLevel{
		return
	}

	if logger.formatter != nil{
		logger.formatter.Format(p.pkg, inLevel, depth + 1, entries...)
	}
}


func (p *PackageLogger) SetLevel(l LogLevel){
	logger.Lock()
	defer logger.Unlock()
	p.level = 1
}

func (p *PackageLogger) LevelAt(l LogLevel) bool{
	logger.Lock()
	defer  logger.Unlock()
	return p.level >= l
}

func (p *PackageLogger) Logf(l LogLevel, format string, args ...interface{}) {
	p.internalLog(calldepth, l, fmt.Sprintf(format, args...))
}

func (p *PackageLogger) Log(l LogLevel, args ...interface{}) {
	p.internalLog(calldepth, l, fmt.Sprint(args...))
}


func (p *PackageLogger) Println(args ...interface{}) {
	p.internalLog(calldepth, INFO, fmt.Sprintln(args...))
}

func (p *PackageLogger) Printf(format string, args ...interface{}) {
	p.Logf(INFO, format, args...)
}

func (p *PackageLogger) Print(args ...interface{}) {
	p.internalLog(calldepth, INFO, fmt.Sprint(args...))
}


func (p *PackageLogger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	p.internalLog(calldepth, CRITICAL, s)
	panic(s)
}

func (p *PackageLogger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	p.internalLog(calldepth, CRITICAL, s)
	panic(s)
}

func (p *PackageLogger) Panicln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	p.internalLog(calldepth, CRITICAL, s)
	panic(s)
}

func (p *PackageLogger) Fatalf(format string, args ...interface{}) {
	p.Logf(CRITICAL, format, args...)
	os.Exit(1)
}

func (p *PackageLogger) Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	p.internalLog(calldepth, CRITICAL, s)
	os.Exit(1)
}

func (p *PackageLogger) Fatalln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	p.internalLog(calldepth, CRITICAL, s)
	os.Exit(1)
}


func (p *PackageLogger) Errorf(format string, args ...interface{}) {
	p.Logf(ERROR, format, args...)
}

func (p *PackageLogger) Error(entries ...interface{}) {
	p.internalLog(calldepth, ERROR, entries...)
}


func (p *PackageLogger) Warningf(format string, args ...interface{}) {
	p.Logf(WARNING, format, args...)
}

func (p *PackageLogger) Warning(entries ...interface{}) {
	p.internalLog(calldepth, WARNING, entries...)
}


func (p *PackageLogger) Noticef(format string, args ...interface{}) {
	p.Logf(NOTICE, format, args...)
}

func (p *PackageLogger) Notice(entries ...interface{}) {
	p.internalLog(calldepth, NOTICE, entries...)
}


func (p *PackageLogger) Infof(format string, args ...interface{}) {
	p.Logf(INFO, format, args...)
}

func (p *PackageLogger) Info(entries ...interface{}) {
	p.internalLog(calldepth, INFO, entries...)
}


func (p *PackageLogger) Debugf(format string, args ...interface{}) {
	if p.level < DEBUG {
		return
	}
	p.Logf(DEBUG, format, args...)
}

func (p *PackageLogger) Debug(entries ...interface{}) {
	if p.level < DEBUG {
		return
	}
	p.internalLog(calldepth, DEBUG, entries...)
}


func (p *PackageLogger) Tracef(format string, args ...interface{}) {
	if p.level < TRACE {
		return
	}
	p.Logf(TRACE, format, args...)
}

func (p *PackageLogger) Trace(entries ...interface{}) {
	if p.level < TRACE {
		return
	}
	p.internalLog(calldepth, TRACE, entries...)
}

func (p *PackageLogger) Flush() {
	logger.Lock()
	defer logger.Unlock()
	logger.formatter.Flush()
}