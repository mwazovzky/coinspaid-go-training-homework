package pubsub

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type SubConfig struct {
	Host                string
	Port                string
	MaxAttempts         uint16
	MaxInFlight         int
	MaxRequeueDelay     time.Duration
	DefaultRequeueDelay time.Duration
	Topic               string
	Channel             string
}

func NewSubConfig() *SubConfig {
	return &SubConfig{}
}

func Sub(sc *SubConfig, handler nsq.Handler) {
	cfg := nsq.NewConfig()

	host := sc.Host
	port := sc.Port
	addr := fmt.Sprintf("%s:%s", host, port)

	topic := sc.Topic
	channel := sc.Channel
	cfg.MaxAttempts = sc.MaxAttempts
	cfg.MaxInFlight = sc.MaxInFlight
	cfg.MaxRequeueDelay = sc.MaxRequeueDelay
	cfg.DefaultRequeueDelay = sc.DefaultRequeueDelay

	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Register message handler
	consumer.AddHandler(handler)

	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	}

	// Listen channel until interrupted via console command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
