package parser

import (
	"context"
	"time"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/ethereum"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/logger"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/storage"
)

func NewObserver(
	transactionsStorage storage.Transactions,
	addressesStorage storage.Addresses,
	blocksStorage storage.Blocks,
	ethClient ethereum.Client,
	log logger.Logger,
) Observer {
	return Observer{
		transactionsStorage: transactionsStorage,
		addressesStorage:    addressesStorage,
		blocksStorage:       blocksStorage,
		ethClient:           ethClient,
		log:                 log,
	}
}

type Observer struct {
	transactionsStorage storage.Transactions
	addressesStorage    storage.Addresses
	blocksStorage       storage.Blocks
	ethClient           ethereum.Client
	log                 logger.Logger
}

func (p Observer) GetCurrentBlock(ctx context.Context) (int, error) {
	return p.blocksStorage.GetCurrentBlock(ctx)
}

func (p Observer) Subscribe(ctx context.Context, address string) error {
	return p.addressesStorage.Subscribe(ctx, address)
}

func (p Observer) GetTransactions(ctx context.Context, address string) ([]ethereum.Transaction, error) {
	return p.transactionsStorage.GetTransactions(ctx, address)
}

func (p Observer) StartPooling(ctx context.Context) {
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

func (p Observer) loadNewTransactions(ctx context.Context) {
	latestBlock, err := p.ethClient.GetBlockNumber(ctx)
	if err != nil {
		p.log.Errorf("error fetching latest block: %v", err)
		return
	}

	currentBlock, err := p.blocksStorage.GetCurrentBlock(ctx)
	if err != nil {
		p.log.Errorf("failed to load last processed block: %v", err)
		return
	}

	if currentBlock == 0 {
		currentBlock = latestBlock - 1
	}

	// Process new blocks
	for i := currentBlock + 1; i <= latestBlock; i++ {
		block, err := p.ethClient.GetBlockByNumber(ctx, i)
		if err != nil {
			p.log.Errorf("failed to load block %x : %v", i, err)
			return
		}

		p.log.Infof("processing block %x with %d new transactions", i, len(block.Transactions))

		// Process transactions in the block
		for _, tx := range block.Transactions {
			p.saveTxForAddress(ctx, tx, tx.To)
			p.saveTxForAddress(ctx, tx, tx.From)
		}

		// Update the current block
		if err := p.blocksStorage.SetCurrentBlock(ctx, i); err != nil {
			p.log.Errorf("failed to set current processed block to %x: %v", i, err)
		}
	}
}

func (p Observer) saveTxForAddress(ctx context.Context, tx ethereum.Transaction, address string) {
	subscribed, err := p.addressesStorage.IsSubscribed(ctx, address)
	if err != nil {
		p.log.Errorf("failed to check if address %s is subscribed: %v", address, err)
		return
	}

	if !subscribed {
		return
	}

	if err := p.transactionsStorage.SaveTransaction(ctx, address, tx); err != nil {
		p.log.Errorf("failed to save transaction with hash %s for address %s: %v", tx.Hash, address, err)
	}
}
