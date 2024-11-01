package engine

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"strconv"
	"strings"
)

// ParseAction
// IF VALUE_OF(a.b.c[0]) = 1 THEN DELETE(x.y.z)
// IF VALUE_OF(a.b.c[0]) = 1 THEN SET(x.y.z, "shit")
// IF VALUE_OF(a.b.c[0]) = 1 THEN REPLACE_PART(x.y.z, "shit","shit0")
// IF VALUE_OF(a.b.c[0]) = 1 THEN SET(x.y.z, VALUE_OF(a.b.c))
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

// parsePureAction
// DELETE(x.y.z)
// SET(x.y.z, "shit")
// REPLACE_PART(x.y.z, "shit","shit0")
// SET(x.y.z, VALUE_OF(a.b.c))
// SET_WITH_VALUE_OF(x.y.z, zy.c)
func parsePureAction(expression string) (action.Action, error) {
	method, args, err := splitPureActionExpression(expression)
	if err != nil {
		return nil, err
	}

	switch method {
	case "DELETE":
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid delete pure action: %s", expression)
		}
		return action.NewDeleteAction(args[0]), nil
	case "SET":
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid set pure action: %s", expression)
		}
		return action.NewSetAction(unwrapQuotaIfNeeded(args[0]), parseToNumberIfPossible(args[1])), nil
	case "REPLACE_PART":
		if len(args) != 3 {
			return nil, fmt.Errorf("invalid replace pure action: %s", expression)
		}
		return action.NewReplacePartAction(unwrapQuotaIfNeeded(args[0]), unwrapQuotaIfNeeded(args[1]), unwrapQuotaIfNeeded(args[2])), nil
	case "SET_WITH_VALUE_OF":
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid set with value of pure action: %s", expression)
		}
		return action.NewSetWithValueOfAction(unwrapQuotaIfNeeded(args[0]), unwrapQuotaIfNeeded(args[1])), nil
	default:
		return nil, fmt.Errorf("invalid replace pure action: %s", expression)
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
	if strings.HasPrefix(expression, "IF ") {
		index := strings.Index(expression, " THEN ")
		if index == -1 {
			return "", "", fmt.Errorf("invalid action expression: %s", expression)
		}

		conditionPart = strings.TrimSpace(expression[3:index])
		pureActionPart = strings.TrimSpace(expression[index+6:])
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
	if argsString == "" {
		return method, nil, nil
	}
	args = strings.Split(argsString, ",")

	for x, arg := range args {
		args[x] = strings.TrimSpace(arg)
	}
	return method, args, nil
}

// parseCondition
// (VALUE(a.b[2].c=="true"))&&(VALUE(x.y.z)>123)
func parseCondition(expression string) (conditions.Condition, error) {
	if expression == "" {
		return nil, nil
	}

	expression = unwrapIfNeeded(expression)
	if isSimpleCondition(expression) {
		return parseSimpleCondition(expression)
	}
	return recursiveParseCondition(nil, "", expression)
}

func isSimpleCondition(expression string) bool {
	if strings.HasPrefix(expression, "(") {
		return false
	}
	if strings.Index(expression, "&&") != -1 || strings.Index(expression, "||") != -1 {
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

	conditionStr = unwrapIfNeeded(conditionStr)
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
		case "&&":
			headCondition = headCondition.And(nextCondition)
		case "||":
			headCondition = headCondition.Or(nextCondition)
		default:
			return nil, fmt.Errorf("invalid operator: %s", operator)
		}
	}

	if _operator == "" && rest == "" {
		return headCondition, nil
	} else {
		return recursiveParseCondition(headCondition, _operator, rest)
	}

}

var operators = []string{"==", "!=", ">=", "<=", ">", "<"}

func splitSimpleConditionExpression(expression string) (left, right, operator string, err error) {
	index := -1
	for _, _operator := range operators {
		index = strings.Index(expression, _operator)
		if index != -1 && !indexInQuota(index, expression) {
			operator = _operator
			break
		}
	}
	if index == -1 {
		return "", "", "", fmt.Errorf("invalid expression: %s", expression)
	}
	return strings.TrimSpace(expression[:index]), strings.TrimSpace(expression[index+len(operator):]), operator, nil
}

