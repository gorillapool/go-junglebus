package models

// Transaction struct
type Transaction struct {
	ID          ByteString `json:"id"`
	Transaction []byte     `json:"transaction"`
	BlockHash   ByteString `json:"block_hash"`
	BlockHeight uint32     `json:"block_height"`
	BlockTime   uint32     `json:"block_time"`
	BlockIndex  uint64     `json:"block_index"`

	// index data
	// input/output types are
	// p2pkh, p2sh, token-stas, opreturn, tokenized, metanet, bitcom, run, map, bap, non-standard etc.
	Addresses   []string     `json:"addresses"`
	Inputs      []ByteString `json:"inputs"`
	Outputs     []ByteString `json:"outputs"`
	InputTypes  []string     `json:"input_types"`
	OutputTypes []string     `json:"output_types"`
	Contexts    []string     `json:"contexts"`     // optional contexts of output types, only for known protocols
	SubContexts []string     `json:"sub_contexts"` // optional sub-contexts of output types, only for known protocols
	Data        []string     `json:"data"`         // optional data of output types, only for known protocols

	// the merkle proof in binary
	MerkleProof []byte `json:"merkle_proof"`
}
