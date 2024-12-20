package http

import (
	"encoding/json"
	"net/http"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/logger"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/parser"
)

func NewNaiveHTTPServer(observer parser.Parser, log logger.Logger) http.Server {
	http.HandleFunc("POST /address/{address}/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.PathValue("address")
		if err := observer.Subscribe(r.Context(), address); err != nil {
			log.Errorf("failed to subscribe to address %s: %v", address, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		txs, err := observer.GetTransactions(r.Context(), address)
		if err != nil {
			log.Errorf("failed to get transactions for address %s: %v", address, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(txs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err := json.NewEncoder(w).Encode(txs); err != nil {
			log.Errorf("failed to encode transactions for address %s: %v", address, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("GET /current_block", func(w http.ResponseWriter, r *http.Request) {
		block, err := observer.GetCurrentBlock(r.Context())
		if err != nil {
			log.Errorf("failed to get current block: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]int{"current_block": block}); err != nil {
			log.Errorf("failed to encode current block: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	return http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
}
