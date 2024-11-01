package match

import (
	"testing"
)

func Test_orMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		matchers []Matcher
		want     bool
	}{
		{
			name:     "TEST1",
			matchers: []Matcher{NewFalseMatcher(), NewTrueMatcher()},
			want:     true,
		},
		{
			name:     "TEST2",
			matchers: []Matcher{NewFalseMatcher(), NewFalseMatcher()},
			want:     false,
		},
		{
			name:     "TEST3",
			matchers: []Matcher{NewTrueMatcher(), NewTrueMatcher()},
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewOrMather(tt.matchers...)
			if got := l.Match(nil); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_andMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		matchers []Matcher
		want     bool
	}{
		{
			name:     "TEST1",
			matchers: []Matcher{NewFalseMatcher(), NewTrueMatcher()},
			want:     false,
		},
		{
			name:     "TEST2",
			matchers: []Matcher{NewFalseMatcher(), NewFalseMatcher()},
			want:     false,
		},
		{
			name:     "TEST3",
			matchers: []Matcher{NewTrueMatcher(), NewTrueMatcher()},
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewAndMatcher(tt.matchers...)
			if got := l.Match(nil); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringIsMatch(t *testing.T) {

	tests := []struct {
		name       string
		expression string
		value      string
		want       bool
	}{
		{
			name:       "TEST1",
			expression: "*",
			value:      "123",
			want:       true,
		},
		{
			name:       "TEST2",
			expression: "234",
			value:      "123",
			want:       false,
		},
		{
			name:       "TEST3",
			expression: "234",
			value:      "*",
			want:       false,
		},
		{
			name:       "TEST4",
			expression: "12*",
			value:      "123",
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringIsMatch(tt.expression, tt.value); got != tt.want {
				t.Errorf("stringIsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
