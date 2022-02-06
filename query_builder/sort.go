package query_builder

type Sort struct {
	fld IField
	ord Order
}

func (s *Sort) Field() IField {
	return s.fld
}

func (s *Sort) Order() Order {
	return s.ord
}
