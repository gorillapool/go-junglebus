package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/GorillaPool/go-junglebus"
	"github.com/GorillaPool/go-junglebus/models"
)

func main() {
	junglebusClient, err := junglebus.New(
		junglebus.WithHTTP("https://junglebus.gorillapool.io"),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		panic("no block header or limit given")
	}
	block := argsWithoutProg[0]
	limit, _ := strconv.ParseUint(argsWithoutProg[1], 10, 64)

	var blockHeaders []*models.BlockHeader
	if blockHeaders, err = junglebusClient.GetBlockHeaders(context.Background(), block, uint(limit)); err != nil {
		log.Printf("ERROR: failed getting block headers %s", err.Error())
	} else {
		j, _ := json.Marshal(blockHeaders)
		log.Printf("Got block headers %s", string(j))
	}
	os.Exit(0)
}
