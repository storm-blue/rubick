package keywords

//goland:noinspection ALL
const (
	IF           = "IF"
	THEN         = "THEN"
	VALUE_OF     = "VALUE_OF"
	LENGTH_OF    = "LENGTH_OF"
	EXISTS       = "EXISTS"
	NOT_EXISTS   = "NOT_EXISTS"
	DELETE       = "DELETE"
	SET          = "SET"
	REPLACE_PART = "REPLACE_PART"
	TRIM_PREFIX  = "TRIM_PREFIX"
	TRIM_SUFFIX  = "TRIM_SUFFIX"
	PRINT        = "PRINT"
	REMOVE       = "REMOVE"

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
