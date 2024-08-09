package tonapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"time"

	"github.com/go-faster/errors"
	ht "github.com/ogen-go/ogen/http"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/tlb"
)

type Custom interface {
	Request(ctx context.Context, method, url string, params map[string]string, data []byte) (json.RawMessage, error)
}

func (c *Client) GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error) {
	res, err := c.GetAccountSeqno(ctx, GetAccountSeqnoParams{AccountID: account.ToRaw()})
	if err != nil {
		return 0, err
	}
	return uint32(res.Seqno), nil
}

func (c *Client) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	var req SendBlockchainMessageReq
	req.Boc.SetTo(base64.StdEncoding.EncodeToString(payload))
	err := c.SendBlockchainMessage(ctx, &req)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tlb.ShardAccount, error) {
	res, err := c.GetBlockchainRawAccount(ctx, GetBlockchainRawAccountParams{AccountID: accountID.ToRaw()})
	if err != nil {
		return tlb.ShardAccount{}, err
	}
	var shardAccount tlb.ShardAccount
	shardAccount.LastTransLt = uint64(res.LastTransactionLt)
	switch res.Status {
	case "nonexist":
		shardAccount.Account.SumType = "AccountNone"
	case "uninit":
		shardAccount.Account.SumType = "Account"
		shardAccount.Account.Account.Addr = accountID.ToMsgAddress()
		shardAccount.Account.Account.Storage.Balance.Grams = tlb.Grams(res.Balance)
		shardAccount.Account.Account.Storage.State.SumType = "AccountUninit"
	case "active":
		shardAccount.Account.SumType = "Account"
		shardAccount.Account.Account.Addr = accountID.ToMsgAddress()
		shardAccount.Account.Account.Storage.Balance.Grams = tlb.Grams(res.Balance)
		shardAccount.Account.Account.Storage.State.SumType = "AccountActive"
		shardAccount.Account.Account.Storage.State.AccountActive.StateInit.Code.Exists = len(res.Code.Value) > 0
		shardAccount.Account.Account.Storage.State.AccountActive.StateInit.Data.Exists = len(res.Data.Value) > 0
	case "frozen":

	}
	return shardAccount, nil
}

// Request sends an HTTP request with the given method, URL, parameters, and data,
// and returns the response as a json.RawMessage.
func (c *Client) Request(ctx context.Context, method, endpoint string, query map[string]string, data []byte) (json.RawMessage, error) {
	const contentType = "application/json"

	// Start measuring the request duration
	startTime := time.Now()
	defer func() {
		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedDuration := time.Since(startTime)
		c.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond))
	}()

	// Parse the full URL by resolving the endpoint relative to the server URL
	u := c.serverURL.ResolveReference(&url.URL{Path: endpoint})

	// Add query parameters to the URL if any
	if query != nil {
		q := u.Query()
		for key, value := range query {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Create the request
	req, err := ht.NewRequest(ctx, method, u)
	if err != nil {
		// Increment the error counter
		c.errors.Add(ctx, 1)
		return nil, err
	}
	if data != nil {
		ht.SetBody(req, bytes.NewReader(data), contentType)
	}

	// Set the content type header
	req.Header.Set("Content-Type", contentType)

	// Send the request using the baseClient's HTTP client
	resp, err := c.cfg.Client.Do(req) // Use the appropriate client or config
	if err != nil {
		// Increment the error counter
		c.errors.Add(ctx, 1)
		return nil, err
	}
	defer resp.Body.Close()

	c.requests.Add(ctx, 1)

	// Check if the response status code indicates an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Increment the error counter
		c.errors.Add(ctx, 1)
		return nil, errors.New(resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Increment the error counter
		c.errors.Add(ctx, 1)
		return nil, err
	}

	// Unmarshal the response body into json.RawMessage
	var jsonResponse json.RawMessage
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		// Increment the error counter
		c.errors.Add(ctx, 1)
		return nil, err
	}

	return jsonResponse, nil
}
