package conditions

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type existsCondition struct {
	key string
}

func (e *existsCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     e,
		right:    condition,
		operator: And,
	}
}

func (e *existsCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     e,
		right:    condition,
		operator: Or,
	}
}

func (e *existsCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := object.Get(e.key)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

func (e *existsCondition) String() string {
	return fmt.Sprintf("EXISTS(%v)", e.key)
}
