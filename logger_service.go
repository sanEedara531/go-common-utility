package common

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//ZapLoggerObj Global LOgger Object
var ZapLoggerObj *zap.Logger

//ZapLogger logger init
func ZapLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level.SetLevel(zapcore.InfoLevel)
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	//cfg.OutputPaths = []string{"/tmp/logs/dispatcher.log"}
	logger, err := cfg.Build() // panic()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
