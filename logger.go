package mqb

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

const loggerName = "mqb"

// A Level is a logging priority. Higher levels are more important.
type Level int8

//The same as zap levels, be careful later if want to change this levels
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

var zapMqbLevelMap = map[Level]zapcore.Level{
	DebugLevel: zap.DebugLevel,
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	ErrorLevel: zap.ErrorLevel,
	PanicLevel: zap.PanicLevel,
	FatalLevel: zap.FatalLevel,
}

func logF(level Level, msg string) string {
	return fmt.Sprintf("\t%s\t%s %s\n", levelMap[level], loggerName, msg)
}

//LoggerInterface interface for logging
type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Panic(msg string)
	Fatal(msg string)
	GetLevel() Level
}

//LoggerLog standard log implementation
type LoggerLog struct {
	level Level
}

//NewLoggerLog creates new standard logger
func NewLoggerLog(level Level) LoggerInterface {
	return &LoggerLog{level: level}
}

func (l *LoggerLog) log(level Level, msg string) {
	if level < l.level {
		return
	}
	switch level {
	case FatalLevel:
		log.Fatal(logF(level, msg))
	case PanicLevel:
		log.Panic(logF(level, msg))
	default:
		log.Print(logF(level, msg))
	}
}

//GetLevel get logger level
func (l *LoggerLog) GetLevel() Level {
	return l.level
}

//Debug log with level DebugLevel
func (l *LoggerLog) Debug(msg string) {
	l.log(DebugLevel, msg)
}

//Info log with level InfoLevel
func (l *LoggerLog) Info(msg string) {
	l.log(InfoLevel, msg)
}

//Warning log with level WarnLevel
func (l *LoggerLog) Warning(msg string) {
	l.log(WarnLevel, msg)
}

//Error log with level ErrorLevel
func (l *LoggerLog) Error(msg string) {
	l.log(ErrorLevel, msg)
}

//Fatal log with level FatalLevel
func (l *LoggerLog) Fatal(msg string) {
	l.log(FatalLevel, msg)
}

//Panic log with level PanicLevel
func (l *LoggerLog) Panic(msg string) {
	l.log(PanicLevel, msg)
}

//LoggerZap go.uber.org/zap logging library implementation
type LoggerZap struct {
	logger *zap.Logger
	level  Level
}

//NewLoggerZap creates new zap logger
func NewLoggerZap(level Level) LoggerInterface {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleErrors := zapcore.Lock(os.Stderr)
	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= zapMqbLevelMap[level] })))
	logger := zap.New(core).Named(loggerName)

	return &LoggerZap{logger: logger, level: level}
}

//GetLevel get logger level
func (l *LoggerZap) GetLevel() Level {
	return l.level
}

//Debug log with level DebugLevel
func (l *LoggerZap) Debug(msg string) {
	l.logger.Debug(msg)
}

//Info log with level InfoLevel
func (l *LoggerZap) Info(msg string) {
	l.logger.Info(msg)
}

//Warning log with level WarnLevel
func (l *LoggerZap) Warning(msg string) {
	l.logger.Warn(msg)
}

//Error log with level ErrorLevel
func (l *LoggerZap) Error(msg string) {
	l.logger.Error(msg)
}

//Fatal log with level FatalLevel
func (l *LoggerZap) Fatal(msg string) {
	l.logger.Fatal(msg)
}

//Panic log with level PanicLevel
func (l *LoggerZap) Panic(msg string) {
	l.logger.Panic(msg)
}

//LoggerDummy dummy logger struct for tests
type LoggerDummy struct {
	logger *bytes.Buffer
	level  Level
}

//NewLoggerDummy creates new dummy logger for tests and etc...
func NewLoggerDummy(level Level) *LoggerDummy {
	return &LoggerDummy{logger: bytes.NewBuffer([]byte{}), level: level}
}

func (l *LoggerDummy) log(level Level, msg string) {
	if level < l.level {
		return
	}
	l.logger.Write([]byte(logF(level, msg)))
}

//GetLevel get logger level
func (l *LoggerDummy) GetLevel() Level {
	return l.level
}

//Debug log with level DebugLevel
func (l *LoggerDummy) Debug(msg string) {
	l.log(DebugLevel, msg)
}

//Info log with level InfoLevel
func (l *LoggerDummy) Info(msg string) {
	l.log(InfoLevel, msg)
}

//Warning log with level WarnLevel
func (l *LoggerDummy) Warning(msg string) {
	l.log(WarnLevel, msg)
}

//Error log with level ErrorLevel
func (l *LoggerDummy) Error(msg string) {
	l.log(ErrorLevel, msg)
}

//Fatal log with level FatalLevel
func (l *LoggerDummy) Fatal(msg string) {
	l.log(FatalLevel, msg)
}

//Panic log with level PanicLevel
func (l *LoggerDummy) Panic(msg string) {
	l.log(PanicLevel, msg)
}

//Panic log with level PanicLevel
func (l *LoggerDummy) Read() string {
	s, err := l.logger.ReadString(byte('\n'))
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	l.logger.Reset()
	return s
}
