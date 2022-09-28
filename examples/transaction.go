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
		panic("no transaction id given")
	}
	txID := argsWithoutProg[0]

	var tx *models.Transaction
	if tx, err = junglebusClient.GetTransaction(context.Background(), txID); err != nil {
		log.Printf("ERROR: failed getting transaction %s", err.Error())
	} else {
		j, _ := json.Marshal(tx)
		log.Printf("Got transaction %s", string(j))
	}
	os.Exit(0)
}
