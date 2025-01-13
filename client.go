package tonapi

import (
	"context"
	"github.com/ogen-go/ogen/ogenerrors"
)

type Security struct {
	Token string
}

// BearerAuth returns the Bearer authentication credentials for the client.
// If the token is not provided, it skips the security check by returning `ErrSkipClientSecurity`.
// This method is used to authenticate API requests using the Bearer token.
func (s *Security) BearerAuth(ctx context.Context, operationName OperationName, client *Client) (BearerAuth, error) {
	if s.Token == "" {
		return BearerAuth{}, ogenerrors.ErrSkipClientSecurity
	}
	return BearerAuth{Token: s.Token}, nil
}

// WithToken configures the Security object with the provided token for Bearer authentication.
// The token will be used to authorize API requests to tonapi.io.
// To obtain a token, register and generate an API key at https://tonconsole.com.
//
// Example:
//
// import (
//
//	"github.com/tonkeeper/tonapi-go"
//
// )
//
//	func main() {
//	    token := "your-api-token"
//	    client, err := tonapi.New(tonapi.WithToken(token))
//	    if err != nil {
//	        // handle error
//	    }
//	    // use client
//	}
func WithToken(token string) *Security {
	return &Security{Token: token}
}

// TonApiURL is the endpoint for working with the mainnet.
const TonApiURL = "https://tonapi.io"

// TestnetTonApiURL is the endpoint for working with the testnet.
// Example:
// client, err := NewClient(tonapi.TestnetTonApiURL)
const TestnetTonApiURL = "https://testnet.tonapi.io"

// You can also use a custom client by using WithClient, where you can, for example, pass a throttled client.
//
// Example:
//
// import (
//
//	"github.com/tonkeeper/tonapi-go"
//
// )
//
//	func main() {
//	    throttledClient := &http.Client{
//			Transport: throttled.NewTransport(
//				http.DefaultTransport,
//				rate.NewLimiter(1, 1)), // Set values according to the rate plan of the token
//			}
//	    token := "your-api-token"
//	    client, err := tonapi.New(tonapi.WithToken(token), tonapi.WithClient(throttledClient))
//	    if err != nil {
//	        // handle error
//	    }
//	    // use client
//	}

// New returns a new Client instance.
func New(sec SecuritySource, opts ...ClientOption) (*Client, error) {
	if sec == nil {
		sec = &Security{}
	}
	return NewClient(TonApiURL, sec, opts...)
}
