package common

import "testing"

func TestUnwrapIfNeeded(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       string
	}{
		{
			name:       "TEST1",
			expression: "  ((x))",
			want:       "x",
		},
		{
			name:       "TEST2",
			expression: "(   ( (   aa( ))   ))",
			want:       "aa( )",
		},
		{
			name:       "TEST3",
			expression: "(   ( ()))))   aa( ))   ))",
			want:       "(   ( ()))))   aa( ))   ))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnwrapIfNeeded(tt.expression)
			if got != tt.want {
				t.Errorf("unwrapIfNeeded() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestFindFirstParenthesesPair(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantLeft   int
		wantRight  int
	}{
		{
			name:       "TEST",
			expression: "()()",
			wantLeft:   0,
			wantRight:  1,
		},
		{
			name:       "TEST",
			expression: ")()()",
			wantLeft:   1,
			wantRight:  2,
		},
		{
			name:       "TEST",
			expression: ")(()())",
			wantLeft:   1,
			wantRight:  6,
		},
		{
			name:       "TEST",
			expression: ")(()()",
			wantLeft:   -1,
			wantRight:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeft, gotRight := FindFirstParenthesesPair(tt.expression)
			if gotLeft != tt.wantLeft {
				t.Errorf("FindFirstParenthesesPair() gotLeft = %v, want %v", gotLeft, tt.wantLeft)
			}
			if gotRight != tt.wantRight {
				t.Errorf("FindFirstParenthesesPair() gotRight = %v, want %v", gotRight, tt.wantRight)
			}
		})
	}
}
