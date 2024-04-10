// Package log is a structured logger for Go, based on https://github.com/uber-go/zap.
package log

import (
	"log"

	"go.uber.org/zap"
)

var _globalL *Logger

func init() {
	_globalL = New(NewOptions())
}

// L returns the global logger.
func L() *Logger {
	return _globalL
}

// EncodedFilename returns the filename for logging when DisableFile is false.
func EncodedFilename() string {
	return _globalL.EncodedFilename()
}

// SugaredL returns global sugared logger.
func SugaredL() *zap.SugaredLogger {
	return _globalL.Sugared()
}

// StdLogger returns logger of standard library which writes to supplied zap
// logger at the given level.
func StdLogger(level Level) *log.Logger {
	if _globalL == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(_globalL.log, level); err == nil {
		return l
	}

	return nil
}

type DebugLogger interface {
	Debugt(msg string, fields ...Field)
	Debugf(template string, args ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

type InfoLogger interface {
	DebugLogger

	Infot(msg string, fields ...Field)
	Infof(template string, args ...interface{})
	Info(msg string, keysAndValues ...interface{})
}

type Interface interface {
	InfoLogger

	Warnt(msg string, fields ...Field)
	Warnf(template string, args ...interface{})
	Warn(msg string, keysAndValues ...interface{})

	Errort(msg string, fields ...Field)
	Errorf(template string, args ...interface{})
	Error(msg string, keysAndValues ...interface{})

	Panict(msg string, fields ...Field)
	Panicf(template string, args ...interface{})
	Panic(msg string, keysAndValues ...interface{})

	Fatalt(msg string, fields ...Field)
	Fatalf(template string, args ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	AtLevelt(level Level, msg string, fields ...Field)
	AtLevelf(level Level, template string, args ...interface{})
	AtLevel(level Level, msg string, keysAndValues ...interface{})

	// WithValues creates a child logger and adds some Field of
	// context to this logger.
	WithValues(fields ...Field) *Logger

	// Flush calls the underlying Core's Sync method, flushing any buffered
	// log entries. Applications should take care to call Sync before exiting.
	Flush() error

	// Close implements io.Closer, and closes the current logfile.
	Close() error
}

// Configure sets up the global logger.
func Configure(opts *Options) {
	l := New(opts)
	_globalL = l
	zap.RedirectStdLog(_globalL.log)
}

// Debugt logs a message at DebugLevel.
func Debugt(msg string, fields ...Field) {
	_globalL.Debugt(msg, fields...)
}

// Debugf logs a message at DebugLevel.
func Debugf(template string, args ...interface{}) {
	_globalL.Debugf(template, args...)
}

// Debug logs a message at DebugLevel.
func Debug(msg string, keysAndValues ...interface{}) {
	_globalL.Debug(msg, keysAndValues...)
}

// Infot logs a message at InfoLevel.
func Infot(msg string, fields ...Field) {
	_globalL.Infot(msg, fields...)
}

// Infof logs a message at InfoLevel.
func Infof(template string, args ...interface{}) {
	_globalL.Infof(template, args...)
}

// Info logs a message at InfoLevel.
func Info(msg string, keysAndValues ...interface{}) {
	_globalL.Info(msg, keysAndValues...)
}

// Warnt logs a message at WarnLevel.
func Warnt(msg string, fields ...Field) {
	_globalL.Warnt(msg, fields...)
}

// Warnf logs a message at WarnLevel.
func Warnf(template string, args ...interface{}) {
	_globalL.Warnf(template, args...)
}

// Warn logs a message at WarnLevel.
func Warn(msg string, keysAndValues ...interface{}) {
	_globalL.Warn(msg, keysAndValues...)
}

// Errort logs a message at ErrorLevel.
func Errort(msg string, fields ...Field) {
	_globalL.Errort(msg, fields...)
}

// Errorf logs a message at ErrorLevel.
func Errorf(template string, args ...interface{}) {
	_globalL.Errorf(template, args...)
}

// Error logs a message at ErrorLevel.
func Error(msg string, keysAndValues ...interface{}) {
	_globalL.Error(msg, keysAndValues...)
}

// Panict logs a message at PanicLevel, then panics.
func Panict(msg string, fields ...Field) {
	_globalL.Panict(msg, fields...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	_globalL.Panicf(template, args...)
}

// Panic logs a message with some additional context, then panics.
func Panic(msg string, keysAndValues ...interface{}) {
	_globalL.Panic(msg, keysAndValues...)
}

// Fatalt logs a message at FatalLevel, then calls os.Exit(1).
func Fatalt(msg string, fields ...Field) {
	_globalL.Fatalt(msg, fields...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	_globalL.Fatalf(template, args...)
}

// Fatal logs a message with some additional context, then calls os.Exit.
func Fatal(msg string, keysAndValues ...interface{}) {
	_globalL.Fatal(msg, keysAndValues...)
}

// Print logs a message at level Print.
func Print(args ...interface{}) {
	_globalL.Print(args...)
}

// Println logs a message at level Print.
func Println(args ...interface{}) {
	_globalL.Println(args...)
}

// Printf logs a message at level Print.
func Printf(format string, args ...interface{}) {
	_globalL.Printf(format, args...)
}

// AtLevelt logs a message at Level.
func AtLevelt(level Level, msg string, fields ...Field) {
	_globalL.AtLevelt(level, msg, fields...)
}

// AtLevelf logs a message at Level.
func AtLevelf(level Level, msg string, args ...interface{}) {
	_globalL.AtLevelf(level, msg, args...)
}

// AtLevel logs a message at Level.
func AtLevel(level Level, msg string, keysAndValues ...interface{}) {
	_globalL.AtLevel(level, msg, keysAndValues...)
}

// WithValues creates a child logger and adds some Field of
// context to this logger.
func WithValues(fields ...Field) *Logger {
	return _globalL.WithValues(fields...)
}

// Flush calls the underlying Core's Sync method, flushing any buffered
// log entries. Applications should take care to call Sync before exiting.
func Flush() error { return _globalL.Flush() }

// Close implements io.Closer, and closes the current logfile of default logger.
func Close() error { return _globalL.Close() }

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func Check(lvl Level, msg string) *CheckedEntry {
	return _globalL.Check(lvl, msg)
}
