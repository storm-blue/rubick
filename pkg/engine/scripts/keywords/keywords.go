package keywords

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/common"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strconv"
	"strings"
)

//goland:noinspection ALL
const (
	IF           = "IF"
	THEN         = "THEN"
	VALUE_OF     = "VALUE_OF"
	LENGTH_OF    = "LENGTH_OF"
	EXISTS       = "EXISTS"
	NOT_EXISTS   = "NOT_EXISTS"
	HAS_PREFIX   = "HAS_PREFIX"
	HAS_SUFFIX   = "HAS_SUFFIX"
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

	SINGLE_WORDS_SIMPLE_CONDITION_METHODS = map[string]func(args ...string) (conditions.Condition, error){
		EXISTS: func(args ...string) (conditions.Condition, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 1", EXISTS)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", EXISTS, key)
			}
			return conditions.New().Exists(key), nil
		},
		NOT_EXISTS: func(args ...string) (conditions.Condition, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 1", NOT_EXISTS)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", NOT_EXISTS, key)
			}
			return conditions.New().Not(conditions.New().Exists(key)), nil
		},
		HAS_PREFIX: func(args ...string) (conditions.Condition, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 2", HAS_PREFIX)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", HAS_PREFIX, key)
			}
			prefix := common.UnwrapQuotaIfNeeded(args[1])
			return conditions.New().HasPrefix(key, prefix), nil
		},
		HAS_SUFFIX: func(args ...string) (conditions.Condition, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 2", HAS_SUFFIX)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", HAS_SUFFIX, key)
			}
			suffix := common.UnwrapQuotaIfNeeded(args[1])
			return conditions.New().HasSuffix(key, suffix), nil
		},
	}

	RELATIONAL_SIMPLE_CONDITION_METHODS = map[string]func(operator string, rightValue string, args ...string) (conditions.Condition, error){
		VALUE_OF: func(operator string, rightValue string, args ...string) (conditions.Condition, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 1", VALUE_OF)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", VALUE_OF, key)
			}

			value := common.UnwrapQuotaIfNeeded(rightValue)
			switch operator {
			case OPERATOR_EQ:
				return conditions.New().ValueOf(key).EqualTo(value), nil
			case OPERATOR_NE:
				return conditions.New().ValueOf(key).NotEqual(value), nil
			case OPERATOR_GT:
				return conditions.New().ValueOf(key).GreaterThan(value), nil
			case OPERATOR_GE:
				return conditions.New().ValueOf(key).GreaterThanOrEqual(value), nil
			case OPERATOR_LT:
				return conditions.New().ValueOf(key).LesserThan(value), nil
			case OPERATOR_LE:
				return conditions.New().ValueOf(key).LesserThanOrEqual(value), nil
			default:
				return nil, fmt.Errorf("invalid '%s' condition: invalid relational operator: %s", VALUE_OF, operator)
			}
		},
		LENGTH_OF: func(operator string, rightValue string, args ...string) (conditions.Condition, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid '%s' condition: number of parameters must be 1", LENGTH_OF)
			}
			key := common.UnwrapQuotaIfNeeded(args[0])
			if !objects.IsValidKey(key) {
				return nil, fmt.Errorf("invalid '%s' condition: key is invalid: %s", LENGTH_OF, key)
			}

			rightValue = strings.TrimSpace(rightValue)
			_value, err := strconv.ParseInt(rightValue, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid '%s' condition: invalid value: %s", LENGTH_OF, rightValue)
			}

			value := int(_value)
			switch operator {
			case OPERATOR_EQ:
				return conditions.New().LengthOf(key).EqualTo(value), nil
			case OPERATOR_NE:
				return conditions.New().LengthOf(key).NotEqual(value), nil
			case OPERATOR_GT:
				return conditions.New().LengthOf(key).GreaterThan(value), nil
			case OPERATOR_GE:
				return conditions.New().LengthOf(key).GreaterThanOrEqual(value), nil
			case OPERATOR_LT:
				return conditions.New().LengthOf(key).LesserThan(value), nil
			case OPERATOR_LE:
				return conditions.New().LengthOf(key).LesserThanOrEqual(value), nil
			default:
				return nil, fmt.Errorf("invalid '%s' condition: invalid relational operator: %s", LENGTH_OF, operator)
			}
		},
	}
)
