package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

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
		panic("no subscription id or block height given")
	}
	subscriptionID := argsWithoutProg[0]
	var fromBlock uint64
	if fromBlock, err = strconv.ParseUint(argsWithoutProg[1], 10, 64); err != nil {
		panic("invalid block height given")
	}

	eventHandler := junglebus.EventHandler{
		// do not set this function to leave out mined transactions
		OnTransaction: func(tx *models.TransactionResponse) {
			log.Printf("[TX]: %d: %v", tx.BlockHeight, tx.Id)
		},
		// do not set this function to leave out mempool transactions
		OnMempool: func(tx *models.TransactionResponse) {
			log.Printf("[MEMPOOL TX]: %v", tx.Id)
		},
		OnStatus: func(status *models.ControlResponse) {
			log.Printf("[STATUS]: %v", status)
		},
		OnError: func(err error) {
			log.Printf("[ERROR]: %v", err)
		},
	}

	var subscription *junglebus.Subscription
	if subscription, err = junglebusClient.Subscribe(context.Background(), subscriptionID, fromBlock, eventHandler); err != nil {
		log.Printf("ERROR: failed getting subscription %s", err.Error())
	} else {
		time.Sleep(10 * time.Second) // stop after 10 seconds
		if err = subscription.Unsubscribe(); err != nil {
			log.Printf("ERROR: failed unsubscribing %s", err.Error())
		}
	}
	os.Exit(0)
}
