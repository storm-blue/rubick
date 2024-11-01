package match

import (
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

type Matcher interface {
	Match(objects.StructuredObject) bool
}

func NewOrMather(matchers ...Matcher) Matcher {
	return &orMatcher{matchers}
}

type orMatcher struct {
	matchers []Matcher
}

func (l *orMatcher) Match(object objects.StructuredObject) bool {
	if len(l.matchers) == 0 {
		return false
	}
	for _, matcher := range l.matchers {
		if matcher.Match(object) {
			return true
		}
	}
	return false
}

func NewAndMatcher(matchers ...Matcher) Matcher {
	return &andMatcher{matchers}
}

type andMatcher struct {
	matchers []Matcher
}

func (l *andMatcher) Match(object objects.StructuredObject) bool {
	if len(l.matchers) == 0 {
		return false
	}
	for _, matcher := range l.matchers {
		if !matcher.Match(object) {
			return false
		}
	}
	return true
}

func NewStringMatcher(expression, key string) Matcher {
	return &stringMatcher{
		expression: expression,
		key:        key,
	}
}

type stringMatcher struct {
	expression string
	key        string
}

func (s *stringMatcher) Match(object objects.StructuredObject) bool {
	v, err := object.GetString(s.key)
	if err != nil {
		return false
	}

	return stringIsMatch(s.expression, v)
}

func stringIsMatch(expression, value string) bool {
	if expression == "*" {
		return true
	}
	return expression == value
}

func NewTrueMatcher() Matcher {
	return &trueMatcher{}
}

type trueMatcher struct {
	expression string
	key        string
}

func (s *trueMatcher) Match(_ objects.StructuredObject) bool {
	return true
}

func NewFalseMatcher() Matcher {
	return &falseMatcher{}
}

type falseMatcher struct {
	expression string
	key        string
}

func (s *falseMatcher) Match(_ objects.StructuredObject) bool {
	return false
}
