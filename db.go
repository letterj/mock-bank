package main

import "database/sql"

// CreateDB sets up the database
func CreateDB(location string) (*sql.DB, string, error) {
	currencyTable := `

	CREATE TABLE IF NOT EXISTS currencies(
   		currency_code     TEXT PRIMARY KEY,
   		decimal_places    INTEGER NOT NULL,
   		active_saturday   BOOLEAN DEFAULT FALSE,
   		time_zone	   	  TEXT NOT NULL,
   		recon_time        DATE NOT NULL
 	);`

	customerTable := `
 	CREATE TABLE IF NOT EXISTS customers(
   		id       INTEGER PRIMARY KEY,
   		lei      TEXT NOT NULL,
   		name     TEXT NOT NULL,
		qaccount TEXT NOT NULL,
		UNIQUE (lei),
		UNIQUE (qaccount)    
 	);`

	accountTable := `
 	CREATE TABLE IF NOT EXISTS accounts(
   		acct_number      TEXT PRIMARY KEY,
   		qaccount         TEXT NOT NULL,
   		currency_code    TEXT NOT NULL,
		balance          NUMERIC NOT NULL,
		UNIQUE(qaccount),
		UNIQUE(currency_code)   
 	);`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions(
   		id               INTEGER PRIMARY KEY,
		type             TEXT NOT NULL,
		currency		 TEXT NOT NULL,   
   		trans_date       DATE NOT NULL DEFAULT (datetime('now', 'localtime')),
		account_number   TEXT NOT NULL,
		quorum_account   TEXT,  
   		customer_id      INTEGER NOT NULL,
  		description      TEXT NOT NULL,
		amount           NUMERIC NOT NULL,
		start_date       TEXT,
		end_date         TEXT,
		rate             NUMBERIC,
		status           NUMERIC NOT NULL DEFAULT "PENDING"
 	);`

	notificationTable := `
	CREATE TABLE IF NOT EXISTS notifications(
		id               INTEGER PRIMARY KEY,
    	type             TEXT NOT NULL,
		notice_date      DATE NOT NULL DEFAULT (datetime('now', 'localtime')),
		currency		 TEXT NOT NULL,
   		customer_id      INTEGER NOT NULL,
		transaction_id   INTEGER NOT NULL,
		quorum_account   TEXT,
    	message          TEXT NOT NULL,
		amount           NUMERIC NOT NULL,
		start_date       TEXT,
		end_date 		 TEXT, 
		rate             NUMBERIC,
   		status           TEXT NOT NULL DEFAULT "POSTED",
    	ack              BOOLEAN DEFAULT FALSE
	);`

	bankDB, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, "Connection", err
	}

	// Create Tables
	sqlstmt, err := bankDB.Prepare(currencyTable)
	if err != nil {
		return nil, "Prep Currency Table", err
	}
	sqlstmt.Exec()
	sqlstmt, err = bankDB.Prepare(customerTable)
	if err != nil {
		return nil, "Customer Table", err
	}
	sqlstmt.Exec()
	sqlstmt, err = bankDB.Prepare(accountTable)
	if err != nil {
		return nil, "Account Table", err
	}
	sqlstmt.Exec()
	sqlstmt, err = bankDB.Prepare(transactionTable)
	if err != nil {
		return nil, "Transaction Table", err
	}
	sqlstmt.Exec()
	sqlstmt, err = bankDB.Prepare(notificationTable)
	if err != nil {
		return nil, "Notification Table", err
	}
	sqlstmt.Exec()

	return bankDB, "", nil
}
