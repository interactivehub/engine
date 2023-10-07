package logger

import (
	"os"

	"go.uber.org/zap"
)

func Init() *zap.Logger {
	lgr := zap.Must(zap.NewProduction())

	if os.Getenv("APP_ENV") == "development" {
		lgr = zap.Must(zap.NewDevelopment())
	}

	return lgr
}
