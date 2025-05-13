package timeout

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
	"github.com/vovancho/lingua-cat-go/pkg/response"
)

func TimeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)
		tw := &timeoutWriter{ResponseWriter: w}

		done := make(chan struct{})
		go func() {
			next.ServeHTTP(tw, r)
			close(done)
		}()

		select {
		case <-done:
			tw.FlushToOriginal()
		case <-ctx.Done():
			tw.mu.Lock()
			tw.timedOut = true
			tw.mu.Unlock()

			slog.Info("TimeoutMiddleware: context deadline exceeded2")
			appError := _internalError.NewAppError(http.StatusGatewayTimeout, "Request timed out", ctx.Err())
			response.HandleError(w, appError, nil)
		}
	})
}

type timeoutWriter struct {
	http.ResponseWriter
	wroteHeader bool
	mu          sync.Mutex
	buf         bytes.Buffer
	statusCode  int
	timedOut    bool
}

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.statusCode = statusCode
	tw.wroteHeader = true
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return 0, nil
	}
	if !tw.wroteHeader {
		tw.statusCode = http.StatusOK
		tw.wroteHeader = true
	}
	return tw.buf.Write(b)
}

func (tw *timeoutWriter) FlushToOriginal() {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if !tw.timedOut {
		if !tw.wroteHeader {
			tw.ResponseWriter.WriteHeader(http.StatusOK)
		} else {
			tw.ResponseWriter.WriteHeader(tw.statusCode)
		}
		tw.ResponseWriter.Write(tw.buf.Bytes())
	}
}
