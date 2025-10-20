package logger

import (
	"fmt"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type Logger struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	warningLogger *log.Logger
	fatalLogger   *log.Logger
	logFile       *os.File
	mu            sync.Mutex
	level        LogLevel
	timeFormat   string
	withCaller   bool
	jsonMode     bool
}

func NewLogger(logFilePath string) (*Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// auto color: disable if NO_COLOR or not a TTY; enable if FORCE_COLOR
	if os.Getenv("FORCE_COLOR") != "" {
		color.NoColor = false
	} else if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	return &Logger{
		infoLogger:    log.New(os.Stdout, color.CyanString("[INFO] "), log.Ltime),
		errorLogger:   log.New(os.Stdout, color.RedString("[ERROR] "), log.Ltime),
		debugLogger:   log.New(os.Stdout, color.GreenString("[DEBUG] "), log.Ltime),
		warningLogger: log.New(os.Stdout, color.YellowString("[WARNING] "), log.Ltime),
		fatalLogger:   log.New(os.Stdout, color.MagentaString("[FATAL] "), log.Ltime),
		logFile:       logFile,
		level:        DEBUG,
		timeFormat:   "2006-01-02 15:04:05",
	}, nil
}

func (l *Logger) log(level LogLevel, msg string) {
	// level guard
	if level < l.level {
		return
	}
	// Вывод в консоль
	switch level {
	case DEBUG:
		l.debugLogger.Println(msg)
	case INFO:
		l.infoLogger.Println(msg)
	case WARNING:
		l.warningLogger.Println(msg)
	case ERROR:
		l.errorLogger.Println(msg)
	case FATAL:
		l.fatalLogger.Println(msg)
	}

	// Запись в файл
	if l.logFile == nil {
		return
	}
	var levelStr string
	switch level {
	case DEBUG:
		levelStr = "[DEBUG] "
	case INFO:
		levelStr = "[INFO] "
	case WARNING:
		levelStr = "[WARNING] "
	case ERROR:
		levelStr = "[ERROR] "
	case FATAL:
		levelStr = "[FATAL] "
	}

	// caller info
	var caller string
	if l.withCaller {
		if _, file, line, ok := runtime.Caller(2); ok {
			caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	nowStr := time.Now().Format(l.timeFormat)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.jsonMode {
		entry := map[string]any{
			"time":  nowStr,
			"level": levelStr[1:len(levelStr)-2], // remove brackets and space
			"msg":   msg,
		}
		if caller != "" {
			entry["caller"] = caller
		}
		enc := json.NewEncoder(l.logFile)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(entry)
		return
	}

	line := nowStr + " " + levelStr + msg
	if caller != "" {
		line += " (" + caller + ")"
	}
	_, _ = fmt.Fprintln(l.logFile, line)
}

func (l *Logger) Debug(msg string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	fullMsg := msg
	if len(v) > 0 {
		fullMsg = fmt.Sprintf(msg, v...)
	}

	l.log(DEBUG, fullMsg)
}

func (l *Logger) Info(msg string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	fullMsg := msg
	if len(v) > 0 {
		fullMsg = fmt.Sprintf(msg, v...)
	}

	l.log(INFO, fullMsg)
}

func (l *Logger) Warning(msg string, v ...interface{}) {
	if l.level > WARNING {
		return
	}
	fullMsg := msg
	if len(v) > 0 {
		fullMsg = fmt.Sprintf(msg, v...)
	}

	l.log(WARNING, fullMsg)
}

func (l *Logger) Error(msg string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	fullMsg := msg
	if len(v) > 0 {
		fullMsg = fmt.Sprintf(msg, v...)
	}
	l.log(ERROR, fullMsg)
}

func (l *Logger) Fatal(msg string, v ...interface{}) {
	if l.level > FATAL {
		return
	}
	fullMsg := msg
	if len(v) > 0 {
		fullMsg = fmt.Sprintf(msg, v...)
	}

	l.log(FATAL, fullMsg)
	os.Exit(1)
}

func (l *Logger) Close() error {
	if l.logFile == nil {
		return nil
	}
	return l.logFile.Close()
}

// Configuration helpers (kept dead-simple)

func (l *Logger) SetLevel(level LogLevel) { l.level = level }
func (l *Logger) WithColors(enable bool)  { color.NoColor = !enable }
func (l *Logger) WithTimeFormat(layout string) { if layout != "" { l.timeFormat = layout } }
func (l *Logger) WithCaller(enable bool) { l.withCaller = enable }
func (l *Logger) WithJSON()              { l.jsonMode = true }

// File control
func (l *Logger) DisableFile() {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.logFile != nil {
        _ = l.logFile.Close()
        l.logFile = nil
    }
}

func (l *Logger) EnableFile(path string) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.logFile != nil {
        _ = l.logFile.Close()
    }
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil { return err }
    l.logFile = f
    return nil
}
