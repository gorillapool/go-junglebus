package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

// GetAddressTransactions get transaction meta data for the given address
func (jb *JungleBusClient) GetAddressTransactions(ctx context.Context, address string) ([]*models.Address, error) {
	return jb.transport.GetAddressTransactions(ctx, address)
}

// GetAddressTransactionDetails get full transaction data for the given address
func (jb *JungleBusClient) GetAddressTransactionDetails(ctx context.Context, address string) ([]*models.Transaction, error) {
	return jb.transport.GetAddressTransactionDetails(ctx, address)
}
