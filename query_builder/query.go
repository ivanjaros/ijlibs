package query_builder

import (
	"fmt"
)

type query struct {
	exec *executable
}

func newQuery(tbl string, tp Type) IQuery {
	return &query{
		exec: &executable{
			baseTable:  tbl,
			qType:      tp,
			updateCols: map[string]interface{}{},
			insertRows: []map[string]interface{}{},
		},
	}
}

func (q *query) Select(table string, columns ...string) IQuery {
	for _, col := range columns {
		q.exec.selectCols = append(q.exec.selectCols, NewField(table, col))
	}
	return q
}

func (q *query) SelectField(fld IField) IQuery {
	q.exec.selectCols = append(q.exec.selectCols, fld)
	return q
}

func (q *query) Update(column string, value interface{}) IQuery {
	q.exec.updateCols[column] = value
	return q
}

func (q *query) Insert(row map[string]interface{}) IQuery {
	q.exec.insertRows = append(q.exec.insertRows, row)
	return q
}

func (q *query) Join(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error {
	return q.addJoin("inner", srcTable, jTable, srcCol, joinedCol, srcColJoinedColPairs...)
}

func (q *query) LeftJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error {
	return q.addJoin("left", srcTable, jTable, srcCol, joinedCol, srcColJoinedColPairs...)
}

func (q *query) RightJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error {
	return q.addJoin("right", srcTable, jTable, srcCol, joinedCol, srcColJoinedColPairs...)
}

func (q *query) FullJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error {
	return q.addJoin("full", srcTable, jTable, srcCol, joinedCol, srcColJoinedColPairs...)
}

func (q *query) Where() (ICondition, error) {
	if q.exec.where == nil {
		return nil, ErrNoConditions
	}
	return q.exec.where, nil
}

func (q *query) Condition(col IField, operator Operator, value ...interface{}) ICondition {
	if q.exec.where == nil {
		q.exec.where = NewCondition(col, value, operator)
		return q.exec.where
	}
	return q.exec.where.And(col, operator, value)
}

func (q *query) Equals(col IField, value interface{}) ICondition {
	return q.Condition(col, EqualsOperator, value)
}

func (q *query) NotEquals(col IField, value interface{}) ICondition {
	return q.Condition(col, NotEqualsOperator, value)
}

func (q *query) IsNull(col IField) ICondition {
	return q.Condition(col, IsNullOperator)
}

func (q *query) IsNotNull(col IField) ICondition {
	return q.Condition(col, IsNotNullOperator)
}

func (q *query) Larger(col IField, value interface{}) ICondition {
	return q.Condition(col, LargerOperator, value)
}

func (q *query) LargerOrEqual(col IField, value interface{}) ICondition {
	return q.Condition(col, LargerOrEqualOperator, value)
}

func (q *query) Smaller(col IField, value interface{}) ICondition {
	return q.Condition(col, SmallerOperator, value)
}

func (q *query) SmallerOrEqual(col IField, value interface{}) ICondition {
	return q.Condition(col, SmallerOperator, value)
}

func (q *query) Like(col IField, value interface{}) ICondition {
	return q.Condition(col, LikeOperator, value)
}

func (q *query) NotLike(col IField, value interface{}) ICondition {
	return q.Condition(col, NotLikeOperator, value)
}

func (q *query) In(col IField, value interface{}) ICondition {
	return q.Condition(col, InOperator, value)
}

func (q *query) NotIn(col IField, value interface{}) ICondition {
	return q.Condition(col, NotInOperator, value)
}

func (q *query) Between(col IField, value interface{}) ICondition {
	return q.Condition(col, BetweenOperator, value)
}

func (q *query) SortBy(col IField, order Order) IQuery {
	q.exec.sorts = []ISort{NewSort(col, order)}
	return q
}

func (q *query) AddSort(col IField, order Order) IQuery {
	q.exec.sorts = append(q.exec.sorts, NewSort(col, order))
	return q
}

func (q *query) Range(offset, limit uint) IQuery {
	q.exec.offset = offset
	q.exec.limit = limit
	return q
}

func (q *query) Limit(n uint) IQuery {
	q.exec.limit = n
	return q
}

func (q *query) Distinct() IQuery {
	q.exec.distinct = true
	return q
}

func (q *query) Build() IExecutable {
	return q.exec
}

func (q *query) addJoin(kind, srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error {
	on := []IField{
		NewField(srcTable, srcCol),
	}

	if len(srcColJoinedColPairs) > 0 {
		if cols, ok := StringPairs(srcColJoinedColPairs...); ok {
			for k, v := range cols {
				on = append(on, NewField(k, v))
			}
		} else {
			return fmt.Errorf("invalid number of column pairs")
		}
	}

	switch kind {
	case "inner":
		q.exec.innerJoins = append(q.exec.innerJoins, on)
	case "left":
		q.exec.leftJoins = append(q.exec.leftJoins, on)
	case "right":
		q.exec.rightJoins = append(q.exec.rightJoins, on)
	case "full":
		q.exec.fullJoins = append(q.exec.fullJoins, on)
	default:
		return fmt.Errorf("unknown join type '%s'", kind)
	}

	return nil
}
