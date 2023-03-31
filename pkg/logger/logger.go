package logger

import (
	"context"
	"github.com/Mistolotus/tdl/pkg/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
)

type ctxKey struct{}

func From(ctx context.Context) *zap.Logger {
	return ctx.Value(ctxKey{}).(*zap.Logger)
}

func With(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Named(ctx context.Context, name string) context.Context {
	return With(ctx, From(ctx).Named(name))
}

func New(level zapcore.LevelEnabler) *zap.Logger {
	rotate := &lumberjack.Logger{
		Filename:   filepath.Join(consts.LogPath, "latest.log"),
		MaxSize:    10,
		MaxAge:     7,
		MaxBackups: 3,
		LocalTime:  true,
		Compress:   true,
	}

	writer := zapcore.AddSync(rotate)

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(config), writer, level)
	return zap.New(core, zap.AddCaller())
}
