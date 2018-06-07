package main

import (
	"go.uber.org/zap/zapcore"
)

type pLevel struct {
	zapcore.Level
}

func (*pLevel) Type() string {
	return "string"
}
