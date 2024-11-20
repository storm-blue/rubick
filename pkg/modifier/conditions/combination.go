package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type CombinationCondition struct {
	left     Condition
	right    Condition
	operator int
}

func (c *CombinationCondition) Calculate(object objects.StructuredObject) (bool, error) {
	leftResult, err := c.left.Calculate(object)
	if err != nil {
		return false, err
	}
	rightResult, err := c.right.Calculate(object)
	if err != nil {
		return false, err
	}

	switch c.operator {
	case And:
		return leftResult && rightResult, nil
	case Or:
		return leftResult || rightResult, nil
	default:
		return false, fmt.Errorf("calculate error: unsupported operator: %v", c.operator)
	}
}

func (c *CombinationCondition) String() string {
	var leftString, rightString string
	if _, ok := c.left.(*CombinationCondition); ok {
		leftString = c.left.String()
	} else {
		leftString = fmt.Sprintf("(%v)", c.left.String())
	}

	rightString = c.right.String()

	switch c.operator {
	case And:
		return fmt.Sprintf("%v && (%v)", leftString, rightString)
	case Or:
		return fmt.Sprintf("%v || (%v)", leftString, rightString)
	default:
		return fmt.Sprintf("Format error: unsupported operator: %v", c.operator)
	}
}

func (c *CombinationCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *CombinationCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}
