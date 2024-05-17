package models

// AddressTx struct
type AddressTx struct {
	ID            string `json:"id,omitempty" bson:"_id"`
	PKHash        []byte `json:"-" bson:"-"`
	TransactionID string `json:"transaction_id" bson:"transaction_id"`
	BlockHeight   uint32 `json:"block_height" bson:"block_height"`
	BlockHash     string `json:"block_hash" bson:"block_hash"`
	BlockIndex    uint64 `json:"block_index" db:"block_index" bson:"block_index"`
	BlockPage     uint32 `json:"block_page" db:"block_page" bson:"block_page"`
}
