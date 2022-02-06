package query_builder

import (
	"errors"
)

type Type string
type Operator string
type Order string

const (
	SelectQuery            = Type("SELECT")
	InsertQuery            = Type("INSERT")
	UpdateQuery            = Type("UPDATE")
	DeleteQuery            = Type("DELETE")
	Ascending              = Order("ASC")
	Descending             = Order("DESC")
	TableColumnDivider     = `.`
	MatchAll               = `*`
	LikeAny                = `%`
	LikeSingle             = `_`
	Escape                 = `\`
	EqualsOperator         = Operator("=")
	NotEqualsOperator      = Operator("!=")
	IsNullOperator         = Operator("IS NULL")
	IsNotNullOperator      = Operator("IS NOT NULL")
	LargerOperator         = Operator(">")
	LargerOrEqualOperator  = Operator(">=")
	SmallerOperator        = Operator("<")
	SmallerOrEqualOperator = Operator("<=")
	LikeOperator           = Operator("LIKE")
	NotLikeOperator        = Operator("NOT LIKE")
	InOperator             = Operator("IN")
	NotInOperator          = Operator("NOT IN")
	BetweenOperator        = Operator("BETWEEN")
)

var (
	ErrNoConditions = errors.New("no conditions have been added yet")
)

// This interface represents a column on a table or an expression.
type IField interface {
	// Table can be empty, in which case the column should be
	// considered being on the base table or being an expression.
	Table() string
	Column() string
	String() string
}

// This interface is only for building queries. It is up to the database engine
// or the executable parser to validate if it can handle the query parameters.
// For example only columns might be supported, no expressions.
// Only some types of joins are not supported and so on.
type IQuery interface {
	Select(table string, columns ...string) IQuery
	// This supports expressions
	SelectField(fld IField) IQuery
	// Column has to be on the base table
	Update(column string, value interface{}) IQuery
	// String keys are table columns
	Insert(row map[string]interface{}) IQuery

	// Returns records that have matching values in both tables, errors if column pairs are not even
	Join(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error
	// Return all records from the left table, and the matched records from the right table, errors if column pairs are not even
	LeftJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error
	// Return all records from the right table, and the matched records from the left table, errors if column pairs are not even
	RightJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error
	//  Return all records when there is a match in either left or right table, errors if column pairs are not even
	FullJoin(srcTable, jTable, srcCol, joinedCol string, srcColJoinedColPairs ...string) error

	// Returns reference to the root condition, errors into ErrNoConditions if there is none set yet
	Where() (ICondition, error)
	// Returns reference to newly created condition, if it is the first condition, it can be chained since it is the root one.
	Condition(fld IField, operator Operator, value ...interface{}) ICondition
	Equals(fld IField, value interface{}) ICondition         // alias for Condition
	NotEquals(fld IField, value interface{}) ICondition      // alias for Condition
	IsNull(fld IField) ICondition                            // alias for Condition
	IsNotNull(fld IField) ICondition                         // alias for Condition
	Larger(fld IField, value interface{}) ICondition         // alias for Condition
	LargerOrEqual(fld IField, value interface{}) ICondition  // alias for Condition
	Smaller(fld IField, value interface{}) ICondition        // alias for Condition
	SmallerOrEqual(fld IField, value interface{}) ICondition // alias for Condition
	Like(fld IField, value interface{}) ICondition           // alias for Condition
	NotLike(fld IField, value interface{}) ICondition        // alias for Condition
	In(fld IField, value interface{}) ICondition             // alias for Condition
	NotIn(fld IField, value interface{}) ICondition          // alias for Condition
	Between(fld IField, value interface{}) ICondition        // alias for Condition

	SortBy(fld IField, order Order) IQuery  // Sets sorting to only the provided one
	AddSort(fld IField, order Order) IQuery // Adds sort
	Range(offset, limit uint) IQuery
	Limit(n uint) IQuery // Sets only limit if offset can be 0, better dx
	Distinct() IQuery    // Sets query to be distinct
	Build() IExecutable  // Processes all columns, aliases and other data to create usable data set
}

type IExecutable interface {
	QueryType() Type
	BaseTable() string
	InnerJoins() [][]IField
	LeftJoins() [][]IField
	RightJoins() [][]IField
	FullJoins() [][]IField
	Select() []IField
	Update() map[string]interface{}
	Insert() []map[string]interface{}
	Where() ICondition
	SortBy() []ISort
	Range() (offset, limit uint)
	IsDistinct() bool
}

type ICondition interface {
	Table() string
	Column() string
	Value() []interface{}
	Operator() Operator
	Siblings() []ICondition
	Alternatives() []ICondition
	And(expr IField, operator Operator, value ...interface{}) ICondition
	Or(fld IField, operator Operator, value ...interface{}) ICondition
	Nand(fld IField, operator Operator, value ...interface{}) ICondition
	Nor(fld IField, operator Operator, value ...interface{}) ICondition
	IsNegated() bool
}

type ISort interface {
	Field() IField
	Order() Order
}
