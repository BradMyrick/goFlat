package obs

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key int

const (
	keyLogger key = iota
)

func NewZap() (logger *zap.Logger, err error) {

	var level zap.AtomicLevel
	switch viper.GetString("log.level") {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} // switch root.log.level

	if logger, err = (zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    1000,
			Thereafter: 1000,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()); err != nil {
		logger = nil
		return
	}
	zap.RedirectStdLog(logger)
	return

} // func obs.NewZap()

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {

	return context.WithValue(ctx, keyLogger, logger)

} // func obs.WithLogger()

func Logger(ctx context.Context) *zap.Logger {

	if logger, ok := ctx.Value(keyLogger).(*zap.Logger); ok {
		return logger
	}

	return zap.NewNop()

} // func obs.Logger()
