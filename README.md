# log

A structured logger for Go, based on [zap](https://github.com/uber-go/zap). 
Migrated from [golib](https://github.com/shipengqi/golib).

[![test](https://github.com/shipengqi/log/actions/workflows/test.yaml/badge.svg)](https://github.com/shipengqi/log/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/shipengqi/log/branch/main/graph/badge.svg?token=CQKD0I63DQ)](https://codecov.io/gh/shipengqi/log)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/log)](https://goreportcard.com/report/github.com/shipengqi/log)
[![release](https://img.shields.io/github/release/shipengqi/log.svg)](https://github.com/shipengqi/log/releases)
[![license](https://img.shields.io/github/license/shipengqi/log)](https://github.com/shipengqi/log/blob/main/LICENSE)

## Quick Start

```go
opts := log.Newoptions()
errs := opts.Validate()
if len(errs) > 0 {
	// handle errors
	return
}

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

## 🔋 JetBrains OS licenses

`log` had been being developed with **GoLand** under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=log" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo." width="250" align="middle"></a>