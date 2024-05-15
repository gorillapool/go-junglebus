package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

// GetAddressTransactions get transaction meta data for the given address
func (jb *Client) GetAddressTransactions(ctx context.Context, address string, fromHeight uint32) ([]*models.AddressTx, error) {
	return jb.transport.GetAddressTransactions(ctx, address, fromHeight)
}

// GetAddressTransactionDetails get full transaction data for the given address
func (jb *Client) GetAddressTransactionDetails(ctx context.Context, address string, fromHeight uint32) ([]*models.Transaction, error) {
	return jb.transport.GetAddressTransactionDetails(ctx, address, fromHeight)
}
