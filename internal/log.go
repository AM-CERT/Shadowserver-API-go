package internal

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Logger zerolog.Logger

func InitLogger() {
	Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
}

func SetDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
