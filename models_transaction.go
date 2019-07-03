package main

import (
	"database/sql"
	"errors"
	"time"
)

type transaction struct {
	ID          int       `json:"id"`
	TransType   string    `json:"trans_type"`
	Currency    string    `json:"currency"`
	TransDate   time.Time `json:"trans_date"`
	AcctNumber  string    `json:"account_number"`
	CustomerID  int       `json:"customer_id"`
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
	Status      string    `json:"status"`
}

// Transactions
// ************
func getTransactions(db *sql.DB) ([]transaction, error) {
	sqlstmt := `
	SELECT id, type, trans_date, currency, account_number, 
		customer_id, description, amount
	FROM transactions`

	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trans := []transaction{}
	for rows.Next() {
		var tran transaction
		if err := rows.Scan(&tran.ID, &tran.TransType, &tran.TransDate,
			&tran.Currency, &tran.AcctNumber, &tran.CustomerID,
			&tran.Description, &tran.Amount); err != nil {
			return nil, err
		}
		trans = append(trans, tran)
	}

	return trans, nil
}

func (t *transaction) getTransaction(db *sql.DB) error {
	sqlStmt := `
	SELECT id, type, trans_date, currency, account_number, customer_id, 
		description, amount
	FROM transactions WHERE id=$1`
	return db.QueryRow(sqlStmt,
		t.ID).Scan(&t.TransType, &t.TransDate, &t.Currency, &t.AcctNumber, &t.CustomerID,
		&t.Description, &t.Amount)
}

func (t *transaction) createTransaction(db *sql.DB) (int, error) {
	insertTrans := `
	INSERT INTO transactions (type, currency, account_number, customer_id,
		description, amount, status)
		VALUES(?, ?, ?, ?, ?, ?, ?)`

	sqlStmt, err := db.Prepare(insertTrans)
	if err != nil {
		return 0, err
	}

	sqlStmt.Exec(&t.TransType, &t.Currency, &t.AcctNumber,
		&t.CustomerID, &t.Description, &t.Amount, &t.Status)

	rows, _ := db.Query("SELECT last_insert_rowid()")
	var transID int
	for rows.Next() {
		rows.Scan(&transID)
	}
	return transID, nil
}

func (t *transaction) updateTransaction(db *sql.DB) error {
	return errors.New("Not implemented")
}
