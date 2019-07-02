
/*

currencyTable := `
CREATE TABLE IF NOT EXISTS currencies(
   currency_code     TEXT NOT NULL,
   decimal_places     INTEGER NOT NULL,
   active_saturday   BOOLEAN DEFAULT FALSE,
   recon_time        DATE NOT NULL
 );`

insertUSD ;= `
INSERT INTO currencies (currency_code, decimal_places, active_saturday, recon_time)
	 VALUES ('USD', 2, FALSE, DateTime('now'));`
insertCAD ;= `
INSERT INTO currencies (currency_code, decimal_places, active_saturday, recon_time)
	VALUES ('CAD', 2, FALSE, DateTime('now'));`

customerTable := `
 CREATE TABLE IF NOT EXISTS customers(
   id       INTEGER PRIMARY KEY,
   lei      TEXT NOT NULL,
   name     TEXT NOT NULL,
   qaccount TEXT NOT NULL
 );`

accountTable := `
 CREATE TABLE IF NOT EXISTS accounts(
   id               INTEGER PRIMARY KEY,
   acct_number      TEXT NOT NULL,
   qaccount         TEXT NOT NULL,
   currency_code    TEXT NOT NULL,
   balance          NUMERIC NOT NULL
 );`

transactionTable := `
CREATE TABLE IF NOT EXISTS transactions(
   id               INTEGER PRIMARY KEY,
   type             TEXT NOT NULL,
   trans_date       DATE NOT NULL,
   account_number   INTEGER NOT NULL,
   customer_id      INTEGER NOT NULL,
   description      TEXT NOT NULL,
   amount           NUMERIC NOT NULL
 );`

notificationTable := `
CREATE TABLE IF NOT EXISTS notifications(
   id               INTEGER PRIMARY KEY,
   type				TEXT NOT NULL,
   notice_date      DATE NOT NULL,
   account_id       INTEGER NOT NULL,
   customer_id      INTEGER NOT NULL,
   transaction_id   INTEGER NOT NULL,
   message          TEXT NOT NULL,
   amount           NUMERIC NOT NULL,
   status           TEXT NOT NULL,
   ack              BOOLEAN
);`
*/