package pkglogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-oms/internal/domain/entities"
	"go-oms/utils"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

type PrettyJSONEncoder struct {
	zapcore.Encoder
}

func (e PrettyJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, buf.Bytes(), "", "  ")
	if err != nil {
		return nil, err
	}

	buf = buffer.NewPool().Get()
	buf.AppendBytes(pretty.Bytes())
	return buf, nil
}

func NewPrettyJSONEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return PrettyJSONEncoder{Encoder: zapcore.NewJSONEncoder(cfg)}
}

func NewLogger(level string, development bool) *zap.SugaredLogger {
	var cores []zapcore.Core
	var encoderConfig zapcore.EncoderConfig

	nowStr := strings.Split(time.Now().Format("2006-01-02"), "-")
	logPath := fmt.Sprintf("./assets/logger/%s_%s/%s", nowStr[0], nowStr[1], nowStr[2])
	utils.EnsureFolderExists(logPath)
	pathInfo := fmt.Sprintf("%s/info.json", logPath)
	pathError := fmt.Sprintf("%s/error.json", logPath)
	fileInfo := zapcore.AddSync(&lumberjack.Logger{
		Filename:   pathInfo,
		MaxSize:    L.LogSize,
		MaxBackups: L.LogBackups,
		MaxAge:     L.LogAge,
	})

	fileError := zapcore.AddSync(&lumberjack.Logger{
		Filename:   pathError,
		MaxSize:    L.LogSize,
		MaxBackups: L.LogBackups,
		MaxAge:     L.LogAge,
	})
	if development {
		encoderConfig = developmentEncoderConfig
		cores = append(cores, zapcore.NewCore(NewPrettyJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(levelToZapLevel(level))))
	} else {
		encoderConfig = productionEncoderConfig
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(encoder, fileInfo, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel // Logs with level LOWER than ERROR
	})))

	// Core for error.json: logs at ERROR level and above
	cores = append(cores, zapcore.NewCore(encoder, fileError, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel // Logs with level ERROR and HIGHER
	})))

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar()
}

func NewLoggerFromEnv() *zap.SugaredLogger {
	level := os.Getenv("LOG_LEVEL")
	development := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "development"
	return NewLogger(level, development)
}

func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, entities.LoggerKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(entities.LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

func FromContextWithName(ctx context.Context, name string) *zap.SugaredLogger {
	return FromContext(ctx).Named(name).With("request_id", ctx.Value(entities.RequestId).(string))
}

const (
	timestamp = "timestamp"
	severity  = "severity"
	logger    = "logger"
	caller    = "caller"
	message   = "message"
	function  = "function"

	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "EMERGENCY"
)

var productionEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        timestamp,
	LevelKey:       severity,
	NameKey:        logger,
	CallerKey:      caller,
	MessageKey:     message,
	LineEnding:     zapcore.DefaultLineEnding,
	FunctionKey:    function,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     timeEncoder(),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var developmentEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "T",
	LevelKey:       "L",
	NameKey:        "N",
	CallerKey:      "C",
	FunctionKey:    "F",
	MessageKey:     "M",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func levelToZapLevel(s string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelCritical:
		return zapcore.DPanicLevel
	case levelAlert:
		return zapcore.PanicLevel
	case levelEmergency:
		return zapcore.FatalLevel
	}

	return zapcore.WarnLevel
}

func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(levelDebug)
		case zapcore.InfoLevel:
			enc.AppendString(levelInfo)
		case zapcore.WarnLevel:
			enc.AppendString(levelWarning)
		case zapcore.ErrorLevel:
			enc.AppendString(levelError)
		case zapcore.DPanicLevel:
			enc.AppendString(levelCritical)
		case zapcore.PanicLevel:
			enc.AppendString(levelAlert)
		case zapcore.FatalLevel:
			enc.AppendString(levelEmergency)
		}
	}
}

func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}
