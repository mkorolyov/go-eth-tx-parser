package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/etherium"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/http"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/logger"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/parser"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/storage"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	ethClient := etherium.NewJsonRPCClient()
	inMemStorage := storage.NewInMemoryStorage()
	log := logger.DefaultLogger
	observer := parser.NewObserver(inMemStorage, inMemStorage, inMemStorage, ethClient, log)

	// Start polling for new transactions
	go observer.StartPooling(ctx)

	server := http.NewNaiveHTTPServer(observer, log)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Infof("server stopped: %v", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down server...")
	server.Shutdown(ctx)
	log.Info("exiting...")
}
