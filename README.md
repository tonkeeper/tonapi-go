# TonAPI SDK

## Description

This repository contains [TonAPI](https://tonapi.io) SDK and examples.

The native TON's RPC is very low-level and is not suitable for building applications on top of it. [TonAPI](https://tonapi.io) aims at speeding up development of TON-based applications and provides an API centered around high-level concepts like Jettons, NFTs and so on, while keeping a way to access low-level details.

Check out more details at [TonAPI Documentation](https://docs.tonconsole.com/tonapi/rest-api).

## Installation

To install the TonAPI SDK, run:

```bash
go get github.com/tonkeeper/tonapi-go
```

## Quick Start

Here's a simple example to get you started:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tonkeeper/tonapi-go"
)

func main() {
	// Create a new client with default settings
	// If you want to use testnet, use tonapi.TestnetTonApiURL
	// You can use TonAPI.io without a token by passing &tonapi.Security{} as the second parameter,
	// but note that TonAPI.io has strict rate limits, so it's better to get a Token from tonconsole.com 
	// in the TonAPI section - it's completely free
	client, err := tonapi.NewClient(tonapi.TonApiURL, &tonapi.Security{})
	if err != nil {
		log.Fatal(err)
	}
	
	// Get account information
	account, err := client.GetAccount(context.Background(), tonapi.GetAccountParams{
		AccountID: "EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0",
	})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Account Balance: %v\n", account.Balance)
}
```

## Configuration Options

### Network Selection

You can specify which network to use:

```go
// For mainnet
client, err := tonapi.NewClient(tonapi.TonApiURL, &tonapi.Security{})

// For testnet
client, err := tonapi.NewClient(tonapi.TestnetTonApiURL, &tonapi.Security{})
```

### Authentication

While TonAPI can be used without authentication, it's recommended to obtain an API token to avoid rate limits:

```go
// Using API token (recommended)
token := "your-api-token" // Get your free token from tonconsole.com
client, err := tonapi.NewClient(tonapi.TonApiURL, tonapi.WithToken(token))
```

You can get your free API token from [TON Console](https://tonconsole.com/tonapi/api-keys) in the TonAPI section.

### Custom HTTP Client

You can also use a custom client by using WithClient, where you can, for example, pass a throttled client:

```go
import (
	"net/http"
	"time"
	
	"github.com/tonkeeper/tonapi-go"
	"golang.org/x/time/rate"
)

func main() {
	// Create a throttled client
	throttledClient := &http.Client{
		Transport: throttled.NewTransport(
			http.DefaultTransport,
			rate.NewLimiter(1, 1)), // Set values according to the rate plan of the token
	}
	
	// Use custom client with token
	token := "your-api-token"
	client, err := tonapi.NewClient(
		tonapi.TonApiURL, 
		tonapi.WithToken(token), 
		tonapi.WithClient(throttledClient),
	)
	if err != nil {
		// handle error
	}
	
	// Use client...
}
```

## Common Operations

### Get Account Information

```go
account, err := client.GetAccount(context.Background(), tonapi.GetAccountParams{
	AccountID: "EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0",
})
```

### Get Transactions

```go
transactions, err := client.GetTransactions(context.Background(), tonapi.GetTransactionsParams{
	AccountID: "EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0",
	Limit:     10,
})
```

### Get Jettons

```go
jettons, err := client.GetAccountJettons(context.Background(), tonapi.GetAccountJettonsParams{
	AccountID: "EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0",
})
```

## Error Handling

Always check for errors when making API calls:

```go
result, err := client.GetAccount(context.Background(), params)
if err != nil {
    // Check if it's a TonAPI error
    if apiErr, ok := err.(*tonapi.Error); ok {
        fmt.Printf("API Error: %v\n", apiErr.Error)
    } else {
        fmt.Printf("Error: %s\n", err.Error())
    }
    return
}
```

## Rate Limiting

TonAPI has rate limits based on your authentication:
- Anonymous users: Strict rate limits
- API token users: Higher limits based on your plan

For high-volume applications, consider implementing a throttled client as shown in the examples.

## Best Practices

1. Always use an API token for production applications
2. Implement proper error handling for all API calls
3. Use a custom HTTP client with rate limiting for high-volume applications
4. Consider caching frequently accessed data

## Documentation

For complete API documentation, visit the [Documentation](https://docs.tonconsole.com) website.

## Support

For support and questions, join the [TonApi Tech](https://t.me/tonapitech) on Telegram.

## REST API

You can always find the latest version of TonAPI REST API documentation at [TonAPI Documentation](https://docs.tonconsole.com/tonapi/rest-api).

[TonAPI SDK example](examples/tonapi-sdk/main.go) shows how to work with REST API in golang.

## Streaming API

Usually, an application needs to monitor the blockchain for specific events and act accordingly.    
TonAPI offers two ways to do it: SSE and Websocket.

The advantage of Websocket is that it can be reconfigured dynamically to subscribe/unsubscribe to/from specific events,
whereas SSE has to reconnect to TonAPI to change the list of events it is subscribed to.

Take a look at [SSE example](examples/sse/main.go) and [Websocket example](examples/websocket/main.go) to see how to work with TonAPI Streaming API in golang.

More details can be found at [TonAPI Streaming API Documentation](https://docs.tonconsole.com/tonapi/streaming-api).