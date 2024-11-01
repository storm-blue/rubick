package action

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

type Action interface {
	DoAction(context Context, object objects.StructuredObject)
	String() string
}

// -- delete action --

func NewDeleteAction(key string) Action {
	return &DeleteAction{Key: key}
}

type DeleteAction struct {
	Key string
}

func (d *DeleteAction) DoAction(context Context, object objects.StructuredObject) {
	if err := object.Delete(d.Key); err != nil {
		context.Log(object, d, err)
	}
}

func (d *DeleteAction) String() string {
	return fmt.Sprintf("DeleteAction: key=%v", d.Key)
}

// -- set action --

func NewSetAction(key string, value interface{}) Action {
	return &SetAction{Key: key, Value: value}
}

type SetAction struct {
	Key   string
	Value interface{}
}

func (s *SetAction) DoAction(context Context, object objects.StructuredObject) {
	if err := object.Set(s.Key, s.Value); err != nil {
		context.Log(object, s, err)
	}
}

func (s *SetAction) String() string {
	return fmt.Sprintf("DeleteAction: key=%v, value=%v", s.Key, s.Value)
}

// -- set with value of action --

func NewSetWithValueOfAction(key, valueOf string) Action {
	return &SetWithValueOfAction{Key: key, ValueOf: valueOf}
}

type SetWithValueOfAction struct {
	Key     string
	ValueOf string
}

func (s *SetWithValueOfAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := object.Get(s.ValueOf)
	if err != nil {
		context.Log(object, s, err)
		return
	}

	if v != nil {
		if err = object.Set(s.Key, v); err != nil {
			context.Log(object, s, err)
		}
	}
}

func (s *SetWithValueOfAction) String() string {
	return fmt.Sprintf("SetWithValueOfAction: key=%v, valueOf=%v", s.Key, s.ValueOf)
}

// -- replace part action --

func NewReplacePartAction(key, old, new string) Action {
	return &ReplacePartAction{Key: key, Old: old, New: new}
}

type ReplacePartAction struct {
	Key string
	Old string
	New string
}

func (r *ReplacePartAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := object.GetString(r.Key)
	if err != nil {
		context.Log(object, r, err)
		return
	}

	v_ := strings.ReplaceAll(v, r.Old, r.New)

	if err := object.Set(r.Key, v_); err != nil {
		context.Log(object, r, err)
	}
}

func (r *ReplacePartAction) String() string {
	return fmt.Sprintf("ReplacePartAction: key=%v, Old=%v, New=%v", r.Key, r.Old, r.New)
}

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
