package tonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

// JsonRPCRequest represents a request in the JSON-RPC protocol supported by "/v2/websocket" endpoint.
type JsonRPCRequest struct {
	ID      uint64   `json:"id,omitempty"`
	JSONRPC string   `json:"jsonrpc,omitempty"`
	Method  string   `json:"method,omitempty"`
	Params  []string `json:"params,omitempty"`
}

// JsonRPCResponse represents a response in the JSON-RPC protocol supported by "/v2/websocket" endpoint.
type JsonRPCResponse struct {
	ID      uint64          `json:"id,omitempty"`
	JSONRPC string          `json:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type websocketConnection struct {
	// mu protects the handler fields below.
	mu                 sync.Mutex
	requestID          uint64
	conn               *websocket.Conn
	mempoolHandler     MempoolHandler
	transactionHandler TransactionHandler
	traceHandler       TraceHandler
	blockHandler       BlockHandler
}

func (w *websocketConnection) SubscribeToTransactions(accounts []string, operations []string) error {
	params := accounts
	if len(operations) > 0 {
		params = make([]string, 0, len(accounts))
		ops := fmt.Sprintf("operations=%s", strings.Join(operations, ","))
		for _, account := range accounts {
			params = append(params, fmt.Sprintf("%s;%s", account, ops))
		}
	}
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "subscribe_account", Params: params}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) UnsubscribeFromTransactions(accounts []string) error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "unsubscribe_account", Params: accounts}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) SubscribeToTraces(accounts []string) error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "subscribe_trace", Params: accounts}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) UnsubscribeFromTraces(accounts []string) error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "unsubscribe_trace", Params: accounts}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) SubscribeToMempool(accounts []string) error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "subscribe_mempool"}
	if len(accounts) > 0 {
		request.Params = []string{
			fmt.Sprintf("accounts=%s", strings.Join(accounts, ",")),
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) UnsubscribeFromMempool() error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "unsubscribe_mempool"}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) SubscribeToBlocks(workchain *int) error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "subscribe_block"}
	if workchain != nil {
		request.Params = []string{
			fmt.Sprintf("workchain=%d", *workchain),
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) UnsubscribeFromBlocks() error {
	request := JsonRPCRequest{ID: w.currentRequestID(), JSONRPC: "2.0", Method: "unsubscribe_block"}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(request)
}

func (w *websocketConnection) SetMempoolHandler(handler MempoolHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mempoolHandler = handler
}

func (w *websocketConnection) SetTransactionHandler(handler TransactionHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.transactionHandler = handler
}

func (w *websocketConnection) SetTraceHandler(handler TraceHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.traceHandler = handler
}

func (w *websocketConnection) SetBlockHandler(handler BlockHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.blockHandler = handler
}

func websocketConnect(ctx context.Context, endpoint string, apiKey string) (*websocketConnection, error) {
	header := http.Header{}
	if len(apiKey) > 0 {
		header.Set("Authorization", "bearer "+apiKey)
	}
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	switch endpointUrl.Scheme {
	case "http":
		endpointUrl.Scheme = "ws"
	case "https":
		endpointUrl.Scheme = "wss"
	}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, fmt.Sprintf("%s/v2/websocket", endpointUrl.String()), header)
	if err != nil {
		return nil, err
	}
	return &websocketConnection{
		conn:               conn,
		mempoolHandler:     func(data MempoolEventData) {},
		transactionHandler: func(data TransactionEventData) {},
		traceHandler:       func(data TraceEventData) {},
		blockHandler:       func(data BlockEventData) {},
	}, nil
}

func (w *websocketConnection) runJsonRPC(ctx context.Context, fn WebsocketConfigurator) error {
	defer w.conn.Close()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return fn(w)
	})
	g.Go(func() error {
		for {
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			var response JsonRPCResponse
			if err := json.Unmarshal(msg, &response); err != nil {
				return err
			}
			switch response.Method {
			case "trace":
				var traceEvent TraceEventData
				if err := json.Unmarshal(response.Params, &traceEvent); err != nil {
					return err
				}
				w.processHandler(func() {
					w.traceHandler(traceEvent)
				})
			case "account_transaction":
				var txEvent TransactionEventData
				if err := json.Unmarshal(response.Params, &txEvent); err != nil {
					return err
				}
				w.processHandler(func() {
					w.transactionHandler(txEvent)
				})
			case "mempool_message":
				var mempoolEvent MempoolEventData
				if err := json.Unmarshal(response.Params, &mempoolEvent); err != nil {
					return err
				}
				w.processHandler(func() {
					w.mempoolHandler(mempoolEvent)
				})
			case "block":
				var block BlockEventData
				if err := json.Unmarshal(response.Params, &block); err != nil {
					return err
				}
				w.processHandler(func() {
					w.blockHandler(block)
				})
			}
		}
	})
	return g.Wait()
}

func (w *websocketConnection) currentRequestID() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.requestID++
	return w.requestID
}

func (w *websocketConnection) processHandler(fn func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	fn()
}
