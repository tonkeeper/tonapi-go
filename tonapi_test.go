package tonapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
)

var systemAccountID = ton.MustParseAccountID("Ef8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAU")

func TestCustomRequest(t *testing.T) {
	client, err := New()
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
			time.Sleep(time.Millisecond * 100) // rps limit
		})
	}
}
