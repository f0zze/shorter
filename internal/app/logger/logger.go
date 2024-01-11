package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(fileName string) zerolog.Logger {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	logger := zerolog.New(file).With().Timestamp().Logger()

	return logger
}
