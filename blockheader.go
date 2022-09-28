package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

// GetBlockHeader get a block header from JungleBus
func (jb *JungleBusClient) GetBlockHeader(ctx context.Context, block string) (*models.BlockHeader, error) {
	return jb.transport.GetBlockHeader(ctx, block)
}

// GetBlockHeaders get a list of block headers from JungleBus
func (jb *JungleBusClient) GetBlockHeaders(ctx context.Context, block string, limit uint) ([]*models.BlockHeader, error) {
	return jb.transport.GetBlockHeaders(ctx, block, limit)
}
