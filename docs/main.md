
type currency struct {
	CurrencyCode   string    `json:"currency"`
	DecimalPlaces  int       `json:"decimal_places"`
	ActiveSaturday bool      `json:"active_saturday"`
	ReconTime      time.Time `json:"recon"`
}

func main() {
	dbLocation := "/tmp/fcfcbank/bank.db"

	currencyTable := `
	CREATE TABLE IF NOT EXISTS currencies(
   		currency_code     TEXT NOT NULL,
   		decimal_places    INTEGER NOT NULL,
   		active_saturday   BOOLEAN DEFAULT FALSE,
   		recon_time        DATE NOT NULL
 	);`

	insertCurrency := `
	INSERT INTO currencies (currency_code, decimal_places, active_saturday, recon_time) 
	VALUES (?, ?, ?, ?);`

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
    	type             TEXT NOT NULL,
    	notice_date      DATE NOT NULL,
    	account_id       INTEGER NOT NULL,
   		customer_id      INTEGER NOT NULL,
    	transaction_id   INTEGER NOT NULL,
    	message          TEXT NOT NULL,
    	amount           NUMERIC NOT NULL,
   		status           TEXT NOT NULL,
    	ack              BOOLEAN
	);`

	bankDB, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		panic(err)
	}
	defer bankDB.Close()

	// Create Tables
	sqlstmt, _ := bankDB.Prepare(currencyTable)
	sqlstmt.Exec()
	sqlstmt, _ = bankDB.Prepare(customerTable)
	sqlstmt.Exec()
	sqlstmt, _ = bankDB.Prepare(accountTable)
	sqlstmt.Exec()
	sqlstmt, _ = bankDB.Prepare(transactionTable)
	sqlstmt.Exec()
	sqlstmt, _ = bankDB.Prepare(notificationTable)
	sqlstmt.Exec()

	// Insert currencies into the Currency Table
	sqlstmt, _ = bankDB.Prepare(insertCurrency)
	sqlstmt.Exec("USD", 2, false, time.Now())
	sqlstmt, _ = bankDB.Prepare(insertCurrency)
	sqlstmt.Exec("CAD", 2, false, time.Now())

	// Query the Currency Table
	queryStmt := `
		SELECT currency_code, decimal_places, active_saturday, recon_time
		FROM currencies`
	rows, err := bankDB.Query(queryStmt)
	if err != nil {
		panic(err)
	}
	var c currency
	for rows.Next() {
		rows.Scan(&c.CurrencyCode, &c.DecimalPlaces, &c.ActiveSaturday, &c.ReconTime)
		fmt.Println(c)
	}
}