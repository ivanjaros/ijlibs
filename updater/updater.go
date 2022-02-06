package updater

import (
	"database/sql"
	"fmt"
	"github.com/ivanjaros/ijlibs/placeholders"
)

func New(table, primaryKeyColumn string) *updater {
	return &updater{table: table, pkCol: primaryKeyColumn, values: make(map[string][][2]interface{})}
}

type updater struct {
	table  string
	pkCol  string
	values map[string][][2]interface{}
}

func (u *updater) Add(pkVal interface{}, column string, value interface{}) {
	changes := u.values[column]
	if changes == nil {
		changes = make([][2]interface{}, 0, 50)
	}

	for k := range changes {
		// fmt.Sprint fixes issue when comparing byte slices
		if fmt.Sprint(changes[k][0]) == fmt.Sprint(pkVal) {
			changes[k][1] = value
			u.values[column] = changes
			return
		}
	}

	u.values[column] = append(changes, [2]interface{}{pkVal, value})
}

// Build list of queries and their arguments to run.
// Uses the sql CASE approach described in here:
// https://blog.bubble.ro/how-to-make-multiple-updates-using-a-single-query-in-mysql/
//
// UPDATE mytable SET title = CASE
// WHEN id = 1 THEN ‘Great Expectations’
// WHEN id = 2 THEN ‘War and Peace’
// ...
// ELSE col
// END
// WHERE id IN (1,2,...)
func (u *updater) Build() map[string][]interface{} {
	queries := make(map[string][]interface{})

	for col, changes := range u.values {
		ln := len(changes)
		query := "UPDATE " + u.table + " SET " + col + " = CASE "
		args := make([]interface{}, 0, ln*2+ln)
		pks := make([]interface{}, 0, len(changes))
		for _, change := range changes {
			query += "WHEN " + u.pkCol + " = ? THEN ? "
			args = append(args, change[0], change[1])
			pks = append(pks, change[0])
		}
		query += "ELSE " + col
		query += " END WHERE " + u.pkCol + " IN " + placeholders.Group(ln)
		args = append(args, pks...)
		queries[query] = args
	}

	return queries
}

// simple convenience function
func (u *updater) Execute(db interface {
	Exec(string, ...interface{}) (sql.Result, error)
}) error {
	tasks := u.Build()
	for k, v := range tasks {
		if _, err := db.Exec(k, v...); err != nil {
			return err
		}
	}
	return nil
}
