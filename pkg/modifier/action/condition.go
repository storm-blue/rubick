package action

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

// -- condition action --

func NewConditionAction(condition conditions.Condition, action Action) Action {
	return &ConditionAction{
		condition: condition,
		action:    action,
	}
}

type ConditionAction struct {
	condition conditions.Condition
	action    Action
}

func (c *ConditionAction) DoAction(context Context, object objects.StructuredObject) {
	r, err := c.condition.Calculate(object)
	if err != nil {
		context.Log(object, c, err)
		return
	}

	if r {
		c.action.DoAction(context, object)
	}
}

func (c *ConditionAction) String() string {
	return fmt.Sprintf("ConditionAction: condition=%v, action=%v", c.condition.String(), c.action.String())
}
