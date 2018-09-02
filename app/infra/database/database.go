package database

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/taiyoh/labeltile/app"
)

type Database struct {
	app.Database
	db *sql.DB
}

type Transaction struct {
	app.DatabaseTransaction
	committed bool
	txn       *sql.Tx
}

func New(driver, dsn string) (*Database, error) {
	if driver != "sqlite" && driver != "mysql" {
		return nil, errors.New("invalid driver")
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

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

func (d *Database) NewTransaction() *Transaction {
	txn, _ := d.db.Begin()
	return &Transaction{txn: txn}
}

func (t *Transaction) Commit() {
	if t.committed {
		return
	}
	t.txn.Commit()
	t.committed = true
}

func (t *Transaction) Rollback() {
	if t.committed {
		return
	}
	t.txn.Rollback()
}

func (t *Transaction) Context() context.Context {
	return context.WithValue(context.Background(), app.TxnCtxKey, t)
}

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

func (t *Transaction) Mutate(query string, args []interface{}, err error) (app.DatabaseMutateResult, error) {
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
