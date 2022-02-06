package query_builder

type executable struct {
	baseTable  string
	qType      Type
	selectCols []IField
	updateCols map[string]interface{}
	insertRows []map[string]interface{}
	innerJoins [][]IField
	leftJoins  [][]IField
	rightJoins [][]IField
	fullJoins  [][]IField
	where      ICondition
	sorts      []ISort
	offset     uint
	limit      uint
	distinct   bool
}

func (e *executable) QueryType() Type {
	return e.qType
}

func (e *executable) BaseTable() string {
	return e.baseTable
}

func (e *executable) InnerJoins() [][]IField {
	return e.innerJoins
}

func (e *executable) LeftJoins() [][]IField {
	return e.leftJoins
}

func (e *executable) RightJoins() [][]IField {
	return e.rightJoins
}

func (e *executable) FullJoins() [][]IField {
	return e.fullJoins
}

func (e *executable) Select() []IField {
	return e.selectCols
}

func (e *executable) Update() map[string]interface{} {
	return e.updateCols
}

func (e *executable) Insert() []map[string]interface{} {
	return e.insertRows
}

func (e *executable) Where() ICondition {
	return e.where
}

func (e *executable) SortBy() []ISort {
	return e.sorts
}

func (e *executable) Range() (offset, limit uint) {
	return e.offset, e.limit
}

func (e *executable) IsDistinct() bool {
	return e.distinct
}
