package parser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mkorolyov/go-eth-tx-parser/pkg/etherium"
	"github.com/mkorolyov/go-eth-tx-parser/pkg/logger"
)

func TestStartPooling(t *testing.T) {
	mockEthClient := &MockEthClient{}
	mockBlocksStorage := &MockBlocksStorage{}
	mockAddressesStorage := &MockAddressesStorage{}
	mockTransactionsStorage := &MockTransactionsStorage{}
	mockLogger := logger.DefaultLogger

	observer := Observer{
		ethClient:           mockEthClient,
		blocksStorage:       mockBlocksStorage,
		addressesStorage:    mockAddressesStorage,
		transactionsStorage: mockTransactionsStorage,
		log:                 mockLogger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("context done", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 0, nil
		}
		mockEthClient.GetBlockByNumberFunc = func(ctx context.Context, number int) (etherium.EthereumBlock, error) {
			return etherium.EthereumBlock{}, nil
		}
		mockBlocksStorage.SetCurrentBlockFunc = func(ctx context.Context, number int) error {
			return nil
		}
		go func() {
			cancel()
		}()
		observer.StartPooling(ctx)
	})

	t.Run("error fetching latest block", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 0, errors.New("error fetching latest block")
		}
		go func() {
			time.Sleep(time.Second * 1)
			cancel()
		}()
		observer.StartPooling(ctx)
	})

	t.Run("successful pooling", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 100, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 99, nil
		}
		mockEthClient.GetBlockByNumberFunc = func(ctx context.Context, number int) (etherium.EthereumBlock, error) {
			return etherium.EthereumBlock{Transactions: []etherium.Transaction{{Hash: "0x123"}}}, nil
		}
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return true, nil
		}
		mockTransactionsStorage.SaveTransactionFunc = func(ctx context.Context, address string, tx etherium.Transaction) error {
			return nil
		}
		mockBlocksStorage.SetCurrentBlockFunc = func(ctx context.Context, number int) error {
			return nil
		}
		go func() {
			time.Sleep(time.Second * 1)
			cancel()
		}()
		observer.StartPooling(ctx)
	})
}

func TestLoadNewTransactions(t *testing.T) {
	mockEthClient := &MockEthClient{}
	mockBlocksStorage := &MockBlocksStorage{}
	mockAddressesStorage := &MockAddressesStorage{}
	mockTransactionsStorage := &MockTransactionsStorage{}
	mockLogger := logger.DefaultLogger

	observer := Observer{
		ethClient:           mockEthClient,
		blocksStorage:       mockBlocksStorage,
		addressesStorage:    mockAddressesStorage,
		transactionsStorage: mockTransactionsStorage,
		log:                 mockLogger,
	}

	ctx := context.Background()

	t.Run("error fetching latest block", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 0, errors.New("error fetching latest block")
		}

		observer.loadNewTransactions(ctx)
	})

	t.Run("error loading last processed block", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 100, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 0, errors.New("failed to load last processed block")
		}

		observer.loadNewTransactions(ctx)
	})

	t.Run("error loading block", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 100, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 99, nil
		}
		mockEthClient.GetBlockByNumberFunc = func(ctx context.Context, number int) (etherium.EthereumBlock, error) {
			return etherium.EthereumBlock{}, errors.New("failed to load block")
		}

		observer.loadNewTransactions(ctx)
	})

	t.Run("error setting current block", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 100, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 99, nil
		}
		mockEthClient.GetBlockByNumberFunc = func(ctx context.Context, number int) (etherium.EthereumBlock, error) {
			return etherium.EthereumBlock{Transactions: []etherium.Transaction{{Hash: "0x123"}}}, nil
		}
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return true, nil
		}
		mockTransactionsStorage.SaveTransactionFunc = func(ctx context.Context, address string, tx etherium.Transaction) error {
			return nil
		}
		mockBlocksStorage.SetCurrentBlockFunc = func(ctx context.Context, number int) error {
			return errors.New("failed to set current block")
		}

		observer.loadNewTransactions(ctx)
	})

	t.Run("successful load new transactions", func(t *testing.T) {
		mockEthClient.GetBlockNumberFunc = func(ctx context.Context) (int, error) {
			return 100, nil
		}
		mockBlocksStorage.GetCurrentBlockFunc = func(ctx context.Context) (int, error) {
			return 99, nil
		}
		mockEthClient.GetBlockByNumberFunc = func(ctx context.Context, number int) (etherium.EthereumBlock, error) {
			return etherium.EthereumBlock{Transactions: []etherium.Transaction{{Hash: "0x123"}}}, nil
		}
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return true, nil
		}
		mockTransactionsStorage.SaveTransactionFunc = func(ctx context.Context, address string, tx etherium.Transaction) error {
			return nil
		}
		mockBlocksStorage.SetCurrentBlockFunc = func(ctx context.Context, number int) error {
			return nil
		}

		observer.loadNewTransactions(ctx)
	})
}

