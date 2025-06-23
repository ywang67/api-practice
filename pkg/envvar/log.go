package envvar

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// at most two logs a minute, plus one in a thousand after that.
// this sampler is probably overkill, but it's frustrating for our logs to be filled with noise and possibly expensive.
var sampler = &zerolog.BurstSampler{Burst: 3, Period: time.Minute, NextSampler: &zerolog.BasicSampler{N: 1000}}

// this mutex is only really necessary to stop a race during the tests, since we will
// only modify 'logger' after creation there.
var logMux sync.Mutex
var logger zerolog.Logger

func init() {
	if b, _ := lookupBool("LOCAL"); b {
		logger = zerolog.New(&zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel).With().Caller().Logger()
	} else {
		logger = log.Logger.With().Str("source", "genome/pkg/envvar/log.go").Logger().Sample(sampler)
	}
}

func logPanic(key string) *zerolog.Event {
	val, err := Lookup(key)
	if err != nil {
		return logger.Panic().Str("key", key)
	}
	return logger.Panic().Str("key", key).Str("val", val)
}

type logEvent struct {
	// err                   error
	key, caller, fallback string
}

// func caller(skip int) string {
// 	_, f, _, _ := runtime.Caller(skip + 1)
// 	return f
// }

var logQueue = make(chan logEvent, 32)

func init() { go batchLogs(logQueue, time.Minute) }

// we don't want to spam the logs, so we batch them together and send them every minute (if there are any to send).
func batchLogs(queue <-chan logEvent, d time.Duration) {
	t := time.NewTicker(d)
	var events []logEvent
	for {
		select {
		case <-t.C:
			if len(events) == 0 {
				continue
			}
			logMux.Lock()
			info := logger.Info()
			for _, e := range events {

				info = info.Dict(e.key, zerolog.Dict().Str("fallback", e.fallback).Str("caller", e.caller))
			}
			events = events[:0]
			info.Msg("missing environment variables: falling back to default values")
			logMux.Unlock()
		case e := <-queue:
			events = append(events, e)
		}
	}
}
