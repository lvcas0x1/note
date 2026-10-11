package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"notes/internal/api"
	"notes/internal/config"
	"notes/internal/httpx"
	"notes/internal/logx"
	"notes/internal/webui"
)

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the current build over HTTP (NOTES_LISTEN, default 127.0.0.1:8080)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(nil)
			if err != nil {
				return err
			}
			log := logx.New(os.Stderr, slog.LevelInfo)
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			ln, err := net.Listen("tcp", cfg.Listen)
			if err != nil {
				return err
			}
			srv, closeFn, err := newHTTPServer(cfg, log)
			if err != nil {
				_ = ln.Close()
				return err
			}
			defer closeFn()
			log.Info("listening", "addr", ln.Addr().String())
			errc := make(chan error, 1)
			go func() { errc <- srv.Serve(ln) }()
			select {
			case err := <-errc:
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			case <-ctx.Done():
				shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				return srv.Shutdown(shutdown)
			}
		},
	}
}

// newHTTPServer wires the engine: security headers, the API, assets, health, and the web app.
func newHTTPServer(cfg config.Server, log *slog.Logger) (*http.Server, func(), error) {
	current, err := httpx.NewCurrent(cfg.DataDir, log)
	if err != nil {
		return nil, nil, err
	}
	e, err := httpx.NewEngine(log, cfg.TrustedProxies...)
	if err != nil {
		_ = current.Close()
		return nil, nil, err
	}
	e.Use(httpx.SecurityHeaders())
	limiter := httpx.NewRateLimiter(cfg.SearchRate)
	api.Register(e, current, limiter.HumaMiddleware)
	httpx.Mount(e, current, webui.FS())
	return httpx.NewServer(cfg.Listen, e), func() { _ = current.Close() }, nil
}
