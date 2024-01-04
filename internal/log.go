package internal

import (
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

func InitLogger() zerolog.Logger {
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	return zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
}

func SetDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
