package main

import (
	"database/sql"
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
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Rate          float32 `json:"rate"`
	RefID         int     `json:"refid"`
}

type withdraw struct {
	Lei          string  `json:"lei"`
	AcctNumber   string  `json:"account_number"`
	CurrencyCode string  `json:"currency_code"`
	Amount       float32 `json:"amount"`
	BankName     string  `json:"bank_name"`
	Instructions string  `json:"instructions"`
	RefID        int     `json:"refid"`
	Notes        string  `json:"notes"`
}

type message struct {
	Message string `json:"message"`
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
	t.QAccount = d.QuorumAccount
	if d.Type == "WIRE" {
		t.Description = "Customer Deposit"
	} else {
		t.Description = "Interest Rate Deposit"
	}
	t.Amount = d.Amount
	t.StartDate = d.StartDate
	t.EndDate = d.EndDate
	t.Rate = d.Rate
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
	n.QAccount = d.QuorumAccount
	n.StartDate = d.StartDate
	n.EndDate = d.EndDate
	n.Rate = d.Rate

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

	cId, err := lookUpCustomerLEI(db, w.Lei)
	if err != nil {
		return err
	}

	// Post Transaction
	var t transaction

	t.TransType = "WITHDRAW"
	t.Currency = w.CurrencyCode
	t.AcctNumber = w.AcctNumber
	t.CustomerID = cId
	t.Description = w.Instructions
	t.Amount = w.Amount * (-1)
	t.Status = "POSTED"

	w.RefID, err = t.createTransaction(db)
	if err != nil {
		return err
	}

	// Create Notice
	var n notification

	n.NoticeType = t.TransType
	n.CustomerID = cId
	n.TransID = w.RefID
	n.Message = t.Description
	n.Amount = t.Amount
	n.Currency = t.Currency

	err = n.createNotification(db)
	if err != nil {
		return err
	}

	// Update Balance
	err = updateAcctBalance(db, w.AcctNumber, t.Amount)
	if err != nil {
		return err
	}

	return nil
}

// Utils
// *************
func validateDeposit(db *sql.DB, data deposit) string {
	result := ""

	if data.Amount <= 0 {
		result = "Invalid Amount"
	}
	if data.Type != "WIRE" && data.Type != "INTEREST" {
		result = result + "Invalid Transaction Type, "
	}

	// verify currency
	act, err := lookUpAccount(db, data.Currency)
	if err != nil {
		result = result + "SQL Error, "
	}
	if act == "" {
		result = result + "Invalid Currency, "
	}
	// Verify customer
	if data.Type == "WIRE" {
		cID, err := lookUpCustomerID(db, data.Name)
		if err != nil {
			result = result + "SQL Error, "
		}
		if cID == 0 {
			result = result + "Invalid Customer, "
		}
		if data.QuorumAccount == "" {
			result = result + "Quorum Account reqired, "
		}
		if data.StartDate != "" && data.EndDate != "" {
			result = result + "Dates not Needed for WIRE, "
		}
		if data.Rate != 0.0 {
			result = result + "Rate not needed for WIRE, "
		}
	}
	if data.Type == "INTEREST" {
		if data.Name != "" {
			result = result + "Name not needed, "
		}
		if data.QuorumAccount != "" {
			result = result + "Quorum Account not needed, "
		}
		if data.Rate <= 0.0 {
			result = result + "Rate must be > 0, "
		}
		if data.StartDate == "" || data.EndDate == "" {
			result = result + "Start and End Date required.  "
		}
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

func lookUpCustomerLEI(db *sql.DB, lei string) (int, error) {
	// Get account Number
	lookUpCustomerID := `
		SELECT id 
		FROM customers 
		WHERE lei = $1`

	rows, err := db.Query(lookUpCustomerID, lei)
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

func validateWithdraw(db *sql.DB, data withdraw) string {
	result := ""

	if data.Amount <= 0 {
		result = "Invalid Amount, "
	}
	// Validate lei
	custID, err := lookUpCustomerLEI(db, data.Lei)
	if err != nil {
		result = result + "SQL Error, "
	}
	if custID == 0 {
		result = result + "Invalid Customer, "
	}
	// Validate account number
	act, err := lookUpAccount(db, data.CurrencyCode)
	if err != nil {
		result = result + "SQL Error, "
	}
	if act == "" {
		result = result + "Invalid Currency for this bank, "
	}
	if act != data.AcctNumber {
		result = result + "Invalid Account Number"
	}

	return result
}
func getBalanceByCurrency(db *sql.DB, currency string) (float64, error) {
	// Get account Number
	lookUpBalance := `
		SELECT balance
		FROM accounts
		WHERE currency_code = $1`

	rows, err := db.Query(lookUpBalance, currency)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	bal := 0.0
	for rows.Next() {
		if rows.Scan(&bal); err != nil {
			return 0, err
		}
	}

	return bal, nil
}
