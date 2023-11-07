
# Description

We believe that the native TON's RPC is very low-level.
And it is not suitable for building applications on top of it.

TonAPI aims at speeding up development of TON-based applications and
provides an API centered around high-level concepts like Jettons, NFTs and so on,
keeping a way to access low-level details.

# TonAPI SDK Example 

Development of TON-based applications is much easier with TonAPI SDK:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tonkeeper/tonapi-go"
)

func main() {
	client, err := tonapi.New()
	if err != nil {
		log.Fatal(err)
	}
	account, err := client.GetAccount(context.Background(), tonapi.GetAccountParams{AccountID: "EQBszTJahYw3lpP64ryqscKQaDGk4QpsO7RO6LYVvKHSINS0"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Account: %v\n", account.Balance)
}
```

Take a look at more examples at [TonAPI SDK examples](examples).

## HTTP API

[TonAPI SDK example](examples/tonapi-sdk/main.go) shows how to work with TonAPI HTTP API.

## Streaming API

Usually, an application needs to monitor the blockchain for specific events and act accordingly.    
TonAPI offers two ways to do it: SSE and Websocket.

The advantage of Websocket is that Websocket can be reconfigured dynamically to subscribe/unsubscribe to/from specific events.
Where SSE has to reconnect to TonAPI to change the list of events it is subscribed to.

Take a look at [SSE example](examples/sse/main.go) and [Websocket example](examples/websocket/main.go) to see how to work with TonAPI/Opentonapi Streaming API.


