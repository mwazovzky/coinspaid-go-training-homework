package main

import (
	"homework/services/message"
	"homework/services/pubsub"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
)

func init() {
	godotenv.Load()
}

func main() {
	host := os.Getenv("NSQ_LOOKUP_HOST")
	port := os.Getenv("NSQ_LOOKUP_PORT")
	topic := os.Getenv("NSQ_TOPIC")
	channel := os.Getenv("NSQ_CHANNEL")

	cfg := nsq.NewConfig()
	cfg.MaxAttempts = 10
	cfg.MaxInFlight = 5
	cfg.MaxRequeueDelay = time.Second * 900
	cfg.DefaultRequeueDelay = time.Second * 0

	consumer, err := pubsub.NewConsumer(cfg, topic, channel)
	if err != nil {
		log.Fatal(err)
	}

	// Register message handler
	handler := message.NewMessageHandler()
	consumer.AddHandler(handler)

	err = pubsub.Connect(consumer, host, port)
	if err != nil {
		log.Fatal(err)
	}

	// Listen channel until interrupted via console command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
