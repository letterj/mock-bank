// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// App the structure used to hold pieces of the app
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize  the initialzation of the API
func (a *App) Initialize(v Config) {

	// Create database if necessary
	var msg string
	var err error
	a.DB, msg, err = CreateDB(v.DBLocation)
	if err != nil {
		log.Fatalf("Problem on %s with error %v", msg, err)
	}

	// Load in Currencies
	if err = LoadCurrencies(a.DB, v); err != nil {
		log.Fatalf("Problem with loading currencies: %v\n", err)
	}

	// Set Router
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run  Setup and Run the application itself
func (a *App) Run(addr string) {
	fmt.Println("Attempting to load of the FCFC Mock Bank API on port", addr[1:])
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	// Lists
	a.Router.HandleFunc("/api/v1/currency", a.getCurrencies).Methods("GET")
	a.Router.HandleFunc("/api/v1/account", a.getAccounts).Methods("GET")
	a.Router.HandleFunc("/api/v1/customer", a.getCustomers).Methods("GET")
	a.Router.HandleFunc("/api/v1/transaction", a.getTransactions).Methods("GET")
	a.Router.HandleFunc("/api/v1/notification", a.getNotifications).Methods("GET")
	// Specific Items
	a.Router.HandleFunc("/api/v1/account/{id}", a.getAccount).Methods("GET")
	a.Router.HandleFunc("/api/v1/customer/{lei}", a.getCustomer).Methods("GET")
	a.Router.HandleFunc("/api/v1/transaction/{id}", a.getTransaction).Methods("GET")
	a.Router.HandleFunc("/api/v1/notification/{id}", a.getNotification).Methods("GET")

	a.Router.HandleFunc("/api/v1/customer", a.addCustomer).Methods("POST")

	a.Router.HandleFunc("/api/v1/deposit", a.addDeposit).Methods("POST")

	// Account
	// a.Router.HandleFunc("/account", a.getAccounts).Methods("GET")
	// a.Router.HandleFunc("/account/{number}", a.getAccount).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// LoadCurrencies is used to load in currencies
func LoadCurrencies(db *sql.DB, v Config) error {
	var act account

	insertCurrency := `
	INSERT INTO currencies (currency_code, decimal_places, active_saturday, time_zone, recon_time)
	VALUES (?, ?, ?, ?, ?);`

	insertAccount := `
	INSERT INTO accounts (acct_number, qaccount, currency_code, balance)
	VALUES (?, ?, ?, ?);`

	for i := 0; i < len(v.Currencies); i++ {
		fmt.Printf("This Custodial Bank supports the %s.\n", v.Currencies[i].CurrencyCode)
		act.AcctNumber = fmt.Sprintf("%08d", rand.Intn(100000000)) + "-0" + strconv.Itoa(i)
		act.QuorumAccount = "0x"
		act.CurrencyCode = v.Currencies[i].CurrencyCode
		act.Balance = 0

		sqlStmt, err := db.Prepare(insertCurrency)
		if err != nil {
			return err
		}

		sqlStmt.Exec(v.Currencies[i].CurrencyCode, v.Currencies[i].DecimalPlaces,
			v.Currencies[i].ActiveSaturday, v.Currencies[i].CurrencyTimeZone,
			v.Currencies[i].ReconTime)

		sqlStmt, err = db.Prepare(insertAccount)
		if err != nil {
			return err
		}
		sqlStmt.Exec(act.AcctNumber, act.QuorumAccount, act.CurrencyCode,
			act.Balance)

	}
	return nil
}

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
   		qaccount TEXT NOT NULL 
 	);`

	accountTable := `
 	CREATE TABLE IF NOT EXISTS accounts(
   		acct_number      TEXT PRIMARY KEY,
   		qaccount         TEXT NOT NULL,
   		currency_code    TEXT NOT NULL,
   		balance          NUMERIC NOT NULL 
 	);`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions(
   		id               INTEGER PRIMARY KEY,
		type             TEXT NOT NULL,
		currency		 TEXT NOT NULL,   
   		trans_date       DATE NOT NULL DEFAULT (datetime('now', 'localtime')),
   		account_number   TEXT NOT NULL,
   		customer_id      INTEGER NOT NULL,
  		description      TEXT NOT NULL,
		amount           NUMERIC NOT NULL,
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
    	message          TEXT NOT NULL,
    	amount           NUMERIC NOT NULL,
   		status           TEXT NOT NULL DEFAULT "POSTED",
    	ack              BOOLEAN DEFAULT TRUE
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