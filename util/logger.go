package util

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type logger struct {
	level LogLevel
	mutex sync.Mutex
}

var colors = map[LogLevel]*color.Color{
	LevelDebug: color.New(color.FgWhite).Add(color.Bold),
	LevelInfo:  color.New(color.FgHiWhite).Add(color.Bold),
	LevelWarn:  color.New(color.FgYellow).Add(color.Bold),
	LevelError: color.New(color.FgRed).Add(color.Bold),
	LevelFatal: color.New(color.BgRed).Add(color.Bold),
}

var prefixes = map[LogLevel]string{
	LevelDebug: "[DEBUG]",
	LevelInfo:  "[INFO]",
	LevelWarn:  "[WARN]",
	LevelError: "[ERROR]",
	LevelFatal: "[FATAL]",
}

var Logger *logger

func InitLogger(level LogLevel) {
	Logger = &logger{
		level: level,
	}
}

func (logger *logger) Panic(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(logger, LevelFatal, msg)
	panic(msg)
}

func (logger *logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(logger, LevelError, msg)
}

func (logger *logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(logger, LevelWarn, msg)
}

func (logger *logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(logger, LevelInfo, msg)
}

func (logger *logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(logger, LevelDebug, msg)
}

func innerPrintln(logger *logger, level LogLevel, msg string) {
	if level < logger.level {
		return
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	file, line := getCallerFileAndLine()

	colorFunc := colors[level].SprintFunc()
	fmt.Printf("%s %s %s:%d - %s\n",
		colorFunc(prefixes[level]),
		time.Now().Format("2006/01/02 15:04:05"),
		file,
		line,
		msg)
}

func getCallerFileAndLine() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return file, line
}
