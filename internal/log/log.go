package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var ConsoleLogger *zap.Logger
var FileLogger *zap.Logger

func InitLogger(logLevel string) {
	fileCore, consoleCore, combinedCore := buildZapCore(getZapLogLevel(logLevel))
	FileLogger = zap.New(fileCore, zap.AddCaller())
	ConsoleLogger = zap.New(consoleCore, zap.AddCaller())
	Logger = zap.New(combinedCore, zap.AddCaller())
}

func buildZapCore(logLevel zapcore.Level) (fileCore zapcore.Core, consoleCore zapcore.Core, combinedCore zapcore.Core) {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
	consoleWriter := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	fileCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), fileWriter, logLevel)
	consoleCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel)
	combinedCore = zapcore.NewTee(fileCore, consoleCore)
	return
}

func getZapLogLevel(logLevel string) zapcore.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "info":
		return zap.InfoLevel
	default:
		return zap.InfoLevel
	}
}
