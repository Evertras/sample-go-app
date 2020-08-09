package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/evertras/sample-go-app/internal/server"
)

const address = "0.0.0.0:8088"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		fmt.Println("Failed to create logger:", err)
		os.Exit(1)
	}

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
