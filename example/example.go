package main

import (
	"github.com/leanfra/unilog"
	"go.uber.org/zap"
)

func main() {
	unilog.SetSimpleLogger("json", "./a.log", zap.DebugLevel)

	unilog.L().Debug("a log 1")
	unilog.L().Debug("a log 2")

	unilog.SetSimpleLogger("console", "/dev/stderr", zap.DebugLevel)
	unilog.L().Debug("b log 1")
	unilog.L().Debug("b log 2")

	unilog.L().Sugar().Error("b error 1")
}
