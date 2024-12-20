package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
)

// JsonRPCClient interacts with the Ethereum JSON-RPC endpoint
type JsonRPCClient struct {
	endpoint string
	http     *http.Client
	log      *slog.Logger
}

type Option func(*JsonRPCClient)

func WithHTTPClient(http *http.Client) Option {
	return func(c *JsonRPCClient) {
		c.http = http
	}
}

func WithLog(log *slog.Logger) Option {
	return func(c *JsonRPCClient) {
		c.log = log
	}
}

const defaultEndpoint = "https://ethereum-rpc.publicnode.com"

func NewJsonRPCClient(options ...Option) JsonRPCClient {
	c := JsonRPCClient{endpoint: defaultEndpoint, log: slog.New(slog.NewJSONHandler(os.Stdout, nil))}

	for _, option := range options {
		option(&c)
	}

	return c
}

// GetBlockNumber fetches the latest block number
func (c JsonRPCClient) GetBlockNumber(ctx context.Context) (int, error) {
	requestBody := EthereumJSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      rand.Int(),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make eht_blockNumber request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close http response body", "error", err)
		}
	}()

	var rpcResponse EthereumJSONRPCResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if rpcResponse.Result == "" {
		return 0, errors.New("got empty block number")
	}

	// Convert hex to integer
	var blockNumber int
	_, err = fmt.Sscanf(rpcResponse.Result, "0x%x", &blockNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number %s : %w", rpcResponse.Result, err)
	}
	return blockNumber, nil
}

const returnFullTx = true

// GetBlockByNumber fetches a block and its transactions by block number
func (c JsonRPCClient) GetBlockByNumber(ctx context.Context, blockNumber int) (EthereumBlock, error) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)

	requestBody := EthereumJSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumberHex, returnFullTx},
		ID:      rand.Int(),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return EthereumBlock{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return EthereumBlock{}, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return EthereumBlock{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close http response body", "error", err)
		}
	}()

	var rpcResponse EthereumJSONRPCResponse[EthereumBlock]
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return EthereumBlock{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return rpcResponse.Result, nil
}
