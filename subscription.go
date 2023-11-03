package junglebus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/GorillaPool/go-junglebus/models"
	"github.com/centrifugal/centrifuge-go"
	"google.golang.org/protobuf/proto"
)

type Subscription struct {
	SubscriptionID   string
	FromBlock        uint64
	EventHandler     EventHandler
	client           *Client
	centrifugeClient *centrifuge.Client
	subscriptions    map[string]*centrifuge.Subscription
}

type pubEvent struct {
	Channel string
	Data    []byte
}

func (s *Subscription) Unsubscribe() (err error) {
	for _, sub := range s.subscriptions {
		err = sub.Unsubscribe()
	}
	s.centrifugeClient.Close()

	return err
}

func (jb *Client) Unsubscribe() (err error) {
	for _, sub := range jb.subscription.subscriptions {
		if err = sub.Unsubscribe(); err != nil {
			return err
		}
	}

	jb.subscription.centrifugeClient.Close()
	jb.subscription = nil

	return nil
}

var pubChan chan *pubEvent
var lastBlock uint64

func (jb *Client) Subscribe(ctx context.Context, subscriptionID string, fromBlock uint64, eventHandler EventHandler) (*Subscription, error) {

	var subs *Subscription
	lastBlock = fromBlock

	var err error
	token := jb.transport.GetToken()
	if token == "" {
		// get a new subscription token to use for all requests
		if token, err = jb.transport.GetSubscriptionToken(ctx, subscriptionID); err != nil {
			return nil, err
		}
		if token != "" {
			jb.transport.SetToken(token)
		}
	}

	protocol := "wss"
	if !jb.transport.IsSSL() {
		protocol = "ws"
	}
	url := fmt.Sprintf("%s://%s/connection/websocket?format=protobuf", protocol, jb.transport.GetServerURL())
	centrifugeClient := centrifuge.NewProtobufClient(url, centrifuge.Config{
		Token: token,
		GetToken: func(event centrifuge.ConnectionTokenEvent) (string, error) {
			return jb.transport.RefreshToken(ctx)
		},
		Name:               "go-junglebus",
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       2 * time.Second,
		HandshakeTimeout:   30 * time.Second,
		MaxServerPingDelay: 30 * time.Second,
		EnableCompression:  true,
	})

	centrifugeClient.OnConnecting(func(e centrifuge.ConnectingEvent) {
		if jb.subscription != nil {
			eventHandler.OnStatus(&models.ControlResponse{
				StatusCode: uint32(StatusConnecting),
				Status:     "reconnecting",
				Message:    "Reconnecting to server at block " + strconv.FormatUint(lastBlock, 10),
			})
			_ = jb.Unsubscribe()
			_, _ = jb.Subscribe(ctx, subscriptionID, lastBlock, eventHandler)
			return
		}

		jb.subscription = subs
		close(pubChan)
		pubChan = make(chan *pubEvent, 100000)
		go handlePubChan(eventHandler, jb)

		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusConnecting),
			Status:     "connecting",
			Message:    "Connecting to server",
		})
	})

	centrifugeClient.OnConnected(func(e centrifuge.ConnectedEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusConnected),
			Status:     "connected",
			Message:    "Connected to server",
		})
	})

	centrifugeClient.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusDisconnected),
			Status:     "disconnected",
			Message:    "Disconnected from server",
		})
	})

	centrifugeClient.OnError(func(e centrifuge.ErrorEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusError),
			Status:     "error",
			Message:    e.Error.Error(),
		})
	})

	centrifugeClient.OnMessage(func(e centrifuge.MessageEvent) {
		log.Printf("Message from server: %s", string(e.Data))
	})

	centrifugeClient.OnSubscribed(func(e centrifuge.ServerSubscribedEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusSubscribed),
			Status:     "subscribed",
			Message:    "Subscribed to " + e.Channel,
		})
	})

	centrifugeClient.OnSubscribing(func(e centrifuge.ServerSubscribingEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusSubscribing),
			Status:     "subscribing",
			Message:    "Subscribing to " + e.Channel,
		})
	})

	centrifugeClient.OnUnsubscribed(func(e centrifuge.ServerUnsubscribedEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusUnsubscribed),
			Status:     "unsubscribed",
			Message:    "Unsubscribed from " + e.Channel,
		})
	})

	centrifugeClient.OnPublication(func(e centrifuge.ServerPublicationEvent) {
		log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
		if strings.Contains(e.Channel, ":control") {
			pubChan <- &pubEvent{
				Channel: "control",
				Data:    e.Data,
			}
		} else if strings.Contains(e.Channel, ":mempool") {
			pubChan <- &pubEvent{
				Channel: "mempool",
				Data:    e.Data,
			}
		} else {
			pubChan <- &pubEvent{
				Channel: "main",
				Data:    e.Data,
			}
		}
	})

	centrifugeClient.OnJoin(func(e centrifuge.ServerJoinEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusJoin),
			Status:     "join",
			Message:    "Joined " + e.Channel,
		})
	})

	centrifugeClient.OnLeave(func(e centrifuge.ServerLeaveEvent) {
		eventHandler.OnStatus(&models.ControlResponse{
			StatusCode: uint32(StatusLeave),
			Status:     "leave",
			Message:    "Left " + e.Channel,
		})
	})

	subs = &Subscription{
		SubscriptionID:   subscriptionID,
		FromBlock:        fromBlock,
		EventHandler:     eventHandler,
		client:           jb,
		centrifugeClient: centrifugeClient,
		subscriptions:    map[string]*centrifuge.Subscription{},
	}

	pubChan = make(chan *pubEvent, 100000)
	go handlePubChan(eventHandler, jb)

	if subs.subscriptions["control"], err = subs.startSubscription(`query:` + subscriptionID + `:control`); err != nil {
		return nil, err
	}
	subs.subscriptions["control"].OnPublication(func(e centrifuge.PublicationEvent) {
		pubChan <- &pubEvent{
			Channel: "control",
			Data:    e.Data,
		}
	})

	if eventHandler.OnTransaction != nil {
		if subs.subscriptions["main"], err = subs.startSubscription(`query:` + subscriptionID + `:` + strconv.FormatUint(fromBlock, 10)); err != nil {
			return nil, err
		}
		subs.subscriptions["main"].OnPublication(func(e centrifuge.PublicationEvent) {
			pubChan <- &pubEvent{
				Channel: "main",
				Data:    e.Data,
			}
		})
	}

	if eventHandler.OnMempool != nil {
		if subs.subscriptions["mempool"], err = subs.startSubscription(`query:` + subscriptionID + `:mempool`); err != nil {
			return nil, err
		}
		subs.subscriptions["mempool"].OnPublication(func(e centrifuge.PublicationEvent) {
			pubChan <- &pubEvent{
				Channel: "mempool",
				Data:    e.Data,
			}
		})
	}

	if err = centrifugeClient.Connect(); err != nil {
		return nil, err
	}

	for _, sub := range subs.subscriptions {
		if err = sub.Subscribe(); err != nil {
			return nil, err
		}
	}

	return subs, nil
}

