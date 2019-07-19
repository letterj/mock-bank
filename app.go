// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	log.Printf("The location of the bank database '%s'\n", v.DBLocation)
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
func (a *App) Run(addr string, tls bool) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)

	if tls {
		log.Println("TLS set to TRUE")
		log.Fatal(http.ListenAndServeTLS(addr, "cert.pem", "key.pem", loggedRouter))
	} else {
		log.Println("TLS set to FALSE")
		log.Fatal(http.ListenAndServe(addr, loggedRouter))
	}

	log.Println("Loaded of the FCFC Mock Bank API on port", addr[1:])
}

func (a *App) initializeRoutes() {
	// Lists
	a.Router.HandleFunc("/api/v1/version", a.version).Methods("GET")
	a.Router.HandleFunc("/api/v1/currency", a.getCurrencies).Methods("GET")
	a.Router.HandleFunc("/api/v1/account", a.getAccounts).Methods("GET")
	a.Router.HandleFunc("/api/v1/customer", a.getCustomers).Methods("GET")
	a.Router.HandleFunc("/api/v1/transaction", a.getTransactions).Methods("GET")
	a.Router.HandleFunc("/api/v1/notification", a.getNotifications).Methods("GET")

	// Specific Items
	a.Router.HandleFunc("/api/v1/account/{currency}", a.getAccount).Methods("GET")
	a.Router.HandleFunc("/api/v1/customer/{lei}", a.getCustomer).Methods("GET")
	a.Router.HandleFunc("/api/v1/transaction/{id}", a.getTransaction).Methods("GET")
	a.Router.HandleFunc("/api/v1/notification/{id}", a.getNotification).Methods("GET")

	// Creates
	a.Router.HandleFunc("/api/v1/customer", a.addCustomer).Methods("POST")
	a.Router.HandleFunc("/api/v1/deposit", a.addDeposit).Methods("POST")
	a.Router.HandleFunc("/api/v1/withdraw", a.addWithdraw).Methods("POST")
	// Updates
	a.Router.HandleFunc("/api/v1/account/{acctNumber}", a.updateAccount).Methods("PUT")
	a.Router.HandleFunc("/api/v1/notification/{id}", a.updateNotification).Methods("PUT")

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
