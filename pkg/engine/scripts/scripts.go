package scripts

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/common"
	"github.com/storm-blue/rubick/pkg/engine/scripts/keywords"
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strconv"
	"strings"
)

// ParseAction
// IF ... THEN ...
func ParseAction(expression string) (action.Action, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, fmt.Errorf("empty action expression")
	}

	conditionPart, pureActionPart, err := splitExpression(expression)
	if err != nil {
		return nil, err
	}
	condition, err := parseCondition(conditionPart)
	if err != nil {
		return nil, err
	}
	pureAction, err := parsePureAction(pureActionPart)
	if err != nil {
		return nil, err
	}

	if condition == nil {
		return pureAction, nil
	} else {
		return action.NewConditionAction(condition, pureAction), nil
	}
}

// parsePureAction like:
// DELETE(...)
// SET(..., "...")
func parsePureAction(expression string) (action.Action, error) {
	method, args, err := splitPureActionExpression(expression)
	if err != nil {
		return nil, err
	}

	switch method {
	case keywords.DELETE:
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 1: %s", keywords.DELETE, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.DELETE, expression)
		}
		return action.NewDeleteAction(key), nil
	case keywords.SET:
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 2: %s", keywords.SET, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.SET, expression)
		}
		value, err := preprocessArgument(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid '%s' expression: second parameter is invalid: %s", keywords.SET, expression)
		}
		return action.NewSetAction(key, value), nil
	case keywords.REPLACE_PART:
		if len(args) != 3 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 3: %s", keywords.REPLACE_PART, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.REPLACE_PART, expression)
		}
		return action.NewReplacePartAction(key, unwrapQuotaIfNeeded(args[1]), unwrapQuotaIfNeeded(args[2])), nil
	case keywords.TRIM_PREFIX:
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 2: %s", keywords.TRIM_PREFIX, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.TRIM_PREFIX, expression)
		}
		prefix, err := preprocessArgument(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid '%s' expression: second parameter is invalid: %s", keywords.TRIM_PREFIX, expression)
		}
		return action.NewTrimPrefixAction(key, prefix), nil
	case keywords.TRIM_SUFFIX:
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 2: %s", keywords.TRIM_SUFFIX, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.TRIM_SUFFIX, expression)
		}
		suffix, err := preprocessArgument(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid '%s' expression: second parameter is invalid: %s", keywords.TRIM_SUFFIX, expression)
		}
		return action.NewTrimSuffixAction(key, suffix), nil
	case keywords.PRINT:
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 1: %s", keywords.PRINT, expression)
		}
		key := unwrapQuotaIfNeeded(args[0])
		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid '%s' expression: key is invalid: %s", keywords.PRINT, expression)
		}
		return action.NewPrintAction(key), nil
	case keywords.REMOVE:
		if len(args) != 0 {
			return nil, fmt.Errorf("invalid '%s' expression: number of parameters must be 0: %s", keywords.REMOVE, expression)
		}
		return action.NewMarkRemovedAction(), nil
	default:
		return nil, fmt.Errorf("invalid action: %s", expression)
	}
}

func preprocessArgument(argument string) (action.Valuable, error) {
	if strings.HasPrefix(argument, keywords.VALUE_OF+"(") {
		if !strings.HasSuffix(argument, ")") {
			return nil, fmt.Errorf("invalid argument: %s", argument)
		}
		key := strings.TrimSpace(argument[len(keywords.VALUE_OF)+1 : len(argument)-1])
		key = unwrapQuotaIfNeeded(key)

		if !objects.IsValidKey(key) {
			return nil, fmt.Errorf("invalid argument: key is invalid: %s", argument)
		}
		return action.ValueOf(key), nil
	} else {
		v := parseToNumberIfPossible(argument)
		return action.Original(v), nil
	}
}

func parseToNumberIfPossible(arg string) interface{} {
	if isWrappedByQuota(arg) {
		return unwrapQuotaIfNeeded(arg)
	}
	if i, err := strconv.Atoi(arg); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(arg, 64); err == nil {
		return f
	}
	return arg
}

// TODO support quotation
func splitExpression(expression string) (conditionPart string, pureActionPart string, err error) {
	if strings.HasPrefix(expression, keywords.IF+" ") {
		index := strings.Index(expression, " "+keywords.THEN+" ")
		if index == -1 {
			return "", "", fmt.Errorf("invalid scripts expression: keywords '%s' not found: %s", keywords.THEN, expression)
		}

		conditionPart = strings.TrimSpace(expression[len(keywords.IF)+1 : index])
		pureActionPart = strings.TrimSpace(expression[index+len(keywords.THEN)+2:])
	} else {
		pureActionPart = strings.TrimSpace(expression)
	}
	return conditionPart, pureActionPart, nil
}

