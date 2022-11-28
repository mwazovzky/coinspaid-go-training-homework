package messagehandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"homework/models/address"
	"homework/models/message"
	"homework/models/transaction"
	"log"
	"runtime"

	"github.com/nsqio/go-nsq"
)

type MessageHandler struct {
	db     *sql.DB
	buffer chan message.Message
}

const bufferSize = 100

func NewMessageHandler(db *sql.DB) *MessageHandler {
	buffer := make(chan message.Message, bufferSize)

	// g := 1
	g := runtime.GOMAXPROCS(0)

	for i := 0; i < g; i++ {
		go func() {
			for msg := range buffer {
				fmt.Println("buffer size", len(buffer))
				processMessage(db, msg)
			}
		}()
	}

	return &MessageHandler{db, buffer}
}

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
// Check nsq config for processing REQ, currenctly there us a delay id Error return
func (mh *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		log.Println("Error: empty message body")
		return nil
	}

	var msg message.Message
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		log.Println("Error: unmarshaling message body", err)
		return nil
	}

	// validate message, otherwise FIN

	// check if address exists, otherwise FIN
	if !existAddress(mh.db, msg) {
		log.Println("Error: unexpected address", msg.Currency, msg.Address)
		return nil
	}

	// ignore duplicate messages, if transaction already exists and tx.Status == msg.status retun nil

	// if message status is not created check if transaction exists, otherwise wait for create message - REQ
	if msg.Status != message.StatusCreated && !existTransaction(mh.db, msg) {
		log.Println("Error: unexpected status", msg.Status)
		m.RequeueWithoutBackoff(0)
		return nil
	}

	mh.buffer <- msg

	return nil
}

func processMessage(db *sql.DB, msg message.Message) {
	createMessage(db, msg)

	switch msg.Status {
	case message.StatusCreated:
		createTransaction(db, msg)
	case message.StatusCancelled:
		cancelTransaction(db, msg)
	case message.StatusConfirmed:
		confirmTransaction(db, msg)
	}
}

func createMessage(db *sql.DB, msg message.Message) error {
	messageRepository := message.NewMessageRepository(db)
	err := messageRepository.Create(msg)

	return err
}

func existAddress(db *sql.DB, msg message.Message) bool {
	addressRepository := address.NewAddressRepository(db)
	exist := addressRepository.Exist(msg.Address, msg.Currency)

	return exist
}

func existTransaction(db *sql.DB, msg message.Message) bool {
	txRepository := transaction.NewTransactionRepository(db)
	exist := txRepository.Exist(msg.Txid, msg.Address, msg.Currency)

	return exist
}

func createTransaction(db *sql.DB, msg message.Message) error {
	var address address.Address
	getAddress(db, msg, &address)

	txRepository := transaction.NewTransactionRepository(db)
	err := txRepository.Create(msg.Type, msg.Status, address.ID, msg.Txid, msg.Amount)

	return err
}

func confirmTransaction(db *sql.DB, msg message.Message) error {
	return updateTransaction(db, msg)
}

func cancelTransaction(db *sql.DB, msg message.Message) error {
	return updateTransaction(db, msg)
}

func updateTransaction(db *sql.DB, msg message.Message) error {
	var err error

	var address address.Address
	err = getAddress(db, msg, &address)

	txRepository := transaction.NewTransactionRepository(db)
	err = txRepository.Update(msg.Txid, address.ID, msg.Status)

	return err
}

func getAddress(db *sql.DB, msg message.Message, adr *address.Address) error {
	addressRepository := address.NewAddressRepository(db)
	err := addressRepository.Get(adr, msg.Address, msg.Currency)

	return err
}
