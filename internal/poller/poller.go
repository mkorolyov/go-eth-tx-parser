package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mkorolyov/go-eth-tx-parser/internal/ethereum"
)

type TransactionsStorage interface {
	SaveTransaction(ctx context.Context, address string, tx ethereum.Transaction) error
	GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error)
}

type AddressesStorage interface {
	Subscribe(ctx context.Context, address string) error
	IsSubscribed(ctx context.Context, address string) (bool, error)
}

type BlocksStorage interface {
	SetCurrentBlock(ctx context.Context, block int) error
	GetCurrentBlock(ctx context.Context) (int, error)
}

type EthClient interface {
	GetBlockNumber(ctx context.Context) (int, error)
	GetBlockByNumber(ctx context.Context, blockNumber int) (ethereum.EthereumBlock, error)
}

func NewTransactionPoller(
	transactionsStorage TransactionsStorage,
	addressesStorage AddressesStorage,
	blocksStorage BlocksStorage,
	ethClient EthClient,
	log *slog.Logger,
) TransactionPoller {
	return TransactionPoller{
		transactionsStorage: transactionsStorage,
		addressesStorage:    addressesStorage,
		blocksStorage:       blocksStorage,
		ethClient:           ethClient,
		log:                 log,
	}
}

type TransactionPoller struct {
	transactionsStorage TransactionsStorage
	addressesStorage    AddressesStorage
	blocksStorage       BlocksStorage
	ethClient           EthClient
	log                 *slog.Logger
}

func (p TransactionPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 12)
	// for go versions prior 1.23
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.log.Info("stopping observer")
			return
		default:
			p.loadNewTransactions(ctx)
			<-ticker.C
		}
	}
}

func (p TransactionPoller) loadNewTransactions(ctx context.Context) {
	latestBlock, err := p.ethClient.GetBlockNumber(ctx)
	if err != nil {
		p.log.Error("error fetching latest block", "error", err)
		return
	}

	currentBlock, err := p.blocksStorage.GetCurrentBlock(ctx)
	if err != nil {
		p.log.Error("failed to load last processed block", "error", err)
		return
	}

	if currentBlock == 0 {
		currentBlock = latestBlock - 1
	}

	// Process new blocks
	for i := currentBlock + 1; i <= latestBlock; i++ {
		block, err := p.ethClient.GetBlockByNumber(ctx, i)
		if err != nil {
			p.log.Error("failed to load block", "block", fmt.Sprintf("%x", i), "error", err)
			return
		}

		p.log.Info("processing block with new transactions", "block", fmt.Sprintf("%x", i), "transactions_count", len(block.Transactions))

		// Process transactions in the block
		for _, tx := range block.Transactions {
			p.saveTxForAddress(ctx, tx, tx.To)
			p.saveTxForAddress(ctx, tx, tx.From)
		}

		// Update the current block
		if err := p.blocksStorage.SetCurrentBlock(ctx, i); err != nil {
			p.log.Error("failed to set current processed block", "block", fmt.Sprintf("%x", i), "error", err)
		}
	}
}

func (p TransactionPoller) saveTxForAddress(ctx context.Context, tx ethereum.Transaction, address string) {
	subscribed, err := p.addressesStorage.IsSubscribed(ctx, address)
	if err != nil {
		p.log.Error("failed to check if address is subscribed", "address", address, "error", err)
		return
	}

	if !subscribed {
		return
	}

	if err := p.transactionsStorage.SaveTransaction(ctx, address, tx); err != nil {
		p.log.Error("failed to save transaction", "transaction_hash", tx.Hash, "address", address, "error", err)
	}

	p.log.Debug("transaction saved for address", "address", address, "transaction_hash", tx.Hash)
}
