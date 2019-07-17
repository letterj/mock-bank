package main

import (
	"database/sql"
)

type account struct {
	AcctNumber   string  `json:"acct_number"`
	QAccount     string  `json:"quorum_account"`
	CurrencyCode string  `json:"currency_code"`
	Balance      float32 `json:"balance"`
}

// Accounts
// ********
func getAccounts(db *sql.DB) ([]account, error) {
	sqlstmt := `
	SELECT acct_number, qaccount, currency_code, balance
	FROM accounts`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	acts := []account{}
	for rows.Next() {
		var act account
		if err := rows.Scan(&act.AcctNumber, &act.QAccount,
			&act.CurrencyCode, &act.Balance); err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	return acts, nil
}

func (acc *account) getAccount(db *sql.DB) error {
	sqlStmt := `SELECT acct_number, qaccount, currency_code, balance
		FROM accounts WHERE currency_code=$1 limit 1`
	return db.QueryRow(sqlStmt,
		acc.CurrencyCode).Scan(&acc.AcctNumber, &acc.QAccount, &acc.CurrencyCode,
		&acc.Balance)
}

func (a *account) updateAccount(db *sql.DB) error {
	updateAccount := `
	UPDATE accounts
	SET qaccount = $2
	WHERE acct_number = $1`

	sqlStmt, err := db.Prepare(updateAccount)
	if err != nil {
		return err
	}

	sqlStmt.Exec(&a.AcctNumber, &a.QAccount)

	refresh := `
	SELECT acct_number, currency_code, balance
	FROM accounts WHERE acct_number = $1`

	err = db.QueryRow(refresh,
		a.AcctNumber).Scan(&a.AcctNumber, &a.CurrencyCode,
		&a.Balance)
	if err != nil {
		return err
	}

	return nil
}
