package action

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

// -- replace part action --

func NewReplacePartAction(key, old, new string) Action {
	return &replacePartAction{key: key, old: old, _new: new}
}

type replacePartAction struct {
	key  string
	old  string
	_new string
}

func (a *replacePartAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := object.GetString(a.key)
	if err != nil {
		context.Log(object, a, err)
		return
	}

	v_ := strings.ReplaceAll(v, a.old, a._new)

	if err := object.Set(a.key, v_); err != nil {
		context.Log(object, a, err)
	}
}

func (a *replacePartAction) String() string {
	return fmt.Sprintf("replacePartAction: key=%v, old=%v, new=%v", a.key, a.old, a._new)
}

// -- trim prefix action --

func NewTrimPrefixAction(key string, prefix Valuable) Action {
	return &trimPrefixAction{key: key, prefix: prefix}
}

type trimPrefixAction struct {
	key    string
	prefix Valuable
}

func (a *trimPrefixAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := object.GetString(a.key)
	if err != nil {
		context.Log(object, a, err)
		return
	}

	argV, err := a.prefix.getValue(object)
	if err != nil {
		context.Log(object, a, err)
		return
	}
	if p, ok := argV.(string); ok {
		v_ := strings.TrimPrefix(v, p)

		if err := object.Set(a.key, v_); err != nil {
			context.Log(object, a, err)
		}
	} else {
		context.Log(object, a, fmt.Errorf("expected string, got %v", argV))
	}
}

func (a *trimPrefixAction) String() string {
	return fmt.Sprintf("TrimPrefixAction: key=%v, prefix=%v", a.key, a.prefix)
}

// -- trim suffix action --

func NewTrimSuffixAction(key string, suffix Valuable) Action {
	return &trimSuffixAction{key: key, suffix: suffix}
}

type trimSuffixAction struct {
	key    string
	suffix Valuable
}

func (a *trimSuffixAction) DoAction(context Context, object objects.StructuredObject) {
	v, err := object.GetString(a.key)
	if err != nil {
		context.Log(object, a, err)
		return
	}

	argV, err := a.suffix.getValue(object)
	if err != nil {
		context.Log(object, a, err)
		return
	}
	if p, ok := argV.(string); ok {
		v_ := strings.TrimSuffix(v, p)

		if err := object.Set(a.key, v_); err != nil {
			context.Log(object, a, err)
		}
	} else {
		context.Log(object, a, fmt.Errorf("expected string, got %v", argV))
	}
}

func (a *trimSuffixAction) String() string {
	return fmt.Sprintf("TrimPrefixAction: key=%v, prefix=%v", a.key, a.suffix)
}
