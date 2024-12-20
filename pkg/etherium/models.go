package etherium

// EthereumJSONRPCRequest models JSON-RPC requests
type EthereumJSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// EthereumJSONRPCResponse models JSON-RPC responses
type EthereumJSONRPCResponse[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  T      `json:"result"`
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	From string `json:"from"`
	To   string `json:"to"`
	// amount of ETH to transfer from sender to recipient (denominated in WEI, where 1ETH equals 1e+18wei)
	Value string `json:"value"`
	Hash  string `json:"hash"`
}

// EthereumBlock represents an Ethereum block
type EthereumBlock struct {
	Transactions []Transaction `json:"transactions"`
}
