package main

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/evertras/sample-go-app/internal/server"
)

const address = "0.0.0.0:8088"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapLevel,
	))
	defer logger.Sync()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		logger.Warn("Received interrupt")
		cancel()
	}()

	s := server.New(logger, address)

	logger.Info(
		"Listening",
		zap.String("address", address),
	)

	if err := s.ListenAndServe(ctx); err != nil {
		logger.Error("Failed to listen", zap.Error(err))
	}

	logger.Info("Exited gracefully")
}
