package pubsub

import (
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

type ProducerConfig struct {
	Host string
	Port string
}

func NewProducerConfig(host string, port string) *ProducerConfig {
	return &ProducerConfig{host, port}
}

func NewProducer(pc *ProducerConfig) (*nsq.Producer, error) {
	addr := fmt.Sprintf("%s:%s", pc.Host, pc.Port)
	cfg := nsq.NewConfig()
	return nsq.NewProducer(addr, cfg)
}
