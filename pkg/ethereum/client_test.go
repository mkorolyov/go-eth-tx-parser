package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
)

type MockHttpTransport struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestEthClient_GetBlockNumberSmoke(t *testing.T) {
	client := NewJsonRPCClient()

	blockNumber, err := client.GetBlockNumber(context.Background())
	if err != nil {
		t.Fatalf("expected no error, %v", err)
	}

	if blockNumber == 0 {
		t.Fatalf("expected block number > 0, got %d", blockNumber)
	}
}

func TestEthClient_GetBlockNumber(t *testing.T) {
	mockHTTPTransport := &MockHttpTransport{}
	client := NewJsonRPCClient(WithHTTPClient(http.Client{Transport: mockHTTPTransport}))

	randId := rand.Int()
	ctx := context.Background()

	t.Run("successful response", func(t *testing.T) {
		response := `{"jsonrpc":"2.0","id":` + fmt.Sprint(randId) + `,"result":"0x10d4f"}`
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(response)),
			}, nil
		}

		blockNumber, err := client.GetBlockNumber(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if blockNumber != 0x10d4f {
			t.Fatalf("expected block number 0x10d4f, got %d", blockNumber)
		}
	})

	t.Run("empty block number", func(t *testing.T) {
		response := `{"jsonrpc":"2.0","id":` + fmt.Sprint(randId) + `,"result":""}`
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(response)),
			}, nil
		}

		_, err := client.GetBlockNumber(ctx)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})

	t.Run("failed to decode response", func(t *testing.T) {
		response := `invalid json`
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(response)),
			}, nil
		}

		_, err := client.GetBlockNumber(ctx)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})

	t.Run("failed to make request", func(t *testing.T) {
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("request failed")
		}

		_, err := client.GetBlockNumber(ctx)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})
}

func TestEthClient_GetBlockByNumberSmoke(t *testing.T) {
	client := NewJsonRPCClient()
	ctx := context.Background()

	blockNumber, err := client.GetBlockNumber(ctx)
	if err != nil {
		t.Fatalf("expected no error, %v", err)
	}

	block, err := client.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		t.Fatalf("expected no error, %v", err)
	}

	if len(block.Transactions) == 0 {
		t.Fatalf("expected transactions, got none")
	}
}

func TestGetBlockByNumber(t *testing.T) {
	mockHTTPTransport := &MockHttpTransport{}
	client := NewJsonRPCClient(WithHTTPClient(http.Client{Transport: mockHTTPTransport}))
	ctx := context.Background()
	blockNumber := 69007

	t.Run("successful response", func(t *testing.T) {
		response := `{"jsonrpc":"2.0","id":1,"result":{"number":"0x10d4f"}}`
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(response)),
			}, nil
		}

		_, err := client.GetBlockByNumber(ctx, blockNumber)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("failed to decode response", func(t *testing.T) {
		response := `invalid json`
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(response)),
			}, nil
		}

		_, err := client.GetBlockByNumber(ctx, blockNumber)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})

	t.Run("failed to make request", func(t *testing.T) {
		mockHTTPTransport.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("request failed")
		}

		_, err := client.GetBlockByNumber(ctx, blockNumber)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})
}
