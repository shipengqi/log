package log

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestGlobalLogger(t *testing.T) {
	t.Run("Base", func(t *testing.T) {
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()
		os.Stdout = w

		name := "world"
		str := "Hello, world!"
		opts := NewOptions()
		Configure(opts)
		Debugf("Hello, %s", name+"1")
		Infof("Hello, %s", name+"2")
		Warnf("Hello, %s", name+"3")
		Errorf("Hello, %s", name+"4")

		opts.ConsoleLevel = DebugLevel.String()
		Configure(opts)
		Debug(str)
		Info(str)
		Warn(str)
		Error(str)

		opts.DisableConsoleColor = true
		Configure(opts)

		L().Debug(str)
		L().Info(str)
		L().Warn(str)
		L().Error(str)

		expected := []string{
			"\x1b[34mINFO\x1b[0m Hello, world2",
			"\x1b[33mWARN\x1b[0m Hello, world3",
			"\x1b[31mERROR\x1b[0m Hello, world4",
			"\x1b[35mDEBUG\x1b[0m Hello, world!",
			"\x1b[34mINFO\x1b[0m Hello, world!",
			"\x1b[33mWARN\x1b[0m Hello, world!",
			"\x1b[31mERROR\x1b[0m Hello, world!",
			"DEBUG Hello, world!",
			"INFO Hello, world!",
			"WARN Hello, world!",
			"ERROR Hello, world!",
		}
		_ = w.Close()
		stdout, _ := ioutil.ReadAll(r)
		reader := bytes.NewReader(stdout)
		scanner := bufio.NewScanner(reader)
		for _, v := range expected {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			assert.Contains(t, line, v)
		}
	})

	t.Run("With Caller", func(t *testing.T) {
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()
		os.Stdout = w

		name := "world"
		opts := NewOptions()
		opts.ConsoleLevel = DebugLevel.String()
		opts.DisableConsoleColor = true
		opts.DisableConsoleCaller = false
		Configure(opts)
		L().Debugf("Hello, %s", name+"1")
		L().Infof("Hello, %s", name+"2")
		L().Warnf("Hello, %s", name+"3")
		L().Errorf("Hello, %s", name+"4")
		expected := []string{
			"DEBUG log/log_test.go:93 Hello, world1",
			"INFO log/log_test.go:94 Hello, world2",
			"WARN log/log_test.go:95 Hello, world3",
			"ERROR log/log_test.go:96 Hello, world4",
		}
		_ = w.Close()
		stdout, _ := ioutil.ReadAll(r)
		reader := bytes.NewReader(stdout)
		scanner := bufio.NewScanner(reader)
		for _, v := range expected {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			assert.Contains(t, line, v)
		}
	})
	t.Run("With Level encoder", func(t *testing.T) {
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()
		os.Stdout = w

		name := "world"
		opts := NewOptions()
		opts.ConsoleLevel = DebugLevel.String()
		opts.DisableConsoleColor = true
		opts.DisableConsoleCaller = false
		Configure(opts, WithLevelEncoder(zapcore.LowercaseLevelEncoder))
		L().Debugf("Hello, %s", name+"1")
		L().Infof("Hello, %s", name+"2")
		L().Warnf("Hello, %s", name+"3")
		L().Errorf("Hello, %s", name+"4")
		expected := []string{
			"debug log/log_test.go:129 Hello, world1",
			"info log/log_test.go:130 Hello, world2",
			"warn log/log_test.go:131 Hello, world3",
			"error log/log_test.go:132 Hello, world4",
		}
		_ = w.Close()
		stdout, _ := ioutil.ReadAll(r)
		reader := bytes.NewReader(stdout)
		scanner := bufio.NewScanner(reader)
		for _, v := range expected {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			assert.Contains(t, line, v)
		}
	})

	t.Run("With Caller encoder", func(t *testing.T) {
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()
		os.Stdout = w

		name := "world"
		opts := NewOptions()
		opts.ConsoleLevel = DebugLevel.String()
		opts.DisableConsoleColor = true
		opts.DisableConsoleCaller = false
		opts.CallerSkip = -1
		Configure(opts,
			WithLevelEncoder(zapcore.LowercaseLevelEncoder),
			WithCallerEncoder(zapcore.ShortCallerEncoder),
		)
		L().Debugf("Hello, %s", name+"1")
		L().Infof("Hello, %s", name+"2")
		L().Warnf("Hello, %s", name+"3")
		L().Errorf("Hello, %s", name+"4")
		expected := []string{
			"debug log/log_test.go:170 Hello, world1",
			"info log/log_test.go:171 Hello, world2",
			"warn log/log_test.go:172 Hello, world3",
			"error log/log_test.go:173 Hello, world4",
		}
		_ = w.Close()
		stdout, _ := ioutil.ReadAll(r)
		reader := bytes.NewReader(stdout)
		scanner := bufio.NewScanner(reader)
		for _, v := range expected {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			assert.Contains(t, line, v)
		}
	})
}

func TestLoggerPanic(t *testing.T) {
	str := "test panic"
	opts := NewOptions()
	Configure(opts)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, err, str)
		} else {
			t.Fatal("no panic")
		}
	}()
	Panic(str)
}

