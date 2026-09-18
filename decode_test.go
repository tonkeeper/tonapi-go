package tonapi

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/stretchr/testify/require"
)

const (
	jettonPreviewJSON = `{
		"address": "0:97264395BD65A255A429B11326C84128B7D70FFED7949ABAE3036D506BA38621",
		"name": "Some Jetton",
		"symbol": "SJ",
		"decimals": 9,
		"image": "https://example.com/jetton.png",
		"verification": "whitelist",
		"score": 100
	}`
	accountAddressJSON = `{
		"address": "0:410A6B1A9C2A5B3C1C4D8E7F0A1B2C3D4E5F60718293A4B5C6D7E8F901A2B3C4",
		"is_scam": false,
		"is_wallet": true
	}`
)

// TestDecodeIgnoresUnknownFields tests that fields missing from the spec are skipped, not rejected.
func TestDecodeIgnoresUnknownFields(t *testing.T) {
	const payload = `{
		"address": "0:97264395BD65A255A429B11326C84128B7D70FFED7949ABAE3036D506BA38621",
		"name": "Some Jetton",
		"symbol": "SJ",
		"decimals": 9,
		"image": "https://example.com/jetton.png",
		"verification": "whitelist",
		"score": 100,
		"field_added_after_this_sdk_was_released": "whatever",
		"nested_addition": {"deeply": [{"nested": true}, null]}
	}`

	var jetton JettonPreview
	require.NoError(t, jetton.Decode(jx.DecodeStr(payload)))
	require.NoError(t, jetton.Validate())
	require.Equal(t, "SJ", jetton.Symbol)
	require.Equal(t, JettonVerificationTypeWhitelist, jetton.Verification)
	require.Equal(t, 9, jetton.Decimals)
}

// TestDecodeUnknownEnumValue tests that an unlisted enum value decodes as a raw string but fails
// validation, which .ogen.yml runs on every response. A new enum value needs a spec refresh;
// TestSpecEnumsMatchLiveAPI is what catches that in time.
func TestDecodeUnknownEnumValue(t *testing.T) {
	var implementation PoolImplementationType
	require.NoError(t, implementation.Decode(jx.DecodeStr(`"some_new_staking_pool"`)))
	require.Equal(t, PoolImplementationType("some_new_staking_pool"), implementation)
	require.ErrorContains(t, implementation.Validate(), "invalid value")

	require.NoError(t, implementation.Decode(jx.DecodeStr(`"hipo"`)))
	require.Equal(t, PoolImplementationTypeHipo, implementation)
	require.NoError(t, implementation.Validate())
}

// TestDecodeRejectsMissingRequiredField tests that a missing required field is an error naming the field.
func TestDecodeRejectsMissingRequiredField(t *testing.T) {
	tests := []struct {
		name    string
		decode  func(*jx.Decoder) error
		payload string
		missing string
	}{
		{
			name:    "GasRelayFee without amount",
			decode:  func(d *jx.Decoder) error { var v GasRelayFee; return v.Decode(d) },
			payload: `{"jetton": ` + jettonPreviewJSON + `}`,
			missing: "amount",
		},
		{
			name:    "MigrationTransaction without gas_spent",
			decode:  func(d *jx.Decoder) error { var v MigrationTransaction; return v.Decode(d) },
			payload: `{"seqno": 1, "boc": "te6cc", "messages": []}`,
			missing: "gas_spent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.decode(jx.DecodeStr(tt.payload))
			require.Error(t, err)
			require.ErrorContains(t, err, tt.missing)
		})
	}
}

// TestDecodeGasRelayAction tests decoding of is_battery and relayer_fee, with and without them.
func TestDecodeGasRelayAction(t *testing.T) {
	payload := `{
		"amount": 5000000,
		"relayer": ` + accountAddressJSON + `,
		"target": ` + accountAddressJSON + `,
		"is_battery": true,
		"relayer_fee": {"jetton": ` + jettonPreviewJSON + `, "amount": "1000000000"}
	}`

	var action GasRelayAction
	require.NoError(t, action.Decode(jx.DecodeStr(payload)))
	require.NoError(t, action.Validate())
	require.Equal(t, int64(5000000), action.Amount)
	require.True(t, action.IsBattery.Value)

	fee, ok := action.RelayerFee.Get()
	require.True(t, ok)
	require.Equal(t, "1000000000", fee.Amount)
	require.Equal(t, "SJ", fee.Jetton.Symbol)

	const withoutNewFields = `{"amount": 1, "relayer": ` + accountAddressJSON + `, "target": ` + accountAddressJSON + `}`
	var bare GasRelayAction
	require.NoError(t, bare.Decode(jx.DecodeStr(withoutNewFields)))
	require.NoError(t, bare.Validate())
	require.False(t, bare.IsBattery.Set)
	require.False(t, bare.RelayerFee.Set)
}

// TestDecodePrepareMigrationResponse tests both arms of the union PrepareMigration returns.
func TestDecodePrepareMigrationResponse(t *testing.T) {
	const prepared = `{
		"from": "0:aaa",
		"to": "0:bbb",
		"wallet_version": "v4R2",
		"transactions": []
	}`
	const conflict = `{
		"from": "0:aaa",
		"to": "0:bbb",
		"wallet_version": "v4R2",
		"transactions": [],
		"error": "insufficient GRAM for gas",
		"error_code": 50000,
		"details": {"required": 120000000, "available": 3000000}
	}`

	t.Run("200 is a prepared migration", func(t *testing.T) {
		res, err := decodePrepareMigrationResponse(jsonResponse(http.StatusOK, prepared))
		require.NoError(t, err)

		response, ok := res.(*MigrationPrepareResponse)
		require.Truef(t, ok, "expected *MigrationPrepareResponse, got %T", res)
		require.Equal(t, "0:aaa", response.From)
		require.Equal(t, "v4R2", response.WalletVersion)
		require.NotNil(t, response.Transactions)
	})

	t.Run("409 is a conflict carrying the shortfall", func(t *testing.T) {
		res, err := decodePrepareMigrationResponse(jsonResponse(http.StatusConflict, conflict))
		require.NoError(t, err)

		response, ok := res.(*MigrationPrepareConflict)
		require.Truef(t, ok, "expected *MigrationPrepareConflict, got %T", res)
		require.Equal(t, "insufficient GRAM for gas", response.Error)
		require.Equal(t, int64(50000), response.ErrorCode.Value)

		details, ok := response.Details.Get()
		require.True(t, ok)
		require.Equal(t, int64(120000000), details.Required)
		require.Equal(t, int64(3000000), details.Available)
	})
}

// TestGetJettonHoldersSortBy tests validation of the sort_by parameter.
func TestGetJettonHoldersSortBy(t *testing.T) {
	require.NoError(t, GetJettonHoldersSortByBalance.Validate())
	require.NoError(t, GetJettonHoldersSortByAddress.Validate())
	require.Error(t, GetJettonHoldersSortBy("by_vibes").Validate())
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
