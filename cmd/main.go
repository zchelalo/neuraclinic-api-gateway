package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zchelalo/neuraclinic-api-gateway/pkg/bootstrap"
	"go.uber.org/zap"
)

func main() {
	cfg, err := bootstrap.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	logger := bootstrap.GetLogger()
	defer bootstrap.SyncLogger()

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := bootstrap.InitApp(rootCtx, logger, cfg)
	if err != nil {
		logger.Fatal("cannot initialize application", zap.Error(err))
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server starting", zap.Int("port", cfg.Port))
		errCh <- app.Server.Start()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		logger.Info("signal received, shutting down", zap.String("signal", sig.String()))
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", zap.Error(err))
		}
	}

	cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = app.Cleanup(ctx)
	logger.Info("shutdown complete")
}
