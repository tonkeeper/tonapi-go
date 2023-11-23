package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tonkeeper/tonapi-go"
)

func subscribeToMempool(token string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
	for {
		err := streamingAPI.SubscribeToMempool(context.Background(),
			// this "accounts" parameter is optional,
			// if not set, you will receive all mempool events.
			// if defined, you will receive only mempool events that involve these accounts.
			[]string{"-1:5555555555555555555555555555555555555555555555555555555555555555"},
			func(data tonapi.MempoolEventData) {
				value, _ := json.Marshal(data)
				fmt.Printf("mempool event: %#v\n", value)
			})
		if err != nil {
			fmt.Printf("mempool error: %v, reconnecting...\n", err)
		}
	}
}

func subscribeToTransactions(token string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
	for {
		err := streamingAPI.SubscribeToTransactions(context.Background(),
			[]string{"-1:5555555555555555555555555555555555555555555555555555555555555555"},
			// this "operations" is optional,
			// if not set, you will receive all transactions.
			// if defined, you will receive only transactions with these operations.
			nil,
			func(data tonapi.TransactionEventData) {
				fmt.Printf("New tx with hash: %v\n", data.TxHash)
			})
		if err != nil {
			fmt.Printf("tx error: %v, reconnecting...\n", err)
		}
	}
}

func subscribeToTraces(token string) {
	streamingAPI := tonapi.NewStreamingAPI(tonapi.WithStreamingToken(token))
	for {
		err := streamingAPI.SubscribeToTraces(context.Background(), []string{"-1:5555555555555555555555555555555555555555555555555555555555555555"},
			func(data tonapi.TraceEventData) {
				fmt.Printf("New trace with hash: %v\n", data.Hash)
			})
		if err != nil {
			fmt.Printf("trace error: %v, reconnecting...\n", err)
		}
	}
}

func main() {
	// When working with tonapi.io, you should consider getting an API key at https://tonconsole.com/
	// because tonapi.io has per-ip limits for sse and websocket connections.
	//
	// You can configure it with:
	//         streamingAPI = tonapi.NewStreamingAPI(tonapi.WithStreamingToken("<private-key>"))
	//
	// To work with a local version of tonapi.io (opentonapi) use:
	//         streamingAPI = tonapi.NewStreamingAPI(tonapi.WithStreamingEndpoint("http://127.0.0.1:8081"))
	//
	token := ""

	go subscribeToTraces(token)
	go subscribeToMempool(token)
	go subscribeToTransactions(token)
	select {}
}
