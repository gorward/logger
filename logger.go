package logger

import (
	"fmt"
	"io"
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
	All
)

type Logger struct {
	Level Level
	E     io.Writer
	D     io.Writer
	I     io.Writer
	W     io.Writer
}

type Level uint8

type Config struct {
	Level Level
	Err   string
	Debug string
	Info  string
	Warn  string
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
		E:     setWriter(cfg.Err, os.Stderr),
		D:     setWriter(cfg.Debug, os.Stdout),
		I:     setWriter(cfg.Info, os.Stdout),
		W:     setWriter(cfg.Warn, os.Stdout),
	}
}

func (l *Logger) Error(msg string) {
	if l.Level < Error {
		return
	}

	var prefix string = fmt.Sprintf("%s[ERROR] %s", colors["red"], colors["default"])
	l.log(l.E, prefix, msg)
}

func (l *Logger) Warn(msg string) {
	if l.Level < Warn {
		return
	}

	var prefix string = fmt.Sprintf("%s[WARNING] %s", colors["yellow"], colors["default"])
	l.log(l.W, prefix, msg)
}

func (l *Logger) Info(msg string) {
	if l.Level < Info {
		return
	}

	var prefix string = fmt.Sprintf("%s[INFO] %s", colors["blue"], colors["default"])
	l.log(l.I, prefix, msg)
}

func (l *Logger) Debug(msg string) {
	if l.Level < Debug {
		return
	}

	var prefix string = fmt.Sprintf("%s[DEBUG] %s", colors["white"], colors["default"])
	l.log(l.D, prefix, msg)
}

func (l *Logger) log(w io.Writer, prefix string, msg string) {
	fmt.Fprintf(w, time.Now().Format(time.RFC3339)+" "+prefix+msg+"\n")
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