// TODO support quotation
func splitPureActionExpression(expression string) (method string, args []string, err error) {
	i := strings.Index(expression, "(")
	if i == -1 || !strings.HasSuffix(expression, ")") {
		return "", nil, fmt.Errorf("invalid action expression: %s", expression)
	}

	method = expression[:i]
	argsString := expression[i+1 : len(expression)-1]
	argsString = strings.TrimSpace(argsString)

	if argsString == "" {
		return method, nil, nil
	}
	args = strings.Split(argsString, ",")

	for x, arg := range args {
		args[x] = strings.TrimSpace(arg)
	}
	return method, args, nil
}

// parseCondition like:
// (VALUE_OF(...) == "...")) && (VALUE_OF(...) > 0) || EXISTS(...)
func parseCondition(expression string) (conditions.Condition, error) {
	if expression == "" {
		return nil, nil
	}

	expression = common.UnwrapIfNeeded(expression)
	if isSimpleCondition(expression) {
		return parseSimpleCondition(expression)
	}
	return recursiveParseCondition(nil, "", expression)
}

func isSimpleCondition(expression string) bool {
	if strings.HasPrefix(expression, "(") {
		return false
	}
	if strings.Index(expression, keywords.OPERATOR_AND) != -1 || strings.Index(expression, keywords.OPERATOR_OR) != -1 {
		return false
	}
	return true
}

func recursiveParseCondition(headCondition conditions.Condition, operator string, tail string) (conditions.Condition, error) {
	conditionStr, _operator, rest, err := splitCombinationConditionExpression(tail)
	if err != nil {
		return nil, err
	}

	var nextCondition conditions.Condition

	conditionStr = common.UnwrapIfNeeded(conditionStr)
	if isSimpleCondition(conditionStr) {
		nextCondition, err = parseSimpleCondition(conditionStr)
		if err != nil {
			return nil, err
		}
	} else {
		nextCondition, err = recursiveParseCondition(nil, "", conditionStr)
		if err != nil {
			return nil, err
		}
	}

	if headCondition == nil {
		headCondition = nextCondition
	} else {
		switch operator {
		case keywords.OPERATOR_AND:
			headCondition = headCondition.And(nextCondition)
		case keywords.OPERATOR_OR:
			headCondition = headCondition.Or(nextCondition)
		default:
			return nil, fmt.Errorf("parse condition error: invalid logical operator: %s", operator)
		}
	}

	if _operator == "" && rest == "" {
		return headCondition, nil
	} else {
		return recursiveParseCondition(headCondition, _operator, rest)
	}

}

// split express like:
// VALUE_OF(...) == "..."
// LENGTH_OF(...) >= "..."
func splitRelationalSimpleConditionExpression(expression string) (left, right, operator string, err error) {
	index := -1
	for _, _operator := range keywords.RELATIONAL_OPERATORS {
		index = strings.Index(expression, _operator)
		if index != -1 && !indexInQuota(index, expression) {
			operator = _operator
			break
		}
	}
	if index == -1 {
		return "", "", "", fmt.Errorf("invalid relational condition: relational operator not found: %s", expression)
	}
	return strings.TrimSpace(expression[:index]), strings.TrimSpace(expression[index+len(operator):]), operator, nil
}

func parseSimpleCondition(expression string) (conditions.Condition, error) {
	if strings.HasPrefix(expression, keywords.VALUE_OF+"(") {
		return parseValueOfSimpleCondition(expression)
	} else if strings.HasPrefix(expression, keywords.LENGTH_OF+"(") {
		return parseLengthOfSimpleCondition(expression)
	} else if strings.HasPrefix(expression, keywords.EXISTS+"(") {
		return parseExistsSimpleCondition(expression)
	} else if strings.HasPrefix(expression, keywords.NOT_EXISTS+"(") {
		return parseNotExistsSimpleCondition(expression)
	} else {
		return nil, fmt.Errorf("invalid condition expression: must start with '%s' or '%s' or '%s': %s", keywords.VALUE_OF+"(", keywords.EXISTS+"(", keywords.NOT_EXISTS+"(", expression)
	}
}

