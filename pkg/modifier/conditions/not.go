package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type notCondition struct {
	condition Condition
}

func (c *notCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *notCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *notCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := c.condition.Calculate(object)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (c *notCondition) String() string {
	return fmt.Sprintf("!(%v)", c.condition.String())
}
