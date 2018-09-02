package infra

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/taiyoh/labeltile/app"
)

var t = template.Must(template.New("ID-dispenser").Parse("INSERT INTO {{.table}}_id_dispenser (id) VALUES (default)"))

func dispenseID(db app.Database, tablePrefix string) string {
	txn := db.NewTransaction()
	defer txn.Rollback()
	b := bytes.NewBuffer([]byte{})
	t.Execute(b, map[string]interface{}{"table": tablePrefix})
	res, err := txn.Mutate(b.String(), []interface{}{}, nil)
	if err != nil {
		return ""
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ""
	}
	txn.Commit()
	return strconv.FormatInt(id, 10)
}
