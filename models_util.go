package main

import (
	"database/sql"
	"errors"
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
	RefID         int     `json:"refid"`
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
func (d *deposit) postDeposit(db *sql.DB) error {

	acct, err := lookUpAccount(db, d.Currency)
	if err != nil {
		return err
	}

	cId, err := lookUpCustomerID(db, d.Name)
	if err != nil {
		return err
	}

	// Post Transaction
	var t transaction

	t.TransType = d.Type
	t.Currency = d.Currency
	t.AcctNumber = acct
	t.CustomerID = cId
	if d.Type == "WIRE" {
		t.Description = "Depsit QAccount: " + d.QuorumAccount
	} else {
		t.Description = "SWEEP INTEREST PAYMENT"
	}
	t.Amount = d.Amount
	t.Status = "POSTED"

	d.RefID, err = t.createTransaction(db)
	if err != nil {
		return err
	}

	// Create Notice
	var n notification

	n.NoticeType = d.Type
	n.CustomerID = cId
	n.TransID = d.RefID
	n.Message = t.Description
	n.Amount = d.Amount
	n.Currency = d.Currency

	err = n.createNotification(db)
	if err != nil {
		return err
	}

	// Update Balance
	err = updateAcctBalance(db, t.AcctNumber, t.Amount)
	if err != nil {
		return err
	}

	return nil

}

// Withdraws
// *************
func (w *withdraw) postWithdraw(db *sql.DB) error {
	return errors.New("Not implemented")
}

// Utils
// *************
func validateDeposit(db *sql.DB, data deposit) string {
	var result string

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

	if c != 2 {
		result = result + ", Invalid Currency or Customer"
	}

	return result
}

func lookUpAccount(db *sql.DB, currency string) (string, error) {
	// Get account Number
	lookUpAccount := `
		SELECT acct_number 
		FROM accounts 
		WHERE currency_code = $1`

	rows, err := db.Query(lookUpAccount, currency)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	acct := ""
	for rows.Next() {
		if rows.Scan(&acct); err != nil {
			return "", err
		}
	}

	return acct, nil
}

func lookUpCustomerID(db *sql.DB, CustName string) (int, error) {
	// Get account Number
	lookUpCustomerID := `
		SELECT id 
		FROM customers 
		WHERE name = $1`

	rows, err := db.Query(lookUpCustomerID, CustName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	custID := 0
	for rows.Next() {
		if rows.Scan(&custID); err != nil {
			return 0, err
		}
	}

	return custID, nil
}

func updateAcctBalance(db *sql.DB, acctNum string, amt float32) error {
	// Get account Number
	updateBal := `
		UPDATE accounts SET balance = balance + ?
		WHERE acct_number = ?`

	sqlStmt, err := db.Prepare(updateBal)
	if err != nil {
		return err
	}

	sqlStmt.Exec(amt, acctNum)

	return nil
}
