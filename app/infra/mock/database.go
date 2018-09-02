package mock

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/taiyoh/labeltile/app"
)

type SQLCapture struct {
	Query string
	Args  []interface{}
}

type SQLCaptureStore struct {
	captured []*SQLCapture
}

func (c *SQLCaptureStore) Append(q string, a []interface{}) {
	s := &SQLCapture{Query: q, Args: a}
	c.captured = append(c.captured, s)
}

func (c *SQLCaptureStore) ReadAll() []*SQLCapture {
	return c.captured
}

func (c *SQLCaptureStore) Clear() {
	c.captured = []*SQLCapture{}
}

type DB struct {
	app.Database
	captured      *SQLCaptureStore
	mutateResult  *MutateResult
	mutateErr     error
	selectResults map[string][]app.DatabaseSelectResult
	selectErr     error
}

type Txn struct {
	app.DatabaseTransaction
	committed bool
	captured  *SQLCaptureStore
	mResult   *MutateResult
	mErr      error
}

type MutateResult struct {
	app.DatabaseMutateResult
	i  int64
	ie error
	a  int64
	ae error
}

func (r *MutateResult) LastInsertId() (int64, error) {
	return r.i, r.ie
}

func (r *MutateResult) RowsAffected() (int64, error) {
	return r.a, r.ae
}

func (t *Txn) Select(q string, a []interface{}, err error) (app.DatabaseSelectResult, error) {
	t.captured.Append(q, a)
	res := app.DatabaseSelectResult{}
	return res, nil
}

func (t *Txn) Mutate(q string, a []interface{}, err error) (app.DatabaseMutateResult, error) {
	t.captured.Append(q, a)
	return t.mResult, t.mErr
}

func (t *Txn) Rollback() {
	if !t.committed {
		t.captured.Append("ROLLBACK", []interface{}{})
	}
}

func (t *Txn) Commit() {
	if !t.committed {
		t.captured.Append("COMMIT", []interface{}{})
		t.committed = true
	}
}

func LoadDatabase() *DB {
	return &DB{
		captured:      &SQLCaptureStore{captured: []*SQLCapture{}},
		selectResults: map[string][]app.DatabaseSelectResult{},
		selectErr:     nil,
	}
}

func (d *DB) Captures() []*SQLCapture {
	return d.captured.ReadAll()
}

func (d *DB) ClearCapture() {
	d.captured.Clear()
}

func (d *DB) Close() {}

func (d *DB) MutateResult(lastInsertId int64, lastInsertIdErr error, rowsAffected int64, rowsAffectedErr error) {
	d.mutateResult = &MutateResult{
		i:  lastInsertId,
		ie: lastInsertIdErr,
		a:  rowsAffected,
		ae: rowsAffectedErr,
	}
	d.mutateErr = nil
}

func (d *DB) MutateError(err error) {
	d.mutateResult = nil
	d.mutateErr = err
}

func (d *DB) NewTransaction() app.DatabaseTransaction {
	d.captured.Append("BEGIN", []interface{}{})
	txn := &Txn{captured: d.captured, mResult: d.mutateResult, mErr: d.mutateErr}
	return txn
}

var fromT = template.Must(template.New("from").Parse("FROM {{.table}} "))

func (d *DB) Select(q string, a []interface{}, err error) ([]app.DatabaseSelectResult, error) {
	d.captured.Append(q, a)
	for k, pool := range d.selectResults {
		b := bytes.NewBuffer([]byte{})
		fromT.Execute(b, map[string]string{"table": k})
		if strings.Contains(q, b.String()) {
			return pool, nil
		}
	}
	return []app.DatabaseSelectResult{}, nil
}

func (d *DB) ClearData() {
	d.selectResults = map[string][]app.DatabaseSelectResult{}
	d.selectErr = nil
}

func (d *DB) Add(name string, dataList ...app.DatabaseSelectResult) {
	p, exists := d.selectResults[name]
	if !exists {
		p = []app.DatabaseSelectResult{}
	}
	for _, data := range dataList {
		p = append(p, data)
	}
	d.selectResults[name] = p
}
