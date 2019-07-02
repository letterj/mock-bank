package main

import (
	"database/sql"
)

type account struct {
	AcctNumber    string  `json:"acct_number"`
	QuorumAccount string  `json:"quorum_account"`
	CurrencyCode  string  `json:"currency_code"`
	Balance       float32 `json:"balance"`
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
		if err := rows.Scan(&act.AcctNumber, &act.QuorumAccount,
			&act.CurrencyCode, &act.Balance); err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	return acts, nil
}

func (acc *account) getAccount(db *sql.DB) error {
	sqlStmt := `SELECT acct_number, qaccount, currency_code, balance
		FROM accounts WHERE acct_number=$1`
	return db.QueryRow(sqlStmt,
		acc.AcctNumber).Scan(&acc.QuorumAccount, &acc.CurrencyCode, &acc.Balance)
}
