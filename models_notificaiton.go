package main

import (
	"database/sql"
	"errors"
	"time"
)

type notification struct {
	ID         int       `json:"id"`
	NoticeDate time.Time `json:"notice_date"`
	AccountID  int       `json:"account_id"`
	CustomerID int       `json:"customer_id"`
	TransID    int       `json:"transaction_id"`
	Message    string    `json:"message"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	Ack        bool      `json:"ack"`
}

// Notifications
// *************
func getNotifications(db *sql.DB) ([]notification, error) {

	sqlstmt := `
	SELECT id, notice_date, account_id, customer_id, transaction_id,
		message, amount, status, ack
	FROM notifications`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []notification{}
	for rows.Next() {
		var note notification
		if err := rows.Scan(&note.ID, &note.NoticeDate, &note.AccountID,
			&note.CustomerID, &note.TransID, &note.Message,
			&note.Amount, &note.Status, &note.Ack); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (n *notification) getNotification(db *sql.DB) error {
	sqlStmt := `
	SELECT id, notice_date, account_id, customer_id, transaction_id,
	message, amount, status, ack
	FROM notifications WHERE id=$1`
	return db.QueryRow(sqlStmt,
		n.ID).Scan(&n.NoticeDate, &n.AccountID, &n.CustomerID, &n.TransID,
		&n.Message, &n.Amount, &n.Status, &n.Ack)
}

func (n *notification) createNotification(db *sql.DB) error {
	return errors.New("Not implemented")
}
