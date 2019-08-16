package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

// LoadConfig will load configuation variables
func LoadConfig(fileName string) (Config, error) {
	var c Config
	var m Currency

	if fileName == "" || fileName == "default" {

		c.Port = 8080
		c.DBLocation = "./data/bank.db"
		c.TLS = false

		m.CurrencyCode = "USD"
		m.DecimalPlaces = 2
		m.ActiveSaturday = false
		m.CurrencyTimeZone = "NYC"
		m.ReconTime = "16:00"
		c.Currencies = append(c.Currencies, m)

		m.CurrencyCode = "EUR"
		m.DecimalPlaces = 2
		m.ActiveSaturday = false
		m.CurrencyTimeZone = "NYC"
		m.ReconTime = "21:00"
		c.Currencies = append(c.Currencies, m)

		return c, nil
	}

	// Open, Read File and Unmarshall it into json
	cFile, err := os.Open(fileName)
	if err != nil {
		return c, err
	}
	defer cFile.Close()
	jsonParser := json.NewDecoder(cFile)
	if err = jsonParser.Decode(&c); err != nil {
		return c, err
	}

	// TODO:   Log config file

	return c, nil
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
		log.Printf("This Custodial Bank supports the %s.\n", v.Currencies[i].CurrencyCode)
		act.AcctNumber = fmt.Sprintf("%12d", rand.Intn(1000000000000)) + "-0" + strconv.Itoa(i)
		act.QAccount = "0x1111111111111111111111111111111" + strconv.Itoa(i)
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
		sqlStmt.Exec(act.AcctNumber, act.QAccount, act.CurrencyCode,
			act.Balance)

	}
	return nil
}
