package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var l *logrus.Logger

var stdErrFileHandler *os.File

type Log struct {
	UUID uuid.UUID
	log  *logrus.Logger
}

func init() {
	l = logrus.New()
	logFile, err := os.OpenFile("log/"+time.Now().Format("2006_01_02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("can not creat log file: %s", err.Error())
		return
	}
	stdErrFileHandler = logFile
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	setLog(io.MultiWriter(stdErrFileHandler, os.Stdout), GetLevel())
}

func NewLog() *Log {
	return &Log{
		UUID: uuid.New(),
		log:  l,
	}
}

func GetLevel() logrus.Level {
	var level logrus.Level
	switch AppConf.App.LogLevel {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "error":
		level = ErrorLevel
	case "warn":
		level = InfoLevel
	case "fatal":
		level = FatalLevel
	case "panic":
		level = PanicLevel
	default:
		level = InfoLevel
	}
	return level
}

//setLog set log output file and level
func setLog(fi io.Writer, level logrus.Level) {
	l.SetOutput(fi)
	l.SetLevel(level)
}

func (e *Log) buildMsg() string {
	_, file, line, _ := runtime.Caller(3)
	tmp := strings.Split(file, "/")
	size := len(tmp)
	if size > 1 {
		file = tmp[size-1]
	}
	locationInfo := fmt.Sprintf("file:%s line:%d", file, line)
	if e.UUID != [16]byte{00000000 - 0000 - 0000 - 0000 - 000000000000} {
		locationInfo = fmt.Sprintf("%s %s", e.UUID, locationInfo)
	}
	return locationInfo
}

func (e *Log) getParamsln(args ...interface{}) []interface{} {
	base := e.buildMsg()
	params := []interface{}{}
	params = append(params, base)
	return append(params, args...)
}

func (e *Log) getParamsf(format string) string {
	return e.buildMsg() + "[" + format + "]"
}

func (e *Log) Debugf(format string, args ...interface{}) {
	l.Debugf(e.getParamsf(format), args...)
}

func (e *Log) Printf(format string, args ...interface{}) {
	l.Printf(e.getParamsf(format), args...)
}

func (e *Log) Infof(format string, args ...interface{}) {
	l.Infof(e.getParamsf(format), args...)
}

func (e *Log) Warnf(format string, args ...interface{}) {
	l.Warnf(e.getParamsf(format), args...)
}

func (e *Log) Fatalf(format string, args ...interface{}) {
	l.Fatalf(e.getParamsf(format), args...)
}

func (e *Log) Debugln(args ...interface{}) {
	l.Debugln(e.getParamsln(args))
}

func (e *Log) Println(args ...interface{}) {
	l.Println(e.getParamsln(args))
}

func (e *Log) Infoln(args ...interface{}) {
	l.Infoln(e.getParamsln(args))
}

func (e *Log) Warnln(args ...interface{}) {
	l.Warnln(e.getParamsln(args))
}

func (e *Log) Errorln(args ...interface{}) {
	l.Errorln(e.getParamsln(args))
}

func (e *Log) Fatalln(args ...interface{}) {
	l.Fatalln(e.getParamsln(args))
}
