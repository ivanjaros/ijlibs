package query_builder

func NewCondition(fld IField, values []interface{}, operator Operator, negate ...bool) ICondition {
	var not bool
	if len(negate) > 0 && negate[0] {
		not = true
	}

	return &condition{
		fld: fld,
		val: values,
		op:  operator,
		not: not,
	}
}

func NewField(table, column string) IField {
	return &field{table: table, column: column}
}

func Select(table string) IQuery {
	return newQuery(table, SelectQuery)
}

func Insert(table string) IQuery {
	return newQuery(table, InsertQuery)
}

func Update(table string) IQuery {
	return newQuery(table, UpdateQuery)
}

func Delete(table string) IQuery {
	return newQuery(table, DeleteQuery)
}

func NewSort(field IField, order Order) ISort {
	return &Sort{
		fld: field,
		ord: order,
	}
}
