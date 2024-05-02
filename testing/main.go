package main

import (
	"context"
	"encoding/json"

	"github.com/GorillaPool/go-junglebus"
)

func main() {
	ctx := context.Background()

	jb, err := junglebus.New(junglebus.WithHTTP("http://localhost:3008"))
	if err != nil {
		panic(err)
	}

	transport := jb.GetTransport()
	err = (*transport).Login(ctx, "shruggr-testnet", "#&XuW*8k5ZPZrgqW")
	if err != nil {
		panic(err)
	}
	user, err := (*transport).GetUser(ctx)
	if err != nil {
		panic(err)
	}

	// tx, err := (*transport).GetTransaction(ctx, "ac725672c8237af8470e5976212444b67b0ffefc20daf07236ba81d9a70c72cd")
	// if err != nil {
	// 	panic(err)
	// }
	out, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(out))
}
