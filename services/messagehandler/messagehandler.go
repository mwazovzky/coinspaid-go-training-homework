package messagehandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"homework/models/address"
	"homework/models/message"
	"homework/models/transaction"
	"log"

	"github.com/nsqio/go-nsq"
)

type MessageHandler struct {
	db          *sql.DB
	deposits    chan message.Message
	withdrawals chan message.Message
}

const bufferSize = 10

func NewMessageHandler(db *sql.DB) *MessageHandler {
	deposits := make(chan message.Message, bufferSize)
	withdrawals := make(chan message.Message, bufferSize)

	go func() {
		for msg := range deposits {
			process(db, msg)
		}
	}()

	go func() {
		for msg := range withdrawals {
			process(db, msg)
		}
	}()

	return &MessageHandler{db, deposits, withdrawals}
}

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
// Message with an empty body is ignored/discarded
func (mh *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		log.Println("Error: empty message body")
		return nil
	}

	var msg message.Message

	// unmarshall messages
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		log.Println("Error: unmarshaling message body", err)
		return nil
	}

	// validate message, otherwise FIN

	// check if address exists, otherwise FIN
	addressRepository := address.NewAddressRepository(mh.db)
	exist := addressRepository.Exist(msg.Address, msg.Currency)
	if !exist {
		log.Println("Error: unexpected address", msg.Address)
		return nil
	}

	// check if transaction exists if message type is not created, otherwise REQ
	if msg.Status != message.StatusCreated {
		txRepository := transaction.NewTransactionRepository(mh.db)
		exist := txRepository.Exist(msg.Txid, msg.Address, msg.Currency)
		if !exist {
			log.Println("Error: unexpected status", msg.Status)
			m.RequeueWithoutBackoff(0)
			return nil
		}
	}

	switch msg.Type {
	case message.TypeDeposit:
		mh.deposits <- msg
	case message.TypeWithdrawal:
		mh.withdrawals <- msg
	default:
		log.Println("Error: upprocessable message type")
	}

	return nil
}

func process(db *sql.DB, msg message.Message) {
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

func requeue(msg message.Message) error {
	txt := fmt.Sprintf("Error: unexpected status, %s", msg.Status)
	log.Println(txt)
	return errors.New(txt)
}
