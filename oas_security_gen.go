// Code generated by ogen, DO NOT EDIT.

package tonapi

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
)

// SecuritySource is provider of security values (tokens, passwords, etc.).
type SecuritySource interface {
	// BearerAuth provides bearerAuth security value.
	BearerAuth(ctx context.Context, operationName OperationName, client *Client) (BearerAuth, error)
}

func (s *Client) securityBearerAuth(ctx context.Context, operationName OperationName, req *http.Request) error {
	t, err := s.sec.BearerAuth(ctx, operationName, s)
	if err != nil {
		return errors.Wrap(err, "security source \"BearerAuth\"")
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	return nil
}