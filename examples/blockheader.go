package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

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
	if len(argsWithoutProg) == 0 {
		panic("no block header given")
	}
	block := argsWithoutProg[0]

	var blockHeader *models.BlockHeader
	if blockHeader, err = junglebusClient.GetBlockHeader(context.Background(), block); err != nil {
		log.Printf("ERROR: failed getting block header %s", err.Error())
	} else {
		j, _ := json.Marshal(blockHeader)
		log.Printf("Got block header %s", string(j))
	}
	os.Exit(0)
}
