CREATE TABLE IF NOT EXISTS currencies(
	currency_code     TEXT PRIMARY KEY,
	decimal_places    INTEGER NOT NULL,
	active_saturday   BOOLEAN DEFAULT FALSE,
	time_zone	   	  TEXT NOT NULL,
	recon_time        DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts(
	acct_number      TEXT PRIMARY KEY,
	qaccount         TEXT NOT NULL,
	currency_code    TEXT NOT NULL,
	balance          NUMERIC NOT NULL 
);

CREATE UNIQUE INDEX accounts_qaccount ON accounts (qaccount);
CREATE UNIQUE INDEX accounts_currency ON accounts (currency_code);


CREATE TABLE IF NOT EXISTS customers(
	id       SERIAL PRIMARY KEY,
	lei      TEXT NOT NULL,
	name     TEXT NOT NULL,
	qaccount TEXT NOT NULL 
);

CREATE UNIQUE INDEX customers_lei ON customers (lei);
CREATE UNIQUE INDEX customers_qaccount ON customers (qaccount);

CREATE TABLE IF NOT EXISTS transactions(
   	id               SERIAL PRIMARY KEY,
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
	status           TEXT NOT NULL DEFAULT "PENDING"
);

CREATE TABLE IF NOT EXISTS notifications(
	id               SERIAL PRIMARY KEY,
   	type             TEXT NOT NULL,
	notice_date      TIMESTAMPTZ NOT NULL DEFAULT (datetime('now', 'localtime')),
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
);