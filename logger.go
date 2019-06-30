package simpleLogger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
)

const (
	LevelFatal         = "FATAL"
	LevelError         = "ERROR"
	LevelWarning       = "WARN"
	LevelInformational = "INFO"
	LevelDebug         = "DEBUG"
)

type Slogger struct {
	Writer io.Writer
	lock   sync.Mutex
	once    sync.Once
	Levels []string
	Level string
	filterLevels map[string]struct{}
}

type Logger interface {
	Init(config string) error
	Write(p []byte) (n int, err error)
	Flush()
	Close()
}

var adaptersManager = make(map[string]func() Logger)

func Register(name string, logger func() Logger) {
	if logger == nil {
		panic("SimpleLog: Register provide is nil")
	}
	if _, ok := adaptersManager[name]; ok {
		panic("SimpleLog: Register called twice for provider " + name)
	}
	adaptersManager[name] = logger
}


func NewSimpleLog() *Slogger {
	logger := new(Slogger)
	_: logger.SetOutPut(AdapterConsole,"")
    logger.Levels = []string{LevelDebug, LevelInformational, LevelWarning, LevelError, LevelFatal}
	logger.Level = LevelDebug
	return logger
}

func (sl *Slogger) SetOutPut(adapterName string, config string) error {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	newLogger, ok := adaptersManager[adapterName]
	if !ok {
		return fmt.Errorf("SimpleLogs: adaptername %s not registered", adapterName)
	}
	Logger := newLogger()
	err := Logger.Init(config)

	if err != nil {
		fmt.Fprintln(os.Stderr, "logs.BeeLogger.SetLogger: "+err.Error())
		return err
	}
	sl.Writer = Logger
	return nil

}

func (sl *Slogger) SetLevel(l string) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.Level = l
	sl.init()
}

func (sl *Slogger) Check(line []byte) bool {
	sl.once.Do(sl.init)
	var level string
	b := bytes.IndexByte(line, '[')

	if b >= 0 {
		a := bytes.IndexByte(line[b:], ']')
		if a >= 0 {
			level = string(line[b+1 : b+a])
		}
	}

	_, ok := sl.filterLevels[level]
	return ok
}

func (sl *Slogger) init() {
	filterLevels := make(map[string]struct{})
	for _, level := range sl.Levels {
		if level == sl.Level {
			break
		}
		filterLevels[level] = struct{}{}
	}
	sl.filterLevels = filterLevels
}

func (sl *Slogger) Write(p []byte) (n int, err error) {
	if sl.Check(p) {
		return len(p), nil
	}

	return sl.Writer.Write(p)
}