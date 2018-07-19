package zlog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
