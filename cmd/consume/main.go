package main

import (
	"encoding/json"
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

type messageHandler struct{}

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
	consumer.AddHandler(&messageHandler{})

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

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
// Message with an empty body is ignored/discarded
func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	return processMessage(m.Body)
}

func processMessage(body []byte) error {
	var msg message.Message

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		return err
	}

	log.Println("----- Message ------")
	log.Println("Type : ", msg.Type)
	log.Println("Status : ", msg.Status)
	log.Println("Txis : ", msg.Txis)
	log.Println("Currency : ", msg.Currency)
	log.Println("Address : ", msg.Address)
	log.Println("Amount : ", msg.Amount)
	log.Println("Timestamp : ", msg.Timestamp)
	log.Println("--------------------")

	return nil
}
