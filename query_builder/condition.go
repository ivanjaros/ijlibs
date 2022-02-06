package query_builder

type condition struct {
	fld IField
	val []interface{}
	op  Operator
	sbl []ICondition
	alt []ICondition
	not bool
}

func (c *condition) Table() string {
	return c.fld.Table()
}

func (c *condition) Column() string {
	return c.fld.Column()
}

func (c *condition) Value() []interface{} {
	return c.val
}

func (c *condition) Operator() Operator {
	return c.op
}

func (c *condition) Siblings() []ICondition {
	return c.sbl
}

func (c *condition) Alternatives() []ICondition {
	return c.alt
}

func (c *condition) And(fld IField, operator Operator, value ...interface{}) ICondition {
	cond := NewCondition(fld, value, operator)
	c.sbl = append(c.sbl, cond)
	return cond
}

func (c *condition) Or(fld IField, operator Operator, value ...interface{}) ICondition {
	cond := NewCondition(fld, value, operator)
	c.alt = append(c.alt, cond)
	return cond
}

func (c *condition) Nand(fld IField, operator Operator, value ...interface{}) ICondition {
	cond := NewCondition(fld, value, operator, true)
	c.sbl = append(c.sbl, cond)
	return cond
}

func (c *condition) Nor(fld IField, operator Operator, value ...interface{}) ICondition {
	cond := NewCondition(fld, value, operator, true)
	c.alt = append(c.alt, cond)
	return cond
}

func (c *condition) IsNegated() bool {
	return c.not
}
