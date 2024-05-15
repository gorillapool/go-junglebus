package junglebus

import (
	"context"
)

func (jb *Client) GetTxo(ctx context.Context, txID string, vout uint32) ([]byte, error) {
	return jb.transport.GetTxo(ctx, txID, vout)
}

func (jb *Client) GetSpend(ctx context.Context, txID string, vout uint32) ([]byte, error) {
	return jb.transport.GetSpend(ctx, txID, vout)
}
