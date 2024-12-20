package storage

import (
	"context"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/ethereum"
)

type Transactions interface {
	SaveTransaction(ctx context.Context, address string, tx ethereum.Transaction) error
	GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error)
}

type Addresses interface {
	Subscribe(ctx context.Context, address string) error
	IsSubscribed(ctx context.Context, address string) (bool, error)
}

type Blocks interface {
	SetCurrentBlock(ctx context.Context, block int) error
	GetCurrentBlock(ctx context.Context) (int, error)
}
