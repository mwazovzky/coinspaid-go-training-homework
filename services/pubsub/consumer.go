package pubsub

import (
	"fmt"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type ConsumerConfig struct {
	MaxAttempts         uint16
	MaxInFlight         int
	MaxRequeueDelay     time.Duration
	DefaultRequeueDelay time.Duration
}

func NewConsumer(cfg *nsq.Config, topic string, channel string) (*nsq.Consumer, error) {
	return nsq.NewConsumer(topic, channel, cfg)
}

func Connect(consumer *nsq.Consumer, host string, port string) error {
	addr := fmt.Sprintf("%s:%s", host, port)

	return consumer.ConnectToNSQLookupd(addr)
}
