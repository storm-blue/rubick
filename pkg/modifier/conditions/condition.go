package conditions

import (
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type Condition interface {
	Calculate(objects.StructuredObject) (bool, error)
	String() string
	And(condition Condition) Condition
	Or(condition Condition) Condition
}
