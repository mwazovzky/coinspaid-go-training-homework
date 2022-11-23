package main

import (
	"homework/services/message"
	"homework/services/pubsub"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	cfg := pubsub.NewSubConfig()
	cfg.Host = os.Getenv("NSQ_LOOKUP_HOST")
	cfg.Port = os.Getenv("NSQ_LOOKUP_PORT")
	cfg.Topic = os.Getenv("NSQ_TOPIC")
	cfg.Channel = os.Getenv("NSQ_CHANNEL")
	cfg.MaxAttempts = 10
	cfg.MaxInFlight = 5
	cfg.MaxRequeueDelay = time.Second * 900
	cfg.DefaultRequeueDelay = time.Second * 0

	pubsub.Sub(cfg, &message.MessageHandler{})
}