func parseValueOfSimpleCondition(expression string) (conditions.Condition, error) {
	left, right, operator, err := splitRelationalSimpleConditionExpression(expression)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(left, keywords.VALUE_OF+"(") {
		return nil, fmt.Errorf("invalid '%s' condition: must start with '%s': %s", keywords.VALUE_OF, keywords.VALUE_OF+"(", expression)
	}
	if !strings.HasSuffix(left, ")") {
		return nil, fmt.Errorf("invalid '%s' condition: must end with ')': %s", keywords.VALUE_OF, expression)
	}

	key := strings.TrimSpace(left[len(keywords.VALUE_OF)+1 : len(left)-1])
	if !objects.IsValidKey(key) {
		return nil, fmt.Errorf("invalid '%s' condition: invalid object key: %s", keywords.VALUE_OF, expression)
	}
	value := unwrapQuotaIfNeeded(right)

	switch operator {
	case keywords.OPERATOR_EQ:
		return conditions.New().ValueOf(key).EqualTo(value), nil
	case keywords.OPERATOR_NE:
		return conditions.New().ValueOf(key).NotEqual(value), nil
	case keywords.OPERATOR_GT:
		return conditions.New().ValueOf(key).GreaterThan(value), nil
	case keywords.OPERATOR_GE:
		return conditions.New().ValueOf(key).GreaterThanOrEqual(value), nil
	case keywords.OPERATOR_LT:
		return conditions.New().ValueOf(key).LesserThan(value), nil
	case keywords.OPERATOR_LE:
		return conditions.New().ValueOf(key).LesserThanOrEqual(value), nil
	default:
		return nil, fmt.Errorf("invalid '%s' condition: invalid relational operator: %s", keywords.VALUE_OF, operator)
	}
}

func parseLengthOfSimpleCondition(expression string) (conditions.Condition, error) {
	left, right, operator, err := splitRelationalSimpleConditionExpression(expression)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(left, keywords.LENGTH_OF+"(") {
		return nil, fmt.Errorf("invalid '%s' condition: must start with '%s': %s", keywords.LENGTH_OF, keywords.LENGTH_OF+"(", expression)
	}
	if !strings.HasSuffix(left, ")") {
		return nil, fmt.Errorf("invalid '%s' condition: must end with ')': %s", keywords.LENGTH_OF, expression)
	}

	key := strings.TrimSpace(left[len(keywords.LENGTH_OF)+1 : len(left)-1])
	if !objects.IsValidKey(key) {
		return nil, fmt.Errorf("invalid '%s' condition: invalid object key: %s", keywords.LENGTH_OF, expression)
	}
	right = strings.TrimSpace(right)

	_value, err := strconv.ParseInt(right, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid '%s' condition: invalid value: %s", keywords.LENGTH_OF, expression)
	}

	value := int(_value)

	switch operator {
	case keywords.OPERATOR_EQ:
		return conditions.New().LengthOf(key).EqualTo(value), nil
	case keywords.OPERATOR_NE:
		return conditions.New().LengthOf(key).NotEqual(value), nil
	case keywords.OPERATOR_GT:
		return conditions.New().LengthOf(key).GreaterThan(value), nil
	case keywords.OPERATOR_GE:
		return conditions.New().LengthOf(key).GreaterThanOrEqual(value), nil
	case keywords.OPERATOR_LT:
		return conditions.New().LengthOf(key).LesserThan(value), nil
	case keywords.OPERATOR_LE:
		return conditions.New().LengthOf(key).LesserThanOrEqual(value), nil
	default:
		return nil, fmt.Errorf("invalid '%s' condition: invalid relational operator: %s", keywords.LENGTH_OF, operator)
	}
}

func parseExistsSimpleCondition(expression string) (conditions.Condition, error) {
	if !strings.HasPrefix(expression, keywords.EXISTS+"(") || !strings.HasSuffix(expression, ")") {
		return nil, fmt.Errorf("invalid '%s' condition: must start with '%s': %s", keywords.EXISTS, keywords.EXISTS, expression)
	}
	key := strings.TrimSpace(expression[len(keywords.EXISTS)+1 : len(expression)-1])
	if !objects.IsValidKey(key) {
		return nil, fmt.Errorf("invalid '%s' condition: invalid object key: %s", keywords.EXISTS, expression)
	}
	return conditions.New().Exists(key), nil
}

func parseNotExistsSimpleCondition(expression string) (conditions.Condition, error) {
	if !strings.HasPrefix(expression, keywords.NOT_EXISTS+"(") {
		return nil, fmt.Errorf("invalid '%s' condition: must start with %s: %s", keywords.NOT_EXISTS, keywords.NOT_EXISTS+"(", expression)
	}

	if !strings.HasSuffix(expression, ")") {
		return nil, fmt.Errorf("invalid '%s' condition: must end with ')': %s", keywords.NOT_EXISTS, expression)
	}

	key := strings.TrimSpace(expression[len(keywords.NOT_EXISTS)+1 : len(expression)-1])
	if !objects.IsValidKey(key) {
		return nil, fmt.Errorf("invalid '%s' condition: invalid object key: %s", keywords.NOT_EXISTS, expression)
	}
	return conditions.New().Not(conditions.New().Exists(key)), nil
}