func TestSaveTxForAddress(t *testing.T) {
	mockTransactionsStorage := &MockTransactionsStorage{}
	mockAddressesStorage := &MockAddressesStorage{}
	mockLogger := logger.DefaultLogger

	observer := Observer{
		transactionsStorage: mockTransactionsStorage,
		addressesStorage:    mockAddressesStorage,
		log:                 mockLogger,
	}

	ctx := context.Background()
	tx := etherium.Transaction{Hash: "0x123"}

	t.Run("address not subscribed", func(t *testing.T) {
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return false, nil
		}

		observer.saveTxForAddress(ctx, tx, "0xabc")
	})

	t.Run("error checking subscription", func(t *testing.T) {
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return false, errors.New("subscription check error")
		}

		observer.saveTxForAddress(ctx, tx, "0xabc")
	})

	t.Run("error saving transaction", func(t *testing.T) {
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return true, nil
		}
		mockTransactionsStorage.SaveTransactionFunc = func(ctx context.Context, address string, tx etherium.Transaction) error {
			return errors.New("save transaction error")
		}

		observer.saveTxForAddress(ctx, tx, "0xabc")
	})

	t.Run("successful save transaction", func(t *testing.T) {
		mockAddressesStorage.IsSubscribedFunc = func(ctx context.Context, address string) (bool, error) {
			return true, nil
		}
		mockTransactionsStorage.SaveTransactionFunc = func(ctx context.Context, address string, tx etherium.Transaction) error {
			return nil
		}

		observer.saveTxForAddress(ctx, tx, "0xabc")
	})
}

type MockTransactionsStorage struct {
	SaveTransactionFunc func(ctx context.Context, address string, tx etherium.Transaction) error
	GetTransactionsFunc func(ctx context.Context, address string) ([]etherium.Transaction, error)
}

func (m *MockTransactionsStorage) SaveTransaction(ctx context.Context, address string, tx etherium.Transaction) error {
	return m.SaveTransactionFunc(ctx, address, tx)
}

func (m *MockTransactionsStorage) GetTransactions(ctx context.Context, address string) ([]etherium.Transaction, error) {
	return m.GetTransactionsFunc(ctx, address)
}

type MockAddressesStorage struct {
	IsSubscribedFunc func(ctx context.Context, address string) (bool, error)
	SubscribeFunc    func(ctx context.Context, address string) error
}

func (m *MockAddressesStorage) IsSubscribed(ctx context.Context, address string) (bool, error) {
	return m.IsSubscribedFunc(ctx, address)
}

func (m *MockAddressesStorage) Subscribe(ctx context.Context, address string) error {
	return m.SubscribeFunc(ctx, address)
}

type MockEthClient struct {
	GetBlockNumberFunc   func(ctx context.Context) (int, error)
	GetBlockByNumberFunc func(ctx context.Context, number int) (etherium.EthereumBlock, error)
}

func (m *MockEthClient) GetBlockNumber(ctx context.Context) (int, error) {
	return m.GetBlockNumberFunc(ctx)
}

func (m *MockEthClient) GetBlockByNumber(ctx context.Context, number int) (etherium.EthereumBlock, error) {
	return m.GetBlockByNumberFunc(ctx, number)
}

type MockBlocksStorage struct {
	GetCurrentBlockFunc func(ctx context.Context) (int, error)
	SetCurrentBlockFunc func(ctx context.Context, number int) error
}

func (m *MockBlocksStorage) GetCurrentBlock(ctx context.Context) (int, error) {
	return m.GetCurrentBlockFunc(ctx)
}

func (m *MockBlocksStorage) SetCurrentBlock(ctx context.Context, number int) error {
	return m.SetCurrentBlockFunc(ctx, number)
}
