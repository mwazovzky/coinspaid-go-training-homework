package address

import (
	"database/sql"
	"log"
)

type AddressRepository struct {
	db *sql.DB
}

type Address struct {
	ID       int
	Hash     string
	Currency string
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{db}
}

func (ar *AddressRepository) List(addresses *[]Address) error {
	data, err := ar.db.Query("SELECT a.hash, c.iso FROM addresses a LEFT JOIN currencies c ON a.currency_id=c.id")
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

func (ar *AddressRepository) Exist(address string, iso string) bool {
	query := `
		SELECT COUNT(*) FROM addresses 
		LEFT JOIN currencies ON currencies.id=addresses.currency_id
		WHERE addresses.hash=? 
		AND currencies.iso=?
	`

	data := ar.db.QueryRow(query, address, iso)

	var count int
	err := data.Scan(&count)
	if err != nil {
		log.Println("Failed to get addresses count", err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
}

func (ar *AddressRepository) Get(address *Address, hash string, iso string) error {
	query := `
		SELECT addresses.id, addresses.hash, currencies.iso FROM addresses 
		LEFT JOIN currencies ON currencies.id=addresses.currency_id
		WHERE addresses.hash=? 
		AND currencies.iso=?
	`

	data := ar.db.QueryRow(query, hash, iso)

	err := data.Scan(&address.ID, &address.Hash, &address.Currency)
	if err != nil {
		log.Println("Failed to get addresses count", err)
		return err
	}

	return nil
}
