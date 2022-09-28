package models

// BlockHeader is an object representing the BitCoin block header
// This comes from Bux models_block_headers.go
type BlockHeader struct {
	Hash       string `json:"hash"`
	Coin       uint32 `json:"coin"`
	Height     uint32 `json:"height"`
	Time       uint32 `json:"time"`
	Nonce      uint32 `json:"nonce"`
	Version    uint32 `json:"version"`
	MerkleRoot string `json:"merkleroot"`
	Bits       string `json:"bits"`
	Synced     uint64 `json:"synced"`
}
