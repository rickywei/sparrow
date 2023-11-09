package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rickywei/sparrow/project/conf"
)

func init() {
	level := zapcore.DebugLevel
	switch conf.String("log.level") {
	case "error":
		level = zapcore.ErrorLevel
	case "warn":
		level = zapcore.WarnLevel
	case "info":
		level = zapcore.InfoLevel
	case "debug":
		level = zapcore.DebugLevel
	}
	fw := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     15,
		LocalTime:  false,
		Compress:   false,
	}

	cores := []zapcore.Core{zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(fw),
		level,
	)}
	if conf.Bool("log.console") {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(zapcore.Lock(os.Stderr)),
			level,
		),
		)
	}
	zap.ReplaceGlobals(zap.New(zapcore.NewTee(cores...)))
}

func L() *zap.Logger {
	return zap.L()
}
