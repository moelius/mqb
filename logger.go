package mqb

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

const loggerName = "mqb"

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

var levelMap = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Panic(msg string)
	Fatal(msg string)
}

type LoggerLog struct {
	level Level
}

func NewLoggerLog(level Level) LoggerInterface {
	return &LoggerLog{level: level}
}

func (l *LoggerLog) log(level Level, msg string) {
	if level < l.level {
		return
	}
	m := fmt.Sprintf("\t%s\t%s %s\n", levelMap[level], loggerName, msg)
	switch level {
	case FatalLevel:
		log.Fatal(m)
	case PanicLevel:
		log.Panic(m)
	default:
		log.Print(m)
	}
}
func (l *LoggerLog) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l *LoggerLog) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *LoggerLog) Warning(msg string) {
	l.log(WarnLevel, msg)
}

func (l *LoggerLog) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *LoggerLog) Panic(msg string) {
	l.log(PanicLevel, msg)
}

func (l *LoggerLog) Fatal(msg string) {
	l.log(FatalLevel, msg)
}

type LoggerZap struct {
	logger *zap.Logger
}

func NewLoggerZap(level zapcore.Level) LoggerInterface {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleErrors := zapcore.Lock(os.Stderr)
	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= level })))
	logger := zap.New(core).Named(loggerName)
	logger.WithOptions()
	return &LoggerZap{logger: logger}
}

func (l *LoggerZap) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *LoggerZap) Info(msg string) {
	l.logger.Info(msg)
}

func (l *LoggerZap) Warning(msg string) {
	l.logger.Warn(msg)
}

func (l *LoggerZap) Error(msg string) {
	l.logger.Error(msg)
}

func (l *LoggerZap) Panic(msg string) {
	l.logger.Panic(msg)
}

func (l *LoggerZap) Fatal(msg string) {
	l.logger.Fatal(msg)
}