func TestWithValues(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()
	os.Stdout = w
	opts := NewOptions()
	Configure(opts)

	logger := WithValues(String("test key", "test value"))
	logger.Info("Hello, world!")

	_ = w.Close()
	stdout, _ := ioutil.ReadAll(r)
	assert.Contains(t, string(stdout), "Hello, world! {\"test key\": \"test value\"}")
}

func TestGlobalLoggerWithoutTime(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()
	os.Stdout = w
	opts := NewOptions()
	opts.DisableConsoleTime = true
	Configure(opts)

	Info("Hello, world!")
	_ = w.Close()
	stdout, _ := ioutil.ReadAll(r)
	assert.Equal(t, "\u001B[34mINFO\u001B[0m Hello, world!\n", string(stdout))
}

func TestLoggerFile(t *testing.T) {
	tmp := os.TempDir()
	opts := NewOptions()
	opts.DisableConsole = true
	opts.DisableFile = false
	opts.Output = tmp

	Configure(opts, WithFilenameEncoder(func() string {
		return "test.log"
	}))
	Info("Hello, world!")
	assert.Equal(t, filepath.Join(tmp, "test.log"), EncodedFilename())
	_ = os.Remove(EncodedFilename())
}

func TestLoggerClose(t *testing.T) {
	t.Run("close logger without log file", func(t *testing.T) {
		str := "close logger without log file"
		opts := NewOptions()
		Configure(opts, WithTimeEncoder(DefaultTimeEncoder))
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, err, str)
				cerr := Close()
				assert.NoError(t, cerr)
			} else {
				t.Fatal("no panic")
			}
		}()
		Panic(str)
	})
	t.Run("close logger with log file", func(t *testing.T) {
		str := "close logger with log file"
		opts := NewOptions()
		opts.DisableFile = false
		opts.DisableConsole = true
		opts.Output = "testdata/log"
		Configure(opts)
		Info(str)
		Info(EncodedFilename())
		content, err := ioutil.ReadFile(EncodedFilename())
		assert.NoError(t, err)
		strings.Contains(string(content), str)
		err = Close()
		assert.NoError(t, err)
		_ = os.Remove(EncodedFilename())
	})
	t.Run("close logger with rotate log file", func(t *testing.T) {
		str := "close logger with rotate log file"
		opts := NewOptions()
		opts.DisableFile = false
		opts.DisableRotate = false
		opts.DisableConsole = true
		opts.Output = "testdata/log"
		Configure(opts)
		Info(str)
		Info(EncodedFilename())
		err := Close()
		assert.NoError(t, err)
		content, err := ioutil.ReadFile(EncodedFilename())
		assert.NoError(t, err)
		strings.Contains(string(content), str)
		_ = os.Remove(EncodedFilename())
	})
}

func TestErrSlice(t *testing.T) {
	t.Run("Generic ErrSlice", func(t *testing.T) {
		es := NewErrSlice()
		assert.Equal(t, "", es.Error())
		es.Append(errors.New("error1"))
		assert.Equal(t, "error1", es.Error())

		es.AppendStr("error2")
		assert.Equal(t, "error1 : error2", es.Error())

		es.Append(errors.New("error3"))
		assert.Equal(t, "error1 : error2 : error3", es.Error())
	})

	t.Run("Append ErrSlice", func(t *testing.T) {
		es := NewErrSlice()
		es.Append(NewErrSlice())
		assert.Equal(t, 0, es.Len())

		es.Append(errors.New("error1"))
		assert.Equal(t, 1, es.Len())

		es2 := NewErrSlice()
		es2.Append(errors.New("error1"))
		es.Append(es2)
		assert.Equal(t, 2, es.Len())
	})
}

func TestStdInfoLogger(t *testing.T) {
	opts := NewOptions()
	opts.DisableFile = true
	opts.DisableConsole = false
	Configure(opts)
	logger := StdLogger(InfoLevel)
	assert.NotNil(t, logger)

	var (
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	ti := time.Now()
	logger.Printf(traceStr, fileWithLineNum(), float64(ti.Nanosecond())/1e6, "-", "test message")
	logger.Printf(traceErrStr, fileWithLineNum(), "terror", float64(ti.Nanosecond())/1e6, "-", "test error message")

	var (
		Reset       = "\033[0m"
		Green       = "\033[32m"
		Yellow      = "\033[33m"
		BlueBold    = "\033[34;1m"
		MagentaBold = "\033[35;1m"
		RedBold     = "\033[31;1m"
	)
	traceStr = Green + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	traceErrStr = RedBold + "%s " + MagentaBold + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"

	logger.Printf(traceStr, fileWithLineNum(), float64(ti.Nanosecond())/1e6, "-", "color message")
	logger.Printf(traceErrStr, fileWithLineNum(), "terror", float64(ti.Nanosecond())/1e6, "-", "color error message")

	t.Run("Nil StdInfoLogger", func(t *testing.T) {
		tmp := _globalL
		_globalL = nil

		nlogger := StdLogger(InfoLevel)
		assert.Nil(t, nlogger)

		_globalL = tmp
	})
}


// fileWithLineNum return the file name and line number of the current file
func fileWithLineNum() string {
	for i := 4; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.HasSuffix(file, "_test.go") {
			dir, f := filepath.Split(file)
			return filepath.Join(filepath.Base(dir), f) + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
