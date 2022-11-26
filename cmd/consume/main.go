package main

import (
	"homework/services/database"
	"homework/services/messagehandler"
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
	dbcfg := database.NewConfig()
	loadDBConfig(dbcfg)

	db, err := database.Connect(dbcfg)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	defer database.Close(db)

	host := os.Getenv("NSQ_LOOKUP_HOST")
	port := os.Getenv("NSQ_LOOKUP_PORT")
	topic := os.Getenv("NSQ_TOPIC")
	channel := os.Getenv("NSQ_CHANNEL")

	nsqcfg := nsq.NewConfig()
	nsqcfg.MaxAttempts = 10
	nsqcfg.MaxInFlight = 5
	nsqcfg.MaxRequeueDelay = time.Second * 900
	nsqcfg.DefaultRequeueDelay = time.Second * 0

	consumer, err := pubsub.NewConsumer(nsqcfg, topic, channel)
	if err != nil {
		log.Fatal(err)
	}

	// Register message handler
	handler := messagehandler.NewMessageHandler(db)
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

func loadDBConfig(cfg *database.Config) {
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.Database = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
}