func parseSimpleCondition(expression string) (conditions.Condition, error) {
	left, right, operator, err := splitSimpleConditionExpression(expression)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(left, "VALUE_OF(") {
		return nil, fmt.Errorf("invalid expression: %s", expression)
	}
	if !strings.HasSuffix(left, ")") {
		return nil, fmt.Errorf("invalid expression: %s", expression)
	}

	key := left[9 : len(left)-1]
	value := unwrapQuotaIfNeeded(right)

	switch operator {
	case "==":
		return conditions.New().ValueOf(key).EqualTo(value), nil
	case "!=":
		return conditions.New().ValueOf(key).NotEqual(value), nil
	case ">":
		return conditions.New().ValueOf(key).GreaterThan(value), nil
	case ">=":
		return conditions.New().ValueOf(key).GreaterThanOrEqual(value), nil
	case "<":
		return conditions.New().ValueOf(key).LesserThan(value), nil
	case "<=":
		return conditions.New().ValueOf(key).LesserThanOrEqual(value), nil
	default:
		return nil, fmt.Errorf("invalid operator: %s", operator)
	}
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

func unwrapIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)
	if strings.HasPrefix(expression, "(") {
		index := -1
		braceNumber := 0
		for i, char := range expression {
			if char == '(' {
				braceNumber++
			} else if char == ')' {
				braceNumber--
			}
			if braceNumber < 0 {
				return expression
			}
			if braceNumber == 0 {
				index = i
				break
			}
		}

		if index == len(expression)-1 {
			expression = expression[1 : len(expression)-1]
		} else {
			return expression
		}
		return unwrapIfNeeded(expression)
	} else {
		return expression
	}
}

// TODO support quotation
func splitCombinationConditionExpression(expression string) (conditionStr string, operator string, rest string, err error) {
	expression = strings.TrimSpace(expression)
	if strings.HasPrefix(expression, "(") {
		braceNumber := 0
		index := 0
		for i, char := range expression {
			if char == '(' {
				braceNumber++
			} else if char == ')' {
				braceNumber--
			}
			if braceNumber < 0 {
				return "", "", "", fmt.Errorf("splitCombinationConditionExpression: invalid condition expression: %s", expression)
			} else if braceNumber == 0 {
				index = i
				break
			}
		}

		if braceNumber > 0 {
			return "", "", "", fmt.Errorf("splitCombinationConditionExpression: invalid condition expression: %s", expression)
		}

		head := strings.TrimSpace(expression[1:index])
		operatorAndTail := strings.TrimSpace(expression[index+1:])
		operator, rest, err = splitOperatorAndTail(operatorAndTail)
		if err != nil {
			return "", "", "", err
		}
		return head, operator, rest, nil
	} else {
		firstAndIndex := strings.Index(expression, "&&")
		firstOrIndex := strings.Index(expression, "||")

		if firstAndIndex < 0 {
			if firstOrIndex < 0 {
				return expression, "", "", nil
			} else if firstOrIndex == 0 {
				return "", "", "", fmt.Errorf("splitCombinationConditionExpression: invalid condition expression: %s", expression)
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
			return "", "", "", fmt.Errorf("splitCombinationConditionExpression: invalid condition expression: %s", expression)
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
				return "", "", "", fmt.Errorf("splitCombinationConditionExpression: invalid condition expression: %s", expression)
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
func splitOperatorAndTail(operatorAndTail string) (string, string, error) {
	if operatorAndTail == "" {
		return "", "", nil
	}

	if !strings.HasPrefix(operatorAndTail, "&&") && !strings.HasPrefix(operatorAndTail, "||") {
		return "", "", fmt.Errorf("invalid condition expression: invalid operator and tail: %s", operatorAndTail)
	}

	operator := operatorAndTail[:2]

	tail := strings.TrimSpace(operatorAndTail[2:])
	if tail == "" {
		return "", "", fmt.Errorf("invalid condition expression: invalid operator and tail: %s", operatorAndTail)
	}

	return operator, tail, nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
