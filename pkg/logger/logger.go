package logger

import (
	"log"
	"os"
)

var levelColors = map[int]string{
	DEBUG: "\033[36m", // Cyan
	INFO:  "\033[32m", // Green
	WARN:  "\033[33m", // Yellow
	ERROR: "\033[31m", // Red
	FATAL: "\033[35m", // Magenta
}

const resetColor = "\033[0m"

type Logger struct {
	level       int
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	unused      int
}

func NewLogger(level int) *Logger {
	return &Logger{
		level:       level,
		debugLogger: log.New(os.Stdout, levelColors[DEBUG]+"DEBUG: "+resetColor, log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, levelColors[INFO]+"INFO: "+resetColor, log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, levelColors[WARN]+"WARN: "+resetColor, log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, levelColors[ERROR]+"ERROR: "+resetColor, log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, levelColors[FATAL]+"FATAL: "+resetColor, log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Printf(message, args...)
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	if l.level <= INFO {
		l.infoLogger.Printf(message, args...)
	}
}

func (l *Logger) Warn(message string, args ...interface{}) {
	if l.level <= WARN {
		l.warnLogger.Printf(message, args...)
	}
}

func (l *Logger) Error(message string, args ...interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Printf(message, args...)
	}
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	if l.level <= FATAL {
		l.fatalLogger.Printf(message, args...)
		os.Exit(1)
	}
}
