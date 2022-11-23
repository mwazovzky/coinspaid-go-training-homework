package main

import (
	"database/sql"
	"log"
	"os"

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

	getCurrencies(db)

	log.Println("All good")
}

func loadDBConfig(cfg *database.Config) {
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.Database = os.Getenv("DB_DATABASE")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
}

type Currency struct {
	ID  int
	ISO string
}

func getCurrencies(db *sql.DB) error {
	var currencies []Currency

	data, err := db.Query("SELECT id, iso FROM currencies")
	if err != nil {
		log.Println("Failed to select currencies", err)
		return err
	}

	var currency Currency

	for data.Next() {
		err := data.Scan(&currency.ID, &currency.ISO)
		if err != nil {
			log.Println("Failed to scan currency row:", err)
			return err
		}

		currencies = append(currencies, currency)
	}

	log.Println(currencies)

	return nil
}
