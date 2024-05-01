package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

func (jb *Client) GetTransaction(ctx context.Context, txID string) (*models.Transaction, error) {
	return jb.transport.GetTransaction(ctx, txID)
}

func (jb *Client) GetRawTransaction(ctx context.Context, txID string) ([]byte, error) {
	return jb.transport.GetRawTransaction(ctx, txID)
}
