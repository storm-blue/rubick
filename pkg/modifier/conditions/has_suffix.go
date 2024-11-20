package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

type hasSuffixCondition struct {
	key    string
	suffix string
}

func (c *hasSuffixCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *hasSuffixCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *hasSuffixCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := object.GetString(c.key)
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(result, c.suffix), nil
}

func (c *hasSuffixCondition) String() string {
	return fmt.Sprintf("HAS_SUFFIX(%v, \"%v\")", c.key, c.suffix)
}
