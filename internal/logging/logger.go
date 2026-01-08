package logging

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitDefault() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed init logger: %v", err)
	}
	logger = l.Sugar()
}

func Get() *zap.SugaredLogger {
	if logger == nil {
		InitDefault()
	}
	return logger
}
