package main

import (
	"encoding/json"
	"homework/services/message"
	"homework/services/pubsub"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	cfg := pubsub.NewPubConfig()
	cfg.Host = os.Getenv("NSQ_HOST")
	cfg.Port = os.Getenv("NSQ_PORT")
	cfg.Topic = os.Getenv("NSQ_TOPIC")

	msg := message.Message{
		Type:      "Deposit",
		Status:    "Created",
		Txis:      "5d6ftyguihojk-ppjhgfdsd76f8t7igyo-hpijo57689uj",
		Amount:    100,
		Timestamp: time.Now().String(),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	pubsub.Pub(cfg, payload)
}
