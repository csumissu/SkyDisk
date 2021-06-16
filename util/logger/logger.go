package logger

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"sync"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type logger struct {
	level int
	mutex sync.Mutex
}

var colors = map[int]*color.Color{
	LevelDebug: color.New(color.FgWhite).Add(color.Bold),
	LevelInfo:  color.New(color.FgHiWhite).Add(color.Bold),
	LevelWarn:  color.New(color.FgYellow).Add(color.Bold),
	LevelError: color.New(color.FgRed).Add(color.Bold),
	LevelFatal: color.New(color.BgRed).Add(color.Bold),
}

var prefixes = map[int]string{
	LevelDebug: "[DEBUG]",
	LevelInfo:  "[INFO]",
	LevelWarn:  "[WARN]",
	LevelError: "[ERROR]",
	LevelFatal: "[FATAL]",
}

var defaultLogger *logger

func InitLogger(level int) {
	defaultLogger = &logger{
		level: level,
	}
}

func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(LevelFatal, msg)
	panic(msg)
}

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(LevelError, msg)
}

func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(LevelWarn, msg)
}

func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(LevelInfo, msg)
}

func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	innerPrintln(LevelDebug, msg)
}

func innerPrintln(level int, msg string) {
	if level < defaultLogger.level {
		return
	}

	defaultLogger.mutex.Lock()
	defer defaultLogger.mutex.Unlock()

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
