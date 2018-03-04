package logx

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
    "runtime"
)

var (
	defaultLogPath     string
	defaultLogFileName string
	defaultLogLevel    int32
	lastDate           int32
	log                *Log
)

const (
	_ = iota
	LEVEL_NOTICE
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

func init() {
	defaultLogPath = "./log"
	defaultLogFileName = "log"
	defaultLogLevel = LEVEL_NOTICE
}

type Log struct {
	logPath     string
	logFileName string
	logLevel    int32

	Writer io.Writer
}

func Logger(logPath, logFileName string, logLevel int32) *Log {
	if log == nil {
		return New(logPath, logFileName, logLevel)
	}
	return log
}

func DefaultLogger() *Log {
	if log == nil {
		return New(defaultLogPath, defaultLogFileName, defaultLogLevel)
	}
	return log
}

func New(logPath, logFileName string, logLevel int32) *Log {
	l := &Log{logPath: logPath, logFileName: logFileName, logLevel: logLevel}
	l.checkLog()
	return l
}

func (l *Log) Msg(logLevel int32, message string) *Log {
	if logLevel < l.logLevel {
		return l
	}
    l.checkLog()
	var msg string
	var levelStr string
	switch logLevel {
	case LEVEL_NOTICE:
		levelStr = "NOTICE"
	case LEVEL_INFO:
		levelStr = "INFO"
	case LEVEL_WARN:
		levelStr = "WARN"
	case LEVEL_ERROR:
		levelStr = "ERROR"
	case LEVEL_FATAL:
		levelStr = "FATAL"
	default:
		levelStr = "NOTICE"
	}

    pc,file,line,ok := runtime.Caller(1)
    var funcName = "Unkonwn"
    if ok {
        funcName = runtime.FuncForPC(pc).Name()
        file = path.Base(file)
    } else {
        file = "Unkonwn"
        line = 0
    }
    msg = fmt.Sprintf("----------------------\n%s [%s] [%s:%d:%s] %s", time.Now().Format("2006-01-02 15:04:05"), levelStr, file, line, funcName, message)
	fmt.Fprintln(l.Writer, msg)
	return l
}

func (l *Log) checkLog() {
	var err error
	var file *os.File
	curDate := int32(time.Now().Day())
	if curDate != lastDate || l.Writer == nil {
		_, err = os.Stat(l.logPath)
		if err != nil {
			os.Mkdir(l.logPath, 0755)
		}
		fullLogFileName := path.Join(l.logPath, fmt.Sprintf("%s.%s.log", l.logFileName, time.Now().Format("20060102")))
		_, err = os.Stat(fullLogFileName)
		if err != nil {
			file, _ = os.Create(fullLogFileName)
		} else {
			file, _ = os.OpenFile(fullLogFileName, os.O_WRONLY|os.O_APPEND, 0755)
		}
        if l.Writer != nil {
            l.Writer.(*os.File).Sync()
            l.Writer.(*os.File).Close()
        }
		l.Writer = file
	}
	lastDate = curDate
}
