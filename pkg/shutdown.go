package pkg

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Для корректного завершения работы HTTP-сервера
func Graceful(srv *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered in shutdown goroutine", "err", r)
			}
		}()

		<-quit
		slog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown failed", "err", err)
		} else {
			slog.Info("server shutdown completed")
		}
	}()
}
