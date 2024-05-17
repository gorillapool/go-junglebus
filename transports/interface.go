package transports

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
)

// AddressService is the address related requests
type AddressService interface {
	GetAddressTransactions(ctx context.Context, address string, fromHeight uint32) ([]*models.AddressTx, error)
	GetAddressTransactionDetails(ctx context.Context, address string, fromHeight uint32) ([]*models.Transaction, error)
}

// BlockHeaderService is the block header related requests
type BlockHeaderService interface {
	GetBlockHeader(ctx context.Context, block string) (*models.BlockHeader, error)
	GetBlockHeaders(ctx context.Context, fromBlock string, limit uint) ([]*models.BlockHeader, error)
	GetChaintip(ctx context.Context) (*models.BlockHeader, error)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	GetTransaction(ctx context.Context, txID string) (*models.Transaction, error)
	GetRawTransaction(ctx context.Context, txID string) ([]byte, error)
	GetFromBlock(ctx context.Context, subscriptionID string, height uint32, lastIdx uint64) ([]*models.Transaction, error)
	GetLiteFromBlock(ctx context.Context, subscriptionID string, height uint32, lastIdx uint64) ([]*models.TransactionResponse, error)
}

type AuthService interface {
	GetUser(ctx context.Context) (*models.User, error)
}

type TxoService interface {
	GetTxo(ctx context.Context, txID string, vout uint32) ([]byte, error)
	GetSpend(ctx context.Context, txID string, vout uint32) ([]byte, error)
}

// TransportService the transport service interface
type TransportService interface {
	AddressService
	BlockHeaderService
	TransactionService
	TxoService
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
	GetUser(ctx context.Context) (*models.User, error)
}

// LoginResponse response from server on login or token refresh
type LoginResponse struct {
	Token string `json:"token"`
}
