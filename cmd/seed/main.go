package main

import (
	"log"
	"os"

	"homework/models/address"
	"homework/models/message"
	"homework/models/transaction"
	"homework/services/database"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	cfg := database.NewConfig()
	loadDBConfig(cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	defer database.Close(db)

	msg := message.MakeFakeMessage(db)

	var adr address.Address
	addressRepository := address.NewAddressRepository(db)
	err = addressRepository.Get(&adr, msg.Address, msg.Currency)
	if err != nil {
		log.Println("Error: find address", err)
		return
	}

	txRepository := transaction.NewTransactionRepository(db)
	err = txRepository.Create(msg.Type, msg.Status, adr.ID, msg.Txid, msg.Amount)
	if err != nil {
		log.Println("Error: create address", err)
		return
	}

	err = txRepository.Update(msg.Txid, adr.ID, "Ignored")
	if err != nil {
		log.Println("Error: create address", err)
		return
	}

	log.Println("All good")
}

func loadDBConfig(cfg *database.Config) {
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.Database = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
}
