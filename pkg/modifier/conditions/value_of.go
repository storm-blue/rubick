package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strconv"
)

type valueOfCondition struct {
	operator int
	key      string
	value    interface{}
}

func (s *valueOfCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     s,
		right:    condition,
		operator: And,
	}
}

func (s *valueOfCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     s,
		right:    condition,
		operator: Or,
	}
}

func (s *valueOfCondition) Calculate(object objects.StructuredObject) (bool, error) {
	v, err := object.Get(s.key)
	if err != nil {
		return false, err
	}

	return calculate(s.operator, v, s.value)
}

func (s *valueOfCondition) String() string {
	operatorString := ""
	switch s.operator {
	case GreaterThan:
		operatorString = ">"
	case GreaterThanOrEqual:
		operatorString = ">="
	case LesserThan:
		operatorString = "<"
	case LesserThanOrEqual:
		operatorString = "<="
	case EqualTo:
		operatorString = "=="
	case NotEqual:
		operatorString = "!="
	default:
	}

	valueString := ""
	switch tv := s.value.(type) {
	case int8:
		valueString = strconv.Itoa(int(tv))
	case int32:
		valueString = strconv.Itoa(int(tv))
	case int64:
		valueString = strconv.Itoa(int(tv))
	case int:
		valueString = strconv.Itoa(tv)
	case uint:
		valueString = strconv.Itoa(int(tv))
	case float32:
		valueString = strconv.FormatFloat(float64(tv), 'g', 20, 64)
	case float64:
		valueString = strconv.FormatFloat(tv, 'g', 20, 64)
	case string:
		valueString = tv
	case []interface{}:
		valueString = "<array>"
	case map[interface{}]interface{}:
		valueString = "<object>"
	default:
		valueString = "<object>"
	}

	return fmt.Sprintf("%v %v %v", s.key, operatorString, valueString)
}
