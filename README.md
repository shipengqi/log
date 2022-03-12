# log

A structured logger for Go, based on [zap](https://github.com/uber-go/zap). 
Migrated from [golib](https://github.com/shipengqi/golib).

[![ci](https://github.com/shipengqi/log/actions/workflows/ci.yml/badge.svg)](https://github.com/shipengqi/log/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/shipengqi/log/branch/main/graph/badge.svg?token=CQKD0I63DQ)](https://codecov.io/gh/shipengqi/log)
[![release](https://img.shields.io/github/release/shipengqi/log.svg)](https://github.com/shipengqi/log/releases)
[![license](https://img.shields.io/github/license/shipengqi/log)](https://github.com/shipengqi/log/blob/main/LICENSE)

## Quick Start

```go
opts := &log.Newoptions()
log.Configure(opts)
defer func() {
    _ = log.Close()
}()

log.Debug("debug message")
log.Info("info message")
log.Warn("warn message")
log.Error("error message")

log.Debugf("%s message", "debug")
log.Infof("%s message", "info")
log.Warnf("%s message", "warn")
log.Errorf("%s message", "error")

log.Debugt("debug message", log.String("key1", "value1"))
log.Infot("info message", log.Int32("key2", 10))
log.Warnt("warn message", log.Bool("key3", false))
log.Errort("error message", log.Any("key4", "any"))
```

## Documentation

You can find the docs at [go docs](https://pkg.go.dev/github.com/shipengqi/log).
