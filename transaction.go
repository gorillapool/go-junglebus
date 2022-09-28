package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

func (jb *JungleBusClient) GetTransaction(ctx context.Context, txID string) (*models.Transaction, error) {
	return jb.transport.GetTransaction(ctx, txID)
}
