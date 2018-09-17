package database

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/taiyoh/labeltile/app"
)

// Database provides connection interface to database server
type Database struct {
	app.Database
	db *sql.DB
}

// Transaction provides transaction manager for connection
type Transaction struct {
	app.DatabaseTransaction
	committed bool
	txn       *sql.Tx
}

// New returns new database connection object
func New(driver, dsn string) (app.Database, error) {
	if driver != "sqlite" && driver != "mysql" {
		return nil, errors.New("invalid driver")
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// Close provides closing process for database
func (d *Database) Close() {
	d.db.Close()
}

// Select provides selecting query to database
func (d *Database) Select(query string, args []interface{}, err error) ([]app.DatabaseSelectResult, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	results := []app.DatabaseSelectResult{}
	for rows.Next() {
		scans := []interface{}{}
		for range cols {
			var c string
			scans = append(scans, &c)
		}
		rows.Scan(scans...)
		res := map[string]string{}
		for idx, col := range cols {
			res[col.Name()] = scans[idx].(string)
		}
		results = append(results, res)
	}

	return results, nil
}

// NewTransaction returns new transaction manager object
func (d *Database) NewTransaction() app.DatabaseTransaction {
	txn, _ := d.db.Begin()
	return &Transaction{txn: txn}
}

// Commit provides commit query to database server
func (t *Transaction) Commit() {
	if t.committed {
		return
	}
	t.txn.Commit()
	t.committed = true
}

// Rollback provides rollback query to database server
func (t *Transaction) Rollback() {
	if t.committed {
		return
	}
	t.txn.Rollback()
}

// Context returns context object with transaction object
func (t *Transaction) Context() context.Context {
	return context.WithValue(context.Background(), app.TxnCtxKey, t)
}

// Select provides selecting query in transaction
func (t *Transaction) Select(query string, args []interface{}, err error) (app.DatabaseSelectResult, error) {
	stmt, err := t.txn.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	scans := []interface{}{}
	for range cols {
		var c string
		scans = append(scans, &c)
	}
	if !rows.Next() {
		return nil, nil
	}
	row := app.DatabaseSelectResult{}
	rows.Scan(scans...)
	for idx, c := range cols {
		row[c] = scans[idx].(string)
	}
	return row, nil
}

// Mutate provides create/update/delete query in this transaction
func (t *Transaction) Mutate(query string, args []interface{}, err error) (app.DatabaseMutateResult, error) {
	if err != nil {
		return nil, err
	}
	stmt, err := t.txn.Prepare(query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
