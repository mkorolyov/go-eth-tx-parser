package parser

import (
	"context"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/ethereum"
)

type Parser interface {
	// last parsed block
	// In the task methods where defined without error in response
	// while here error was added to propogate possible errors underthehood
	GetCurrentBlock(ctx context.Context) (int, error)
	// add address to observer
	Subscribe(ctx context.Context, address string) error
	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error)
}
