
# TonAPI SDK and API documentation

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

Take a look at more examples at [TonAPI SDK examples](examples/README.md).


[Openapi.yaml](api/openapi.yml) describes the API of both Opentonapi and TonAPI.
