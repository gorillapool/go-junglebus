package models

// Address struct
type Address struct {
	ID            string `json:"id"` // unique od this address record = sha256(Address + TransactionID + BlockIndex)
	Address       string `json:"address"`
	TransactionID string `json:"transaction_id"`
	BlockHash     string `json:"block_hash"`
	BlockIndex    uint64 `json:"block_index"`
}
