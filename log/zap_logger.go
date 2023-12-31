package log

import (
	"context"
	"fmt"
	"os"

	"github.com/easonchen147/foundation/cfg"
	"github.com/easonchen147/foundation/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AccessLogger *zap.Logger
	Logger       *zap.Logger
	SqlLogger    *zap.Logger

	lumberJackLoggerDefault *lumberjack.Logger
	lumberJackLoggerAccess  *lumberjack.Logger
	lumberJackLoggerSql     *lumberjack.Logger
)

func init() {
	InitLog(cfg.AppConf)
}

// InitLog 配置日志模块
func InitLog(cfg *cfg.AppConfig) {
	var level zapcore.Level
	if level.UnmarshalText([]byte(cfg.LogLevel)) != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "time",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		CallerKey:      "location",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	lumberJackLoggerDefault = newLunmberJackLogger(cfg.LogFile)
	lumberJackLoggerAccess = newLunmberJackLogger(cfg.AccessLogFile)
	lumberJackLoggerSql = newLunmberJackLogger(cfg.SqlLogFile)

	var defaultCore, accessCore, sqlCore zapcore.Core
	switch cfg.LogMode {
	case "console":
		defaultCore = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
		accessCore = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
		sqlCore = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	case "file":
		defaultCore = newLoggerCore(lumberJackLoggerDefault, encoderConfig, level)
		accessCore = newLoggerCore(lumberJackLoggerAccess, encoderConfig, level)
		sqlCore = newLoggerCore(lumberJackLoggerSql, encoderConfig, level)
	}

	Logger = zap.New(defaultCore, zap.AddCaller(), zap.AddCallerSkip(1))
	AccessLogger = zap.New(accessCore, zap.AddCaller(), zap.AddCallerSkip(1))
	SqlLogger = zap.New(sqlCore, zap.AddCaller(), zap.AddCallerSkip(1))
}

func newLoggerCore(logger *lumberjack.Logger, encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
	writer := zapcore.AddSync(logger)
	return zapcore.NewTee(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), level))
}

func newLunmberJackLogger(logFilePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 0,
		MaxAge:     30, // days
		LocalTime:  true,
		Compress:   true,
	}
}

func Debug(ctx context.Context, msg string, val ...interface{}) {
	Logger.Debug(fmt.Sprintf(msg, val...), zapDefaultFields(ctx)...)
}

func Info(ctx context.Context, msg string, val ...interface{}) {
	Logger.Info(fmt.Sprintf(msg, val...), zapDefaultFields(ctx)...)
}

func Warn(ctx context.Context, msg string, val ...interface{}) {
	Logger.Warn(fmt.Sprintf(msg, val...), zapDefaultFields(ctx)...)
}

func Error(ctx context.Context, msg string, val ...interface{}) {
	Logger.Error(fmt.Sprintf(msg, val...), zapDefaultFields(ctx)...)
}

func Panic(ctx context.Context, msg string, val ...interface{}) {
	Logger.Panic(fmt.Sprintf(msg, val...), zapDefaultFields(ctx)...)
}

func Access(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zapDefaultFields(ctx)...)
	AccessLogger.Info(msg, fields...)
}

func zapDefaultFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	fields = append(fields, zap.String("traceId", getTraceId(ctx)))
	return fields
}

func getTraceId(ctx context.Context) string {
	obj := ctx.Value(constant.TraceIdKey)
	if obj == nil {
		return ""
	}
	traceId, ok := obj.(string)
	if !ok {
		return ""
	}
	return traceId
}

func Close() {
	_ = lumberJackLoggerDefault.Rotate()
	_ = lumberJackLoggerAccess.Rotate()
	_ = lumberJackLoggerSql.Rotate()
}
