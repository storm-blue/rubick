package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"reflect"
)

type lengthOfCondition struct {
	operator int
	key      string
	value    int
}

func (c *lengthOfCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *lengthOfCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *lengthOfCondition) Calculate(object objects.StructuredObject) (bool, error) {
	v, err := object.Get(c.key)
	if err != nil {
		return false, err
	}

	if v == nil {
		return false, nil
	}

	size := 0
	switch v_ := v.(type) {
	case []interface{}:
		size = len(v_)
	case map[interface{}]interface{}:
		size = len(v_)
	default:
		return false, fmt.Errorf("calculate error: unsupported type: %v", reflect.TypeOf(v))
	}

	return calculateNumber(c.operator, float64(size), c.value)
}

func (c *lengthOfCondition) String() string {
	operatorString := ""
	switch c.operator {
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

	return fmt.Sprintf("%v %v %v", c.key, operatorString, c.value)
}
