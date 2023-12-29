package junglebus

import (
	"context"

	"github.com/GorillaPool/go-junglebus/models"
	"github.com/centrifugal/centrifuge-go"
)

// StatusCode defines the codes that can be returned from the control channel of a subscription
type StatusCode uint

var (
	// StatusConnecting is when connecting to a server
	StatusConnecting StatusCode = 1
	// StatusConnected is when connected to a server
	StatusConnected StatusCode = 2
	// StatusJoin is when joining to a server
	StatusJoin StatusCode = 3
	// StatusLeave is when leaving to a server
	StatusLeave StatusCode = 4
	// StatusDisconnecting is when disconnecting from a server
	StatusDisconnecting StatusCode = 10
	// StatusDisconnected is when disconnected from a server
	StatusDisconnected StatusCode = 11
	// StatusSubscribing is when subscribing to a server
	StatusSubscribing StatusCode = 20
	// StatusSubscribed is when subscribed on a server
	StatusSubscribed StatusCode = 21
	// StatusUnsubscribed is when unsubscribed on a server
	StatusUnsubscribed StatusCode = 29
	// SubscriptionWait is sent when the server is waiting for a new block to be ready to send transactions
	SubscriptionWait StatusCode = 100
	// SubscriptionError is sent when an error was encountered
	SubscriptionError StatusCode = 101
	// SubscriptionPageDone is sent when a block-page is done processing
	SubscriptionPageDone StatusCode = 199
	// SubscriptionBlockDone is sent when a block is done processing
	SubscriptionBlockDone StatusCode = 200
	// SubscriptionReorg is sent when a reorg is initialized
	SubscriptionReorg StatusCode = 300
	// StatusError is sent when an error is found
	StatusError StatusCode = 999
)

type EventHandler struct {
	OnTransaction func(tx *models.TransactionResponse)
	OnMempool     func(tx *models.TransactionResponse)
	OnStatus      func(response *models.ControlResponse)
	OnError       func(err error)
	ctx           context.Context
	debug         bool
}

func (e *EventHandler) OnPublish(event centrifuge.PublicationEvent) {

}
