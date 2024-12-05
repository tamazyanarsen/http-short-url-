package logger

import "go.uber.org/zap"

var Logger zap.SugaredLogger

func InitLogger() {
	if logger, err := zap.NewDevelopment(); err == nil {
		Logger = *logger.Sugar()
	}
}
