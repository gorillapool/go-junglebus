package junglebus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
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
	pubChan          chan *pubEvent
	wg               sync.WaitGroup
	closed           bool
}

type pubEvent struct {
	Channel string
	Data    []byte
}

func (s *Subscription) Unsubscribe() (err error) {

	for _, sub := range s.subscriptions {
		err = sub.Unsubscribe()
	}
	if !s.closed {
		close(s.pubChan)
	}
	s.centrifugeClient.Close()
	s.closed = true
	return err
}

func (s *Subscription) addToQueue(event *pubEvent) {
	s.wg.Add(1)
	s.pubChan <- event
}

func (s *Subscription) handlePubChan() {
	for e := range s.pubChan {
		switch e.Channel {
		case "control":
			status := &models.ControlResponse{}
			if err := proto.Unmarshal(e.Data, status); err != nil {
				s.EventHandler.OnError(err)
			} else {
				if status.StatusCode == uint32(SubscriptionBlockDone) {
					currentBlock = status.Block + 1
					currentPage = 0
				}
				if status.StatusCode == uint32(SubscriptionPageDone) {
					currentBlock = status.Block
					currentPage = status.Transactions + 1
				}
				s.EventHandler.OnStatus(status)
			}
		case "main":
			tx := &models.TransactionResponse{}
			if err := proto.Unmarshal(e.Data, tx); err != nil {
				s.EventHandler.OnError(err)
			} else {
				if len(tx.Transaction) == 0 {
					txData, err := s.client.GetTransaction(context.Background(), tx.Id)
					if err != nil {
						s.EventHandler.OnError(err)
						break
					}
					tx.Transaction = txData.Transaction
				}
				currentBlock = tx.BlockHeight
				// TODO: update protobuf to include block page
				// currentPage = tx.BlockPage

				// log.Printf("[TX]: %d %d - %d: %v", tx.BlockHeight, tx.BlockIndex, len(tx.Transaction), tx.Id)
				s.EventHandler.OnTransaction(tx)
			}
		case "mempool":
			tx := &models.TransactionResponse{}
			if err := proto.Unmarshal(e.Data, tx); err != nil {
				s.EventHandler.OnError(err)
			} else {
				if len(tx.Transaction) == 0 {
					txData, err := s.client.GetTransaction(context.Background(), tx.Id)
					if err != nil {
						s.EventHandler.OnError(err)
						break
					}
					tx.Transaction = txData.Transaction
				}
				s.EventHandler.OnMempool(tx)
			}
		}
		s.wg.Done()
	}
}

func (jb *Client) Unsubscribe() (err error) {
	jb.subscription.Unsubscribe()
	jb.subscription = nil
	return nil
}

// var currentPage uint64
var currentBlock uint32
var currentPage uint64

func (jb *Client) Subscribe(ctx context.Context, subscriptionID string, fromBlock uint64, eventHandler EventHandler) (*Subscription, error) {
	return jb.SubscribeWithQueue(ctx, subscriptionID, fromBlock, 0, eventHandler, 100000)
}

func (jb *Client) SubscribeWithQueue(ctx context.Context, subscriptionID string, fromBlock uint64, fromPage uint64, eventHandler EventHandler, queueSize uint32) (*Subscription, error) {
	var subs *Subscription
	currentBlock = uint32(fromBlock)
	currentPage = fromPage

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
			for _, sub := range subs.subscriptions {
				sub.Unsubscribe()
			}
			subs.wg.Wait()
			eventHandler.OnStatus(&models.ControlResponse{
				StatusCode: uint32(StatusConnecting),
				Status:     "reconnecting",
				Message:    fmt.Sprintf("Reconnecting to server at block %d, page %d", currentBlock, currentPage),
			})
			_ = jb.Unsubscribe()
			time.Sleep(1 * time.Second)
			_, err = jb.SubscribeWithQueue(ctx, subscriptionID, uint64(currentBlock), currentPage, eventHandler, queueSize)
			if err != nil {
				eventHandler.OnError(err)
			}
			return
		}

		jb.subscription = subs

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
	subs = &Subscription{
		SubscriptionID:   subscriptionID,
		FromBlock:        fromBlock,
		EventHandler:     eventHandler,
		client:           jb,
		centrifugeClient: centrifugeClient,
		subscriptions:    map[string]*centrifuge.Subscription{},
		pubChan:          make(chan *pubEvent, queueSize),
	}
	go subs.handlePubChan()

	centrifugeClient.OnPublication(func(e centrifuge.ServerPublicationEvent) {
		log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
		if strings.Contains(e.Channel, ":control") {
			subs.addToQueue(&pubEvent{
				Channel: "control",
				Data:    e.Data,
			})
		} else if strings.Contains(e.Channel, ":mempool") {
			subs.addToQueue(&pubEvent{
				Channel: "mempool",
				Data:    e.Data,
			})
		} else {
			subs.addToQueue(&pubEvent{
				Channel: "main",
				Data:    e.Data,
			})
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

	if subs.subscriptions["control"], err = subs.startSubscription(`query:` + subscriptionID + `:control`); err != nil {
		return nil, err
	}
	subs.subscriptions["control"].OnPublication(func(e centrifuge.PublicationEvent) {
		subs.addToQueue(&pubEvent{
			Channel: "control",
			Data:    e.Data,
		})
	})

	if eventHandler.OnTransaction != nil {
		if subs.subscriptions["main"], err = subs.startSubscription(`query:` + subscriptionID + `:` + strconv.FormatUint(fromBlock, 10) + `:` + strconv.FormatUint(currentPage, 10)); err != nil {
			// if subs.subscriptions["main"], err = subs.startSubscription(`query:` + subscriptionID + `:` + strconv.FormatUint(fromBlock, 10)); err != nil {
			return nil, err
		}
		subs.subscriptions["main"].OnPublication(func(e centrifuge.PublicationEvent) {
			subs.addToQueue(&pubEvent{
				Channel: "main",
				Data:    e.Data,
			})
		})
	}

	if eventHandler.OnMempool != nil {
		if subs.subscriptions["mempool"], err = subs.startSubscription(`query:` + subscriptionID + `:mempool`); err != nil {
			return nil, err
		}
		subs.subscriptions["mempool"].OnPublication(func(e centrifuge.PublicationEvent) {
			subs.addToQueue(&pubEvent{
				Channel: "mempool",
				Data:    e.Data,
			})
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
