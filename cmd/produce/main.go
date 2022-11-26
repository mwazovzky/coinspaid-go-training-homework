package main

import (
	"database/sql"
	"encoding/json"
	"homework/services/database"
	"homework/services/message"
	"homework/services/pubsub"
	"log"
	"math/rand"
	"os"
	"time"

	faker "github.com/bxcodec/faker/v3"
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
)

func init() {
	godotenv.Load()
	rand.Seed(time.Now().UnixNano())
}

type Address struct {
	Hash     string
	Currency string
}

// count is number of transactions to generate
const count = 50

func main() {
	db := connectDB()
	defer database.Close(db)

	producer := getProducer()
	defer producer.Stop()

	topic := os.Getenv("NSQ_TOPIC")

	var transactions []message.Message
	generateTransactions(db, count, &transactions)

	var messages []message.Message
	generateMessages(&transactions, &messages)

	shuffleMessages(&messages)

	for _, msg := range messages {
		publishMessage(msg, topic, producer)
	}
}

func connectDB() *sql.DB {
	cfg := database.NewConfig()
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.Database = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	return db
}

func getProducer() *nsq.Producer {
	host := os.Getenv("NSQ_HOST")
	port := os.Getenv("NSQ_PORT")
	cfg := pubsub.NewProducerConfig(host, port)

	producer, err := pubsub.NewProducer(cfg)
	if err != nil {
		log.Fatal("Could create nsq producer", err)
	}

	return producer
}

// generateTransactions generates batch of fake transactions
func generateTransactions(db *sql.DB, count int, operations *[]message.Message) {
	var addresses []Address

	getAddresses(db, &addresses)

	for i := 0; i < count; i++ {
		*operations = append(*operations, createTransaction(addresses))
	}
}

// getAddresses gets addresses from database
func getAddresses(db *sql.DB, addresses *[]Address) error {
	data, err := db.Query("SELECT a.hash, c.iso FROM addresses a LEFT JOIN currencies c ON a.currency_id=c.id")
	if err != nil {
		log.Println("Failed to select addresses", err)
		return err
	}

	var address Address

	for data.Next() {
		err := data.Scan(&address.Hash, &address.Currency)
		if err != nil {
			log.Println("Failed to scan currency row:", err)
			return err
		}

		*addresses = append(*addresses, address)
	}

	return nil
}

// createTransaction creates fake message in final status - Cancelled of Confirmed
func createTransaction(addresses []Address) message.Message {
	types := []string{"Deposit", "Withdrawal"}
	statuses := []string{"Cancelled", "Confirmed"}
	address := addresses[rand.Intn(len(addresses))]

	return message.Message{
		Type:      types[rand.Intn(2)],
		Status:    statuses[rand.Intn(2)],
		Currency:  address.Currency,
		Address:   address.Hash,
		Txid:      faker.UUIDHyphenated(),
		Amount:    rand.Intn(1000),
		Timestamp: time.Now().String(),
	}
}

// generateMessages creates initial messages with status Created for every fake message
func generateMessages(operations *[]message.Message, messages *[]message.Message) {
	for _, operation := range *operations {
		*messages = append(*messages, operation)
		operation.Status = "Created"
		*messages = append(*messages, operation)
	}
}

func shuffleMessages(messages *[]message.Message) {
	rand.Shuffle(len(*messages), func(i, j int) { (*messages)[i], (*messages)[j] = (*messages)[j], (*messages)[i] })
}

// publishMessage publishes message to specified nsq topic
func publishMessage(msg message.Message, topic string, producer *nsq.Producer) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	err = producer.Publish(topic, payload)
	if err != nil {
		log.Println("Could not connect", err)
		return err
	}

	return nil
}
