package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		headers := make([]string, 0, len(r.Header))
		for k, v := range r.Header {
			headers = append(headers, k+":"+strings.Join(v, ", "))
		}

		rw := &responseRecorder{
			ResponseWriter: w, status: http.StatusOK, size: 0}

		next.ServeHTTP(rw, r)
		m.logger.Debug("запрос",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", rw.status),
			zap.String("headers", strings.Join(headers, ",")),
			zap.Int("size", rw.size),
			zap.Duration("duration", time.Since(start)))

	})
}

func (m Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				m.logger.Error("возникла паника",
					zap.Any("rec", rec),
					zap.ByteString("stack", debug.Stack()),
					zap.String("method", r.Method),
					zap.String("url", r.URL.Path),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
