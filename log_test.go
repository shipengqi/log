package log

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalLogger(t *testing.T) {
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
		"debug Hello, world!",
		"info Hello, world!",
		"warn Hello, world!",
		"error Hello, world!",
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
	opts.FilenameEncoder = func() string {
		return "test.log"
	}
	opts.DisableConsole = true
	opts.DisableFile = false
	opts.Output = tmp

	Configure(opts)
	Info("Hello, world!")
	assert.Equal(t, filepath.Join(tmp, "test.log"), EncodedFilename())
	_ = os.Remove(EncodedFilename())
}

func TestLoggerClose(t *testing.T) {
	t.Run("close logger without log file", func(t *testing.T) {
		str := "close logger without log file"
		opts := NewOptions()
		Configure(opts)
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
	es := NewErrSlice()
	assert.Equal(t, "", es.Error())
	es.Append(errors.New("error1"))
	assert.Equal(t, "error1", es.Error())

	es.AppendStr("error2")
	assert.Equal(t, "error1 : error2", es.Error())

	es.Append(errors.New("error3"))
	assert.Equal(t, "error1 : error2 : error3", es.Error())
}