func indexInQuota(index int, expression string) bool {
	inQuota := false
	for i, char := range expression {
		if char == '"' {
			inQuota = !inQuota
		}
		if i == index {
			return inQuota
		}
	}
	return false
}

func isWrappedByQuota(expression string) bool {
	if len(expression) < 2 {
		return false
	}
	if strings.HasPrefix(expression, "\"") && strings.HasSuffix(expression, "\"") {
		return true
	}
	return false
}

func unwrapQuotaIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)
	if len(expression) < 2 {
		return expression
	}
	if strings.HasPrefix(expression, "\"") && strings.HasSuffix(expression, "\"") {
		expression = expression[1 : len(expression)-1]
		return expression
	} else {
		return expression
	}
}

// TODO support quotation
// "expression0 || expression1 && expression2" => "expression0", "||", "expression1 && expression2"
func splitCombinationConditionExpression(expression string) (conditionStr string, operator string, rest string, err error) {
	expression = strings.TrimSpace(expression)
	if strings.HasPrefix(expression, "(") {
		_, index := common.FindFirstParenthesesPair(expression)
		if index == -1 {
			return "", "", "", fmt.Errorf("invalid condition: can not find corresponding ')': %s", expression)
		}

		head := strings.TrimSpace(expression[1:index])
		operatorAndTail := strings.TrimSpace(expression[index+1:])
		operator, rest, err = splitOperatorAndTail(operatorAndTail)
		if err != nil {
			return "", "", "", err
		}
		return head, operator, rest, nil
	} else {
		firstAndIndex := strings.Index(expression, keywords.OPERATOR_AND)
		firstOrIndex := strings.Index(expression, keywords.OPERATOR_OR)

		if firstAndIndex < 0 {
			if firstOrIndex < 0 {
				return expression, "", "", nil
			} else if firstOrIndex == 0 {
				return "", "", "", fmt.Errorf("invalid condition: must not start with '%s': %s", keywords.OPERATOR_OR, expression)
			} else {
				head := strings.TrimSpace(expression[:firstOrIndex])
				operatorAndTail := strings.TrimSpace(expression[firstOrIndex:])
				operator, rest, err = splitOperatorAndTail(operatorAndTail)
				if err != nil {
					return "", "", "", err
				}
				return head, operator, rest, nil
			}
		} else if firstAndIndex == 0 {
			return "", "", "", fmt.Errorf("invalid condition: must not start with '%s': %s", keywords.OPERATOR_AND, expression)
		} else {
			if firstOrIndex < 0 {
				head := strings.TrimSpace(expression[:firstAndIndex])
				operatorAndTail := strings.TrimSpace(expression[firstAndIndex:])
				operator, rest, err = splitOperatorAndTail(operatorAndTail)
				if err != nil {
					return "", "", "", err
				}
				return head, operator, rest, nil
			} else if firstOrIndex == 0 {
				return "", "", "", fmt.Errorf("invalid condition: must not start with '%s': %s", keywords.OPERATOR_OR, expression)
			} else {
				minIndex := Min(firstOrIndex, firstAndIndex)
				head := strings.TrimSpace(expression[:minIndex])
				operatorAndTail := strings.TrimSpace(expression[minIndex:])
				operator, rest, err = splitOperatorAndTail(operatorAndTail)
				if err != nil {
					return "", "", "", err
				}
				return head, operator, rest, nil
			}
		}
	}
}

// TODO support quotation
// "|| expression1 && expression2" => "||", "expression1 && expression2"
func splitOperatorAndTail(operatorAndTail string) (string, string, error) {
	if operatorAndTail == "" {
		return "", "", nil
	}

	if !strings.HasPrefix(operatorAndTail, keywords.OPERATOR_AND) && !strings.HasPrefix(operatorAndTail, keywords.OPERATOR_OR) {
		return "", "", fmt.Errorf("invalid condition: invalid operator and tail: %s", operatorAndTail)
	}

	operatorLen := 0
	if strings.HasPrefix(operatorAndTail, keywords.OPERATOR_AND) {
		operatorLen = len(keywords.OPERATOR_AND)
	} else if strings.HasPrefix(operatorAndTail, keywords.OPERATOR_OR) {
		operatorLen = len(keywords.OPERATOR_OR)
	}

	operator := operatorAndTail[:operatorLen]

	tail := strings.TrimSpace(operatorAndTail[operatorLen:])
	if tail == "" {
		return "", "", fmt.Errorf("invalid condition: invalid operator and tail: %s", operatorAndTail)
	}

	return operator, tail, nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
