package message

import (
	"database/sql"
	"homework/models/address"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
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

type MessageRepository struct {
	db *sql.DB
}

const TypeDeposit = "Deposit"
const TypeWithdrawal = "Withdrawal"
const StatusCreated = "Created"
const StatusCancelled = "Cancelled"
const StatusConfirmed = "Confirmed"
const DateFormat = "2006-01-02 15:04:05"

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (mr *MessageRepository) Create(msg Message) error {
	cmd := `
		INSERT INTO messages (type, status, currency, address, txid, amount, timestamp) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := mr.db.Exec(cmd, msg.Type, msg.Status, msg.Currency, msg.Address, msg.Txid, msg.Amount, msg.Timestamp)

	return err
}

// MakeFakeMessage makes fake message in final status - Cancelled of Confirmed
func MakeFakeMessage(db *sql.DB) Message {
	addressRepository := address.NewAddressRepository(db)
	var addresses []address.Address
	addressRepository.List(&addresses)

	types := []string{TypeDeposit, TypeWithdrawal}
	statuses := []string{StatusCancelled, StatusConfirmed}
	address := addresses[rand.Intn(len(addresses))]

	return Message{
		Type:      types[rand.Intn(2)],
		Status:    statuses[rand.Intn(2)],
		Currency:  address.Currency,
		Address:   address.Hash,
		Txid:      faker.UUIDHyphenated(),
		Amount:    rand.Intn(1000),
		Timestamp: time.Now().Format(DateFormat),
	}
}
