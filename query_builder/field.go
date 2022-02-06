package query_builder

type field struct {
	table  string
	column string
}

func (e field) Table() string {
	return e.table
}

func (e field) Column() string {
	return e.column
}

func (e field) String() string {
	return e.table + TableColumnDivider + e.column
}
