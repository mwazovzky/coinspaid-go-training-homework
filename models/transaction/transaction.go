package transaction

import (
	"database/sql"
	"log"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (tr *TransactionRepository) Exist(txid string, address string, iso string) bool {
	query := `
		SELECT COUNT(*) FROM transactions 
		LEFT JOIN addresses ON addresses.id=transactions.address_id
		LEFT JOIN currencies ON currencies.id=addresses.currency_id
		WHERE transactions.txid=? 
		AND addresses.hash=?
		AND currencies.iso=?
	`

	data := tr.db.QueryRow(query, txid, address, iso)

	var count int
	err := data.Scan(&count)
	if err != nil {
		log.Println("Failed to get transactions count", err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
}

func (tr *TransactionRepository) Create(txType string, txStatus string, addressID int, txid string, amount int) error {
	cmd := `
		INSERT INTO transactions (type, status, address_id, txid, amount) 
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := tr.db.Exec(cmd, txType, txStatus, addressID, txid, amount)

	log.Println("tx.Create", txid, txType, txStatus, addressID, amount, "error:", err)

	return err
}

func (tr *TransactionRepository) Update(txid string, addressID int, status string) error {
	cmd := `
		UPDATE transactions SET status=? 
		WHERE txid=?
		AND address_id=?
	`

	_, err := tr.db.Exec(cmd, status, txid, addressID)

	log.Println("tx.Update", txid, status, addressID, "error:", err)

	return err
}
