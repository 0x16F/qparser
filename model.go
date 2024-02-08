package qparser

const (
	operatorEqual            = "eq"
	operatorNotEqual         = "neq"
	operatorGreaterThan      = "gt"
	operatorGreaterThanEqual = "gte"
	operatorLowerThan        = "lt"
	operatorLowerThanEqual   = "lte"
	operatorLike             = "like"
	operatorRange            = "rng"
)

const (
	sqlOperatorEqual            = "="
	sqlOperatorNotEqual         = "<>"
	sqlOperatorGreaterThan      = ">"
	sqlOperatorGreaterThanEqual = ">="
	sqlOperatorLowerThan        = "<"
	sqlOperatorLowerThanEqual   = "<="
	sqlOperatorLike             = "ILIKE"
	sqlOperatorRange            = "BETWEEN"
)

type Field struct {
	Name     string
	Value    string
	Operator string
}

type Options struct {
	limit  int
	offset int
	fields []*Field
}
