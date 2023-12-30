package natspub

import (
	"encoding/json"
	"flag"
	"github.com/nats-io/nats.go"
	sctx "github.com/phathdt/service-context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-jetstream/pkg/jetstream"
	"github.com/ThreeDotsLabs/watermill/message"
)

type natsPub struct {
	id        string
	natsURL   string
	publisher *jetstream.Publisher
	logger    sctx.Logger
}

func New(id string) *natsPub {
	return &natsPub{id: id}
}

func (n *natsPub) ID() string {
	return n.id
}

func (n *natsPub) InitFlags() {
	flag.StringVar(&n.natsURL, "nats-sub-uri", "nats://localhost:4222", "nats uri")
}

func (n *natsPub) Activate(sc sctx.ServiceContext) error {
	n.logger = sctx.GlobalLogger().GetLogger(n.id)

	n.logger.Info("Connect to nats at ", n.natsURL, " ...")

	options := []nats.Option{
		nats.RetryOnFailedConnect(true),
		nats.Timeout(30 * time.Second),
		nats.ReconnectWait(1 * time.Second),
	}
	marshaler := &jetstream.GobMarshaler{}
	logger := watermill.NewStdLogger(false, false)

	publisher, err := jetstream.NewPublisher(
		jetstream.PublisherConfig{
			URL:         n.natsURL,
			NatsOptions: options,
			Marshaler:   marshaler,
		},
		logger,
	)
	if err != nil {
		n.logger.Error("Error connect to nats at ", n.natsURL, ". ", err.Error())
		return err
	}

	n.publisher = publisher

	return nil
}

func (n *natsPub) Stop() error {
	if n.publisher != nil {
		if err := n.publisher.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (n *natsPub) Publish(topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	if err = n.publisher.Publish(topic, msg); err != nil {
		return err
	}

	n.logger.Infof("natsPub %s message = %+v\n", topic, string(msg.Payload))

	return nil
}

func (n *natsPub) PublishRaw(topic string, payload []byte) error {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	if err := n.publisher.Publish(topic, msg); err != nil {
		return err
	}

	n.logger.Infof("natsPub %s message = %+v\n", topic, string(msg.Payload))

	return nil
}
