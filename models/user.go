package models

type User struct {
	ID          string  `json:"id" bson:"id"`
	Username    string  `json:"username" bson:"username"`
	FullName    string  `json:"full_name" bson:"full_name"`
	Email       string  `json:"email" bson:"email"`
	Avatar      string  `json:"avatar" bson:"avatar"`
	PrivateKey  string  `json:"private_key" bson:"private_key"`
	PublicKey   string  `json:"public_key" bson:"public_key"`
	FundAddress string  `json:"fund_address" bson:"fund_address"`
	FundHeight  uint32  `json:"fund_height" bson:"fund_height"`
	FundTotal   uint64  `json:"fund_total" bson:"fund_total"`
	FundUsed    float64 `json:"fund_used" bson:"fund_used"`
	FundBalance float64 `json:"fund_balance" bson:"fund_balance"`
}
