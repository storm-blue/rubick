package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

type hasPrefixCondition struct {
	key    string
	prefix string
}

func (c *hasPrefixCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *hasPrefixCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *hasPrefixCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := object.GetString(c.key)
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(result, c.prefix), nil
}

func (c *hasPrefixCondition) String() string {
	return fmt.Sprintf("HAS_PREFIX(%v, \"%v\")", c.key, c.prefix)
}
