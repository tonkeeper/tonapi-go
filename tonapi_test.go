package tonapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/graze/go-throttled"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"testing"
)

var systemAccountID = ton.MustParseAccountID("Ef8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAU")

// TestThrottling tests the client with throttling.
func TestThrottling(t *testing.T) {
	const (
		withTokenName    = "WithToken"
		withoutTokenName = "WithoutToken"
	)
	tests := []struct {
		name      string
		token     string
		rateLimit rate.Limit
		burst     int
	}{
		{
			name:      withTokenName,
			token:     os.Getenv("TONAPI_TOKEN"), // use TonApi token with Lite tier
			rateLimit: 10,
			burst:     10,
		},
		{
			name:      withoutTokenName,
			token:     "",
			rateLimit: 1,
			burst:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			throttledClient := &http.Client{
				Transport: throttled.NewTransport(
					http.DefaultTransport,
					rate.NewLimiter(tt.rateLimit, tt.burst)),
			}
			if tt.token == "" && tt.name == withTokenName {
				t.Skip("skipping test with token")
			}
			var client *Client
			var err error
			if tt.token == "" {
				client, err = NewClient(TonApiURL, &Security{}, WithClient(throttledClient))
			} else {
				client, err = NewClient(TonApiURL, WithToken(tt.token), WithClient(throttledClient))
			}
			if err != nil {
				t.Fatalf("failed to init tonapi client: %v", err)
			}
			for i := 0; i < 30; i++ {
				_, err = client.Status(context.Background())
				require.NoError(t, err)
			}
		})
	}
}

// TestCustomRequest tests the client with custom requests.
func TestCustomRequest(t *testing.T) {
	throttledClient := &http.Client{
		Transport: throttled.NewTransport(
			http.DefaultTransport,
			rate.NewLimiter(1, 1)),
	}
	client, err := NewClient(TonApiURL, &Security{}, WithClient(throttledClient))
	if err != nil {
		t.Fatalf("failed to init tonapi client: %v", err)
	}

	tests := []struct {
		name   string
		method string
		path   string
		query  map[string][]string
		err    error
	}{
		{
			name:   "fail to get account info - method not allowed",
			method: http.MethodPost,
			path:   fmt.Sprintf("v2/accounts/%v", systemAccountID),
			err:    errors.New("405 Method Not Allowed"),
		},
		{
			name:   "ok to get account info",
			method: http.MethodGet,
			path:   fmt.Sprintf("v2/accounts/%v", systemAccountID),
			err:    nil,
		},
		{
			name:   "fail with invalid account ID",
			method: http.MethodGet,
			path:   "v2/accounts/invalidAccountID",
			err:    errors.New("400 Bad Request"),
		},
		{
			name:   "fail with non-existent path",
			method: http.MethodGet,
			path:   "v2/nonexistentpath",
			err:    errors.New("404 Not Found"),
		},
		{
			name:   "ok to get collections",
			method: http.MethodGet,
			path:   "v2/nfts/collections",
			query:  map[string][]string{"limit": {"10"}},
			err:    nil,
		},
		{
			name:   "ok to exec get method",
			method: http.MethodGet,
			path:   "v2/blockchain/accounts/EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs/methods/get_wallet_address",
			query:  map[string][]string{"args": {"UQDNzlh0XSZdb5_Qrlx5QjyZHVAO74v5oMeVVrtF_5Vt1rIt", "UQBVXzBT4lcTA3S7gxrg4hnl5fnsDKj4oNEzNp09aQxkwmCa"}},
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Request(context.Background(), tt.method, tt.path, tt.query, nil)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
