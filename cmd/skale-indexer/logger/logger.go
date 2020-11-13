package logger

import (
	"strings"

	"github.com/bearcherian/rollzap"
	"github.com/rollbar/rollbar-go"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log Logger
)

type RollbarConfig struct {
	RollbarAccessToken string
	RollbarServerRoot  string
	AppEnv             string
	Version            string
	ChainIDs           []string
}

type Logger struct {
	*zap.Logger
}

func GetLogger() *zap.Logger {
	return Log.Logger
}

func Init(encoding, logLevel string, logOutputs []string, rollbarConfig *RollbarConfig) error {
	logConfig := zap.Config{
		OutputPaths: logOutputs,
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(getLevel(logLevel)),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	}

	log, err := logConfig.Build()
	if err != nil {
		return err
	}

	if rollbarConfig != nil && rollbarConfig.RollbarAccessToken != "" {
		rollbar.SetToken(rollbarConfig.RollbarAccessToken)
		rollbar.SetEnvironment(rollbarConfig.AppEnv)
		rollbar.SetServerRoot(rollbarConfig.RollbarServerRoot)
		rollbar.SetCodeVersion(rollbarConfig.Version)
		rollbar.SetCustom(map[string]interface{}{"chainIDs": rollbarConfig.ChainIDs})

		log.Info("Enabling rollbar integration", zap.String("env", rollbarConfig.AppEnv), zap.String("server_root", rollbarConfig.RollbarServerRoot), zap.String("access_token", rollbarConfig.RollbarAccessToken[0:5]))
		// create a new core that sends zapcore.ErrorLevel and above messages to Rollbar
		rollbarCore := rollzap.NewRollbarCore(zapcore.ErrorLevel)

		// Wrap a NewTee to send log messages to both your main logger and to rollbar
		log = log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, rollbarCore)
		}))
	}
	Log.Logger = log

	return nil
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Sync() {
	Log.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Error(err error, fields ...zap.Field) {
	Log.Error(err.Error(), fields...)
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	case "error":
		return zap.ErrorLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.InfoLevel
	}
}
