package message

import (
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
)

type Message struct {
	Type      string
	Status    string
	Txis      string
	Amount    int
	Timestamp string
}

type MessageHandler struct{}

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Message with an empty body is simply ignored/discarded.
		return nil
	}

	return processMessage(m.Body)
}

func processMessage(body []byte) error {
	var msg Message

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		return err
	}

	log.Println("----- Message ------")
	log.Println("Type : ", msg.Type)
	log.Println("Status : ", msg.Status)
	log.Println("Txis : ", msg.Txis)
	log.Println("Amount : ", msg.Amount)
	log.Println("Timestamp : ", msg.Timestamp)
	log.Println("--------------------")

	return nil
}
