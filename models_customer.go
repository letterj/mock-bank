package main

import (
	"database/sql"
	"errors"
)

type customer struct {
	ID            int    `json:"id"`
	Lei           string `json:"lei"`
	Name          string `json:"name"`
	QuorumAccount string `json:"quorum_account"`
}

// Customers
// **********
func getCustomers(db *sql.DB) ([]customer, error) {
	sqlstmt := `
	SELECT id, lei, name, qaccount
	FROM customers`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	custs := []customer{}
	for rows.Next() {
		var cust customer
		if err := rows.Scan(&cust.ID, &cust.Lei, &cust.Name, &cust.QuorumAccount); err != nil {
			return nil, err
		}
		custs = append(custs, cust)
	}
	return custs, nil
}

func (c *customer) getCustomer(db *sql.DB) error {
	sqlStmt := `SELECT id, lei, name, qaccount
	FROM customers WHERE lei=$1`
	return db.QueryRow(sqlStmt,
		c.Lei).Scan(&c.ID, &c.Name, &c.QuorumAccount)
}

func (c *customer) createCustomer(db *sql.DB) error {
	insertCustomer := `INSERT INTO customers (lei, name, qaccount) 
		VALUES($1, $2, $3);`

	sqlStmt, err := db.Prepare(insertCustomer)
	if err != nil {
		return err
	}
	_, err = sqlStmt.Exec(c.Lei, c.Name, c.QuorumAccount)
	if err != nil {
		return err
	}
	return nil
}

func (c *customer) updateCustomer(db *sql.DB) error {
	return errors.New("Not implemented")
}
