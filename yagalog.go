package logger

import (
	"encoding/json"
	"fmt"
	"io"
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
	level         LogLevel
	out           io.Writer // destination for console loggers (defaults to os.Stdout)
	filePath      string    // optional file path set via options
	timeFormat    string
	withCaller    bool
	jsonMode      bool
}

// Option configures a Logger during construction.
type Option func(*Logger)

// WithFilePath sets a file path to open and write logs into. NewLogger will
// attempt to open the file and may return an error if it fails.
func WithFilePath(path string) Option {
	return func(l *Logger) { l.filePath = path }
}

// WithWriter sets the destination used for console loggers (defaults to os.Stdout).
func WithWriter(w io.Writer) Option {
	return func(l *Logger) {
		if w != nil {
			l.out = w
		}
	}
}

func WithLevel(level LogLevel) Option { return func(l *Logger) { l.level = level } }
func WithJSON() Option                { return func(l *Logger) { l.jsonMode = true } }
func WithTimeFormat(layout string) Option {
	return func(l *Logger) {
		if layout != "" {
			l.timeFormat = layout
		}
	}
}
func WithCaller(enable bool) Option { return func(l *Logger) { l.withCaller = enable } }
func WithColors(enable bool) Option { return func(l *Logger) { color.NoColor = !enable } }

// NewLogger constructs a Logger with functional options. Options are applied
// before any file is opened, so options like WithFilePath should be provided
// to have effect. NewLogger may return an error if opening the configured
// file fails.
func NewLogger(opts ...Option) (*Logger, error) {
	// sensible defaults
	l := &Logger{
		out:        os.Stdout,
		level:      DEBUG,
		timeFormat: "2006-01-02 15:04:05",
	}

	// apply options
	for _, o := range opts {
		o(l)
	}

	// auto color: disable if NO_COLOR or not a TTY; enable if FORCE_COLOR
	if os.Getenv("FORCE_COLOR") != "" {
		color.NoColor = false
	} else if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	// initialize console loggers using l.out
	l.infoLogger = log.New(l.out, color.CyanString("[INFO] "), log.Ltime)
	l.errorLogger = log.New(l.out, color.RedString("[ERROR] "), log.Ltime)
	l.debugLogger = log.New(l.out, color.GreenString("[DEBUG] "), log.Ltime)
	l.warningLogger = log.New(l.out, color.YellowString("[WARNING] "), log.Ltime)
	l.fatalLogger = log.New(l.out, color.MagentaString("[FATAL] "), log.Ltime)

	// if a file path was provided, attempt to open it
	if l.filePath != "" {
		f, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		l.logFile = f
	}

	return l, nil
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
			"level": levelStr[1 : len(levelStr)-2], // remove brackets and space
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
func (l *Logger) WithTimeFormat(layout string) {
	if layout != "" {
		l.timeFormat = layout
	}
}
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
	if err != nil {
		return err
	}
	l.logFile = f
	return nil
}
