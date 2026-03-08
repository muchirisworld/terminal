package middleware

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muchirisworld/terminal/internal/logger"
)

// Logger is a middleware that logs the start and end of each request, along with some useful information about it.
func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Initialize Wide Event
			wideEvent := logger.NewWideEvent()

			// Generate request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", requestID)

			// Add base environment/request info
			wideEvent.Add("method", r.Method)
			wideEvent.Add("path", r.URL.Path)
			wideEvent.Add("request_id", requestID)

			// Context with wide event
			ctx := logger.WithEvent(r.Context(), wideEvent)
			r = r.WithContext(ctx)

			ww := &responseWriter{w, http.StatusOK}

			defer func() {
				duration := time.Since(start)
				wideEvent.Add("duration_ms", duration.Milliseconds())
				wideEvent.Add("status_code", ww.status)

				// Determine log level based on status code or error presence
				data := wideEvent.GetAll()
				args := make([]any, 0, len(data)*2)
				hasError := false
				for k, v := range data {
					args = append(args, k, v)
					if k == "error" {
						hasError = true
					}
				}

				if ww.status >= 500 || hasError {
					wideEvent.Add("outcome", "error")
					args = append(args, "outcome", "error")
					log.Error("request_completed", args...)
				} else {
					wideEvent.Add("outcome", "success")
					args = append(args, "outcome", "success")
					log.Info("request_completed", args...)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijack not supported")
	}
	return h.Hijack()
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return fmt.Errorf("push not supported")
}
