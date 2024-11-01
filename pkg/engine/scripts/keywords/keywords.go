package keywords

//goland:noinspection ALL
const (
	IF                = "IF"
	THEN              = "THEN"
	VALUE_OF          = "VALUE_OF"
	EXISTS            = "EXISTS"
	NOT_EXISTS        = "NOT_EXISTS"
	DELETE            = "DELETE"
	SET               = "SET"
	REPLACE_PART      = "REPLACE_PART"
	SET_WITH_VALUE_OF = "SET_WITH_VALUE_OF"

	// operators
	OPERATOR_EQ = "=="
	OPERATOR_NE = "!="
	OPERATOR_LE = "<="
	OPERATOR_GE = ">="
	OPERATOR_LT = "<"
	OPERATOR_GT = ">"

	OPERATOR_AND = "&&"
	OPERATOR_OR  = "||"
	OPERATOR_NOT = "!"
)

//goland:noinspection ALL
var (
	RELATIONAL_OPERATORS = []string{OPERATOR_EQ, OPERATOR_NE, OPERATOR_LE, OPERATOR_GE, OPERATOR_LT, OPERATOR_GT}
	LOGICAL_OPERATORS    = []string{OPERATOR_AND, OPERATOR_OR, OPERATOR_NOT}
)
