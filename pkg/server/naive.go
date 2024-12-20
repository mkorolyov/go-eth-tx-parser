package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mkorolyov/go-eth-tx-parser/internal/ethereum"
)

type Parser interface {
	// GetCurrentBlock returns last parsed block
	// In the task methods where defined without error in response
	// while here error was added to propogate possible errors underthehood
	GetCurrentBlock(ctx context.Context) (int, error)
	// Subscribe add address to observer
	Subscribe(ctx context.Context, address string) error
	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error)
}

func NewNaiveHTTPServer(parser Parser, log *slog.Logger) *http.Server {
	serverMux := http.NewServeMux()

	serverMux.HandleFunc("POST /address/{address}/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := strings.ToLower(r.PathValue("address"))

		log.Info("subscribing to address", "address", address)

		if err := parser.Subscribe(r.Context(), address); err != nil {
			log.Error("failed to subscribe to address", "address", address, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	//naive implementation without paging support
	serverMux.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		txs, err := parser.GetTransactions(r.Context(), address)
		if err != nil {
			log.Error("failed to get transactions for address", "address", address, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(txs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err := json.NewEncoder(w).Encode(txs); err != nil {
			log.Error("failed to encode transactions for address", "address", address, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	serverMux.HandleFunc("GET /current_block", func(w http.ResponseWriter, r *http.Request) {
		block, err := parser.GetCurrentBlock(r.Context())
		if err != nil {
			log.Error("failed to get current block", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]int{"current_block": block}); err != nil {
			log.Error("failed to encode current block", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	return &http.Server{Addr: ":8080", Handler: serverMux}
}
