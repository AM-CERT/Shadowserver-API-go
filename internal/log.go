package internal

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func InitLogger() zerolog.Logger {
	return zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
}

func SetDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
