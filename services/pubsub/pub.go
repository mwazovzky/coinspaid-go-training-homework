package pubsub

import (
	"fmt"
	"log"

	nsq "github.com/nsqio/go-nsq"
)

type PubConfig struct {
	Host  string
	Port  string
	Topic string
}

func NewPubConfig() *PubConfig {
	return &PubConfig{}
}

func Pub(pc *PubConfig, payload []byte) {
	cfg := nsq.NewConfig()
	addr := fmt.Sprintf("%s:%s", pc.Host, pc.Port)
	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		log.Fatal("Could create nsq producer", err)
	}

	err = producer.Publish(pc.Topic, payload)
	if err != nil {
		log.Fatal("Could not connect", err)
	}

	producer.Stop()
}
