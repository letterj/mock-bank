package main

import (
	"database/sql"
	"time"
)

type notification struct {
	ID         int       `json:"id"`
	NoticeType string    `json:"type"`
	NoticeDate time.Time `json:"notice_date"`
	Currency   string    `json:"currency"`
	CustomerID int       `json:"customer_id"`
	TransID    int       `json:"transaction_id"`
	Message    string    `json:"message"`
	Amount     float32   `json:"amount"`
	QAccount   string    `json:"quorum_account"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	Rate       float32   `json:"rate"`
	Status     string    `json:"status"`
	Ack        bool      `json:"ack"`
}

// Notifications
// *************
func getNotifications(db *sql.DB) ([]notification, error) {

	sqlstmt := `
	SELECT id, type, notice_date, currency, customer_id, 
		transaction_id, message, amount, quorum_account,
		start_date, end_date, rate, status, ack
	FROM notifications WHERE ack = false`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []notification{}
	for rows.Next() {
		var note notification
		if err := rows.Scan(&note.ID, &note.NoticeType, &note.NoticeDate,
			&note.Currency, &note.CustomerID, &note.TransID,
			&note.Message, &note.Amount, &note.QAccount,
			&note.StartDate, &note.EndDate, &note.Rate, &note.Status,
			&note.Ack); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (n *notification) getNotification(db *sql.DB) error {
	sqlStmt := `
	SELECT id, type, notice_date, currency, customer_id, transaction_id,
	message, amount, quorum_account, start_date, end_date, rate, status, ack
	FROM notifications WHERE id=$1`
	return db.QueryRow(sqlStmt,
		n.ID).Scan(&n.NoticeDate, &n.NoticeType, &n.Currency,
		&n.CustomerID, &n.TransID, &n.Message, &n.Amount,
		&n.QAccount, &n.StartDate, &n.EndDate,
		&n.Rate, &n.Status, &n.Ack)
}

func (n *notification) createNotification(db *sql.DB) error {
	insertNotice := `
	INSERT INTO notifications (type, customer_id,
		transaction_id, message, amount, currency, quorum_account,
		start_date, end_date, rate, ack)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlStmt, err := db.Prepare(insertNotice)
	if err != nil {
		return err
	}

	sqlStmt.Exec(&n.NoticeType, &n.CustomerID, &n.TransID,
		&n.Message, &n.Amount, &n.Currency, &n.QAccount,
		&n.StartDate, &n.EndDate, &n.Rate, &n.Ack)

	return nil
}

func (n *notification) updateNotification(db *sql.DB) error {
	updateNotice := `
	UPDATE notifications
	SET ack = true
	WHERE id = ? AND ack = false`

	sqlStmt, err := db.Prepare(updateNotice)
	if err != nil {
		return err
	}

	sqlStmt.Exec(&n.ID)

	lookupNotice := `
	SELECT id, type, notice_date, currency, customer_id, transaction_id,
	message, amount, quorum_account, start_date, end_date, rate, status, ack
	FROM notifications WHERE id=$1`

	err = db.QueryRow(lookupNotice,
		n.ID).Scan(&n.NoticeDate, &n.NoticeType, &n.Currency,
		&n.CustomerID, &n.TransID, &n.Message, &n.Amount,
		&n.QAccount, &n.StartDate, &n.EndDate,
		&n.Rate, &n.Status)

	if err != nil {
		return err
	}

	return nil
}
