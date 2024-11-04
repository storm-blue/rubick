package action

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type Action interface {
	DoAction(context Context, object objects.StructuredObject)
	String() string
}

type Valuable interface {
	getValue(object objects.StructuredObject) (interface{}, error)
}

func Original(v interface{}) Valuable {
	return &originalValue{value: v}
}

type originalValue struct {
	value interface{}
}

func (o *originalValue) getValue(_ objects.StructuredObject) (interface{}, error) {
	return o.value, nil
}

func ValueOf(key string) Valuable {
	return &valueOfKeyValue{key: key}
}

type valueOfKeyValue struct {
	key string
}

func (o *valueOfKeyValue) getValue(object objects.StructuredObject) (interface{}, error) {
	return object.Get(o.key)
}

// -- delete action --

func NewDeleteAction(key string) Action {
	return &deleteAction{key: key}
}

type deleteAction struct {
	key string
}

func (d *deleteAction) DoAction(context Context, object objects.StructuredObject) {
	if err := object.Delete(d.key); err != nil {
		context.Log(object, d, err)
	}
}

func (d *deleteAction) String() string {
	return fmt.Sprintf("DeleteAction: key=%v", d.key)
}

// -- set action --

func NewSetAction(key string, value Valuable) Action {
	return &setAction{key: key, value: value}
}

type setAction struct {
	key   string
	value Valuable
}

func (s *setAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := s.value.getValue(object)
	if err != nil {
		context.Log(object, s, err)
		return
	}

	if err := object.Set(s.key, v); err != nil {
		context.Log(object, s, err)
	}
}

func (s *setAction) String() string {
	return fmt.Sprintf("SetAction: key=%v, value=%v", s.key, s.value)
}

// -- print action --

func NewPrintAction(key string) Action {
	return &printAction{key: key}
}

type printAction struct {
	key string
}

func (s *printAction) DoAction(context Context, object objects.StructuredObject) {
	if v, err := object.Get(s.key); err != nil {
		context.Log(object, s, err)
	} else {
		fmt.Printf("%v: %v\n", s.key, v)
	}
}

func (s *printAction) String() string {
	return fmt.Sprintf("PrintAction: key=%v", s.key)
}

// -- mark removed action --

func NewMarkRemovedAction() Action {
	return &markRemovedAction{}
}

type markRemovedAction struct{}

func (s *markRemovedAction) DoAction(_ Context, object objects.StructuredObject) {
	object.Metadata().MarkRemoved(true)
}

func (s *markRemovedAction) String() string {
	return fmt.Sprintf("MarkRemovedAction")
}
