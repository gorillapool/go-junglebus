package transports

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

// AddressService is the address related requests
type AddressService interface {
	GetAddressTransactions(ctx context.Context, address string) ([]*models.Address, error)
	GetAddressTransactionDetails(ctx context.Context, address string) ([]*models.Transaction, error)
}

// BlockHeaderService is the block header related requests
type BlockHeaderService interface {
	GetBlockHeader(ctx context.Context, block string) (*models.BlockHeader, error)
	GetBlockHeaders(ctx context.Context, fromBlock string, limit uint) ([]*models.BlockHeader, error)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	GetTransaction(ctx context.Context, txID string) (*models.Transaction, error)
}

// TransportService the transport service interface
type TransportService interface {
	AddressService
	BlockHeaderService
	TransactionService
	Login(ctx context.Context, username string, password string) error
	IsDebug() bool
	SetDebug(debug bool)
	GetToken() string
	GetSubscriptionToken(ctx context.Context, subscriptionID string) (string, error)
	RefreshToken(ctx context.Context) (string, error)
	SetToken(token string)
	SetVersion(version string)
	UseSSL(useSSL bool)
	IsSSL() bool
	GetServerURL() string
}

// LoginResponse response from server on login or token refresh
type LoginResponse struct {
	Token string `json:"token"`
}
