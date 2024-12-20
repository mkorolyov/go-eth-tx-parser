package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/mkorolyov/go-eth-tx-parser/internal/ethereum"
	"github.com/mkorolyov/go-eth-tx-parser/internal/poller"
	"github.com/mkorolyov/go-eth-tx-parser/internal/server"
	"github.com/mkorolyov/go-eth-tx-parser/internal/storage"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ethClient := ethereum.NewJsonRPCClient(ethereum.WithHTTPClient(&http.Client{}), ethereum.WithLog(logger))
	inMemStorage := storage.NewInMemoryStorage()
	transactionPoller := poller.NewTransactionPoller(inMemStorage, inMemStorage, inMemStorage, ethClient, logger)

	// Start polling for new transactions
	go transactionPoller.Start(ctx)

	httpServer := server.NewNaiveHTTPServer(inMemStorage, logger)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Info("server stopped", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server...")
	// ctx already closed but it is not a problem here
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("shutdown server", "error", err)
	}
	logger.Info("exiting...")
}
