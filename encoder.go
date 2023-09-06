package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// FilenameEncoder log filename encoder,
// return the full name of the log file.
type FilenameEncoder func() string

// DefaultFilenameEncoder return <process name>-<date>.log.
func DefaultFilenameEncoder() string {
	return fmt.Sprintf("%s-%s.log", filepath.Base(os.Args[0]), time.Now().Format("20060102"))
}

func DefaultTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func rollingFileEncoder(opts *Options, encoder FilenameEncoder) (zapcore.WriteSyncer, io.Closer, string) {
	encoded := encoder()
	f := filepath.Join(opts.Output, encoded)
	if opts.DisableRotate {
		fd, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o644)
		if err != nil {
			panic(err)
		}
		return zapcore.AddSync(fd), fd, f
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	jackl := &lumberjack.Logger{
		Filename:   f,
		MaxSize:    opts.MaxSize,
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxBackups,
	}
	return zapcore.AddSync(jackl), jackl, f
}
