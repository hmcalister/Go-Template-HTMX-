package main

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Perform logging of requests using Zerolog's global logger.
//
// Should be placed early if not first in middleware stack.
//
// Inspiration taken from https://github.com/ironstar-io/chizerolog/blob/master/main.go
func zerologLogger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		requestTimeReceived := time.Now()

		next.ServeHTTP(wrappedWriter, r)

		requestTimeResolved := time.Now()
		log.Info().
			Str("URL", r.URL.Path).
			Str("Protocol", r.Proto).
			Str("RemoteIP", r.RemoteAddr).
			Int("Status", wrappedWriter.Status()).
			Str("UserAgent", r.Header.Get("User-Agent")).
			Float32("Latency_ms", float32(requestTimeResolved.Sub(requestTimeReceived).Nanoseconds()/1_000_000.0)).
			Int64("BytesReceived", r.ContentLength).
			Int("BytesSent", wrappedWriter.BytesWritten()).
			Send()
	})
}

// Catches any unintended panics in the server and recovers, logging the result for future reference
//
// This middleware should be placed first (or at least early) such that all errors are recovered from.
func recoverWithInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer so that any panics are caught when logger is resolved
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Interface("RecoverInformation", rec).
					Bytes("DebugStack", debug.Stack()).
					Msg("ErrorRecoveryDuringZerologMiddleware")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
