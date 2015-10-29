package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	_ Level = iota
	None
	Error
	Debug
	Warn
	Info
	Access
	All
)

type Logger struct {
	Level Level
	E     LevelLogger
	D     LevelLogger
	I     LevelLogger
	W     LevelLogger
	A     LevelLogger
}

type LevelLogger struct {
	FileWriter io.Writer
	StdWriter  io.Writer
}

type AccessLog struct {
	Time           time.Time     `json:"time,omitempty"`
	Protocol       string        `json:"protocol,omitempty"`
	HTTPStatusCode int           `json:"http_status_code,omitempty"`
	ResponseTime   time.Duration `json:"response_time,omitempty"`
	UserAgent      string        `json:"user_agent,omitempty"`
	URL            string        `json:"url, omitempty"`
	IP             string        `json:"ip,omitempty"`
	Method         string        `json:"method,omitempty"`
}

type GenericLog struct {
	Time    time.Time              `json:"time,omitempty"`
	Level   string                 `json:"level,omitempty"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type Data map[string]interface{}

type Level uint8

type Config struct {
	Level  Level
	Err    string
	Debug  string
	Info   string
	Warn   string
	Access string
}

var colors = map[string]string{
	"default": "\x1b[0m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"purple":  "\x1b[35m",
	"cyan":    "\x1b[36m",
	"grey":    "\x1b[37m",
}

func New(cfg Config) *Logger {
	var level Level = All

	if cfg.Level != 0 {
		level = cfg.Level
	}

	return &Logger{
		Level: level,
		E:     LevelLogger{FileWriter: setWriter(cfg.Err), StdWriter: os.Stderr},
		D:     LevelLogger{FileWriter: setWriter(cfg.Debug), StdWriter: os.Stdout},
		I:     LevelLogger{FileWriter: setWriter(cfg.Info), StdWriter: os.Stdout},
		W:     LevelLogger{FileWriter: setWriter(cfg.Warn), StdWriter: os.Stderr},
		A:     LevelLogger{FileWriter: setWriter(cfg.Access), StdWriter: os.Stdout},
	}
}

func (l *Logger) Error(msg string, data Data) {
	if l.Level < Error {
		return
	}
	gLog := GenericLog{
		Time:    time.Now(),
		Level:   "ERROR",
		Message: msg,
		Data:    data,
	}

	l.log(l.E, "ERROR", gLog)
}

func (l *Logger) Warn(msg string, data Data) {
	if l.Level < Warn {
		return
	}
	gLog := GenericLog{
		Time:    time.Now(),
		Level:   "WARN",
		Message: msg,
		Data:    data,
	}

	l.log(l.W, "WARN", gLog)
}

func (l *Logger) Debug(msg string, data Data) {
	if l.Level < Debug {
		return
	}
	gLog := GenericLog{
		Time:    time.Now(),
		Level:   "DEBUG",
		Message: msg,
		Data:    data,
	}

	l.log(l.D, "DEBUG", gLog)
}

func (l *Logger) Info(msg string, data Data) {
	if l.Level < Info {
		return
	}
	gLog := GenericLog{
		Time:    time.Now(),
		Level:   "INFO",
		Message: msg,
		Data:    data,
	}

	l.log(l.I, "INFO", gLog)
}

func (l *Logger) Access(startTime time.Time, w http.ResponseWriter, r *http.Request) {
	al := AccessLog{
		Time:           startTime,
		Protocol:       r.Proto,
		HTTPStatusCode: 200,
		ResponseTime:   time.Since(startTime),
		UserAgent:      r.UserAgent(),
		IP:             "123.123.123.123",
		Method:         r.Method,
		URL:            r.URL.String(),
	}

	l.log(l.A, "ACCESS", al)
}

func (l *Logger) log(w LevelLogger, level string, log interface{}) {

	var out_message string

	switch t := log.(type) {
	case AccessLog:
		out_message = fmt.Sprintf("\n%s - - [%s] \"%s %s %s\" %d - %s\n", t.IP, t.Time, t.Method, t.URL, t.Protocol, t.HTTPStatusCode, t.ResponseTime)
	case GenericLog:
		var coloredLogLevel string

		switch level {
		case "ERROR":
			coloredLogLevel = fmt.Sprintf("%s[%s]%s", colors["red"], level, colors["default"])
		case "WARN":
			coloredLogLevel = fmt.Sprintf("%s[%s]%s", colors["yellow"], level, colors["default"])
		case "INFO":
			coloredLogLevel = fmt.Sprintf("%s[%s]%s", colors["green"], level, colors["default"])
		case "DEBUG":
			coloredLogLevel = fmt.Sprintf("%s[%s]%s", colors["cyan"], level, colors["default"])
		}

		data, _ := json.Marshal(t.Data)
		out_message = fmt.Sprintf("\n%s %s %s %v\n", coloredLogLevel, t.Time, t.Message, string(data))
	default:
		_ = t
		out_message = "default"
	}

	// Print to std
	fmt.Fprintf(w.StdWriter, out_message)

	json, err := json.Marshal(log)
	if err != nil {
		return
	}

	// Print to file
	fmt.Fprintf(w.FileWriter, "\n%s\n", string(json))
}

func setWriter(filePath string, defaultWriters ...io.Writer) io.Writer {
	if filePath != "" {
		fileWriter, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			panic("xx")
		}

		defaultWriters = append(defaultWriters, fileWriter)
	}

	return io.MultiWriter(defaultWriters...)
}
