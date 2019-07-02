package main

import (
	"database/sql"
	"errors"
	"fmt"
)

// Currency used for config file
type Currency struct {
	CurrencyCode     string `json:"currency"`
	DecimalPlaces    int    `json:"decimal_places"`
	ActiveSaturday   bool   `json:"active_saturday"`
	CurrencyTimeZone string `json:"time_zone"`
	ReconTime        string `json:"recon_time"`
}

type deposit struct {
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	QuorumAccount string  `json:"quorum_account"`
	Currency      string  `json:"currency_code"`
	Amount        float32 `json:"amount"`
	RefID         string  `json:"refid"`
}

type withdraw struct {
	Lei             string  `json:"lei"`
	CurrencyCode    string  `json:"currency_code"`
	Amount          float32 `json:"amount"`
	HomeBankName    string  `json:"home_bank_name"`
	WireInstruction string  `json:"wire_instructions"`
}

// Currencies
// **********
func getCurrencies(db *sql.DB) ([]Currency, error) {
	sqlstmt := `
		SELECT currency_code, decimal_places, active_saturday,
			time_zone, recon_time
		FROM currencies`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	curs := []Currency{}
	for rows.Next() {
		var c Currency
		if err := rows.Scan(&c.CurrencyCode, &c.DecimalPlaces, &c.ActiveSaturday,
			&c.CurrencyTimeZone, &c.ReconTime); err != nil {
			return nil, err
		}
		curs = append(curs, c)
	}
	return curs, nil
}

// Deposit
// *************
func (d *deposit) postDeposit(db *sql.DB) (string, error) {
	return "112112", nil
}

// Withdraws
// *************
func (w *withdraw) postWithdraw(db *sql.DB) (int, error) {
	return 0, errors.New("Not implemented")
}

// Utils
// *************
func validateDeposit(db *sql.DB, data deposit) string {
	var result string

	fmt.Println("The value of data is ", data)

	if data.Amount <= 0 {
		result = "Invalid Amount"
	}
	if data.Type != "WIRE" && data.Type != "INTEREST" {
		result = result + ", Invalid Transaction Type"
	}

	sqlStmt := `
		SELECT COUNT(results) as data
			FROM (SELECT 1 as results
		        FROM currencies
		        WHERE currency_code = $1
		    	UNION ALL
		    	SELECT 1 as results
		    	FROM customers
		    	WHERE name = $2)`

	rows, err := db.Query(sqlStmt, data.Currency, data.Name)
	if err != nil {
		return "SQL Error"
	}
	defer rows.Close()

	c := 0
	for rows.Next() {
		if rows.Scan(&c); err != nil {
			result = result + ", SQL Error"
			return result
		}
	}

	fmt.Println("The value of count is ", c)

	if c != 2 {
		result = result + ", Invalid Currency or Customer"
	}

	fmt.Println("The value of result is ", result)

	return result
}
