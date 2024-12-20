package storage

import (
	"context"
	"sync"

	"github.com/mkorolyov/go-eth-tx-parser/internal/ethereum"
)

// InMemoryStorage is a thread-safe in-memory storage for transactions
type InMemoryStorage struct {
	mu sync.RWMutex
	// address -> transactions
	transactions        map[string][]ethereum.Transaction
	currentBlock        int
	subscribedAddresses map[string]struct{}
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		transactions:        make(map[string][]ethereum.Transaction),
		subscribedAddresses: make(map[string]struct{}),
	}
}

// SaveTransaction AddTransaction stores a transaction for an address
func (s *InMemoryStorage) SaveTransaction(_ context.Context, address string, tx ethereum.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[address] = append(s.transactions[address], tx)
	return nil
}

// GetTransactions fetches transactions for a given address. Paging is not supported for simplicity
func (s *InMemoryStorage) GetTransactions(_ context.Context, address string) ([]ethereum.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transactions[address], nil
}

// Subscribe adds an address to be observed
func (s *InMemoryStorage) Subscribe(_ context.Context, address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribedAddresses[address] = struct{}{}
	return nil
}

// IsSubscribed checks if an address is being observed
func (s *InMemoryStorage) IsSubscribed(_ context.Context, address string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.subscribedAddresses[address]
	return ok, nil
}

// SetCurrentBlock updates the current block
func (s *InMemoryStorage) SetCurrentBlock(_ context.Context, block int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentBlock = block
	return nil
}

// GetCurrentBlock retrieves the current block
func (s *InMemoryStorage) GetCurrentBlock(_ context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentBlock, nil
}