func (s *Subscription) startSubscription(subscription string) (*centrifuge.Subscription, error) {
	sub, err := s.centrifugeClient.NewSubscription(subscription, centrifuge.SubscriptionConfig{
		Recoverable: true,
	})
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func handlePubChan(eventHandler EventHandler, jb *Client) {
	for e := range pubChan {
		switch e.Channel {
		case "control":
			controlResponse := &models.ControlResponse{}
			if err := proto.Unmarshal(e.Data, controlResponse); err != nil {
				eventHandler.OnError(err)
			} else {
				if controlResponse.StatusCode == uint32(SubscriptionBlockDone) {
					lastBlock = uint64(controlResponse.Block)
				}
				eventHandler.OnStatus(controlResponse)
			}
		case "main":
			transaction := &models.TransactionResponse{}
			if err := proto.Unmarshal(e.Data, transaction); err != nil {
				eventHandler.OnError(err)
			} else {
				if len(transaction.Transaction) == 0 {
					txData, err := jb.GetTransaction(context.Background(), transaction.Id)
					if err != nil {
						eventHandler.OnError(err)
						continue
					}
					transaction.Transaction = txData.Transaction
				}
				eventHandler.OnTransaction(transaction)
			}
		case "mempool":
			transaction := &models.TransactionResponse{}
			if err := proto.Unmarshal(e.Data, transaction); err != nil {
				eventHandler.OnError(err)
			} else {
				if len(transaction.Transaction) == 0 {
					txData, err := jb.GetTransaction(context.Background(), transaction.Id)
					if err != nil {
						eventHandler.OnError(err)
						continue
					}
					transaction.Transaction = txData.Transaction
				}
				eventHandler.OnMempool(transaction)
			}
		}
	}
}
