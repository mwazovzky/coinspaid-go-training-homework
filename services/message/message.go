package message

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

type Message struct {
	Type      string
	Status    string
	Currency  string
	Address   string
	Txid      string
	Amount    int
	Timestamp string
}

type MessageHandler struct {
	deposits    chan Message
	withdrawals chan Message
}

const bufferSize = 10
const typeDeposit = "Deposit"
const typeWithdrawal = "Withdrawal"

func NewMessageHandler() *MessageHandler {
	deposits := make(chan Message, bufferSize)
	withdrawals := make(chan Message, bufferSize)

	go func() {
		for msg := range deposits {
			log.Println("process deposit, buffer:", len(deposits))
			processDeposit(msg)
		}
	}()

	go func() {
		for msg := range withdrawals {
			log.Println("process withdrawals, buffer:", len(deposits))
			processWithdrawals(msg)
		}
	}()

	return &MessageHandler{deposits, withdrawals}
}

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
// Message with an empty body is ignored/discarded
func (mh *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		log.Println("Error: empty message body")
		return nil
	}

	var msg Message

	// unmarshall messages
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		log.Println("Error: unmarshaling message body", err)
		return nil
	}

	// validate message, otherwise FIN
	// check if address exists, otherwise FIN
	// check if transaction exists if message type is not created, otherwise REQ

	// push message to buffer
	switch msg.Type {
	case typeDeposit:
		mh.deposits <- msg
	case typeWithdrawal:
		mh.withdrawals <- msg
	default:
		log.Println("Error: upprocessable message type")
	}

	return nil
}

func processDeposit(msg Message) {
	// create event record
	// process event create or update transaction status
	log.Println("----- Deposit ------")
	log.Println("Type : ", msg.Type)
	log.Println("Status : ", msg.Status)
	log.Println("Txid : ", msg.Txid)
	log.Println("Currency : ", msg.Currency)
	log.Println("Address : ", msg.Address)
	log.Println("Amount : ", msg.Amount)
	log.Println("Timestamp : ", msg.Timestamp)
	log.Println("--------------------")

	time.Sleep(time.Millisecond * 300)
}

func processWithdrawals(msg Message) {
	// create event record
	// process event create or update transaction status
	log.Println("----- Withdrawal ------")
	log.Println("Type : ", msg.Type)
	log.Println("Status : ", msg.Status)
	log.Println("Txid : ", msg.Txid)
	log.Println("Currency : ", msg.Currency)
	log.Println("Address : ", msg.Address)
	log.Println("Amount : ", msg.Amount)
	log.Println("Timestamp : ", msg.Timestamp)
	log.Println("--------------------")

	time.Sleep(time.Millisecond * 100)
}
