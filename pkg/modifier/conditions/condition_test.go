package conditions

import (
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		want      Condition
	}{
		{
			name:      "TEST_SIMPLE_1",
			condition: New().ValueOf("a").LesserThan(1),
			want: &simpleCondition{
				operator: LesserThan,
				key:      "a",
				value:    1,
			},
		},
		{
			name:      "TEST_SIMPLE_2",
			condition: New().ValueOf("a").LesserThan("simple"),
			want: &simpleCondition{
				operator: LesserThan,
				key:      "a",
				value:    "simple",
			},
		},
		{
			name:      "TEST_NOT",
			condition: New().Not(New().ValueOf("a").LesserThan(1)),
			want: &notCondition{
				condition: &simpleCondition{operator: LesserThan, key: "a", value: 1},
			},
		},
		{
			name: "TEST_MULTI_1",
			condition: New().ValueOf("a").LesserThan(1).
				And(New().ValueOf("b").GreaterThan(2)).
				And(New().ValueOf("c").EqualTo(3)),
			want: &CombinationCondition{
				left: &CombinationCondition{
					left:     &simpleCondition{operator: LesserThan, key: "a", value: 1},
					right:    &simpleCondition{operator: GreaterThan, key: "b", value: 2},
					operator: And,
				},
				right:    &simpleCondition{operator: EqualTo, key: "c", value: 3},
				operator: And,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculate(t *testing.T) {

	tests := []struct {
		name           string
		operator       int
		objectValue    interface{}
		conditionValue interface{}
		want           bool
		wantErr        bool
	}{
		{
			name:           "TEST_NUMBER_EQ_NUMBER",
			operator:       EqualTo,
			objectValue:    130,
			conditionValue: 130,
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_EQ_STRING",
			operator:       EqualTo,
			objectValue:    130,
			conditionValue: "130",
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_GT_NUMBER",
			operator:       GreaterThan,
			objectValue:    132,
			conditionValue: 130,
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_GT_STRING",
			operator:       GreaterThan,
			objectValue:    132,
			conditionValue: "130",
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_LT_NUMBER",
			operator:       LesserThan,
			objectValue:    132,
			conditionValue: 130,
			want:           false,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_LT_STRING",
			operator:       LesserThan,
			objectValue:    132,
			conditionValue: "130",
			want:           false,
			wantErr:        false,
		},
		{
			name:           "TEST_NUMBER_EQ_OBJECT_SHOULD_ERROR",
			operator:       LesserThan,
			objectValue:    132,
			conditionValue: map[interface{}]interface{}{},
			want:           false,
			wantErr:        true,
		},
		{
			name:           "TEST_STRING_EQ_STRING",
			operator:       EqualTo,
			objectValue:    "130",
			conditionValue: "130",
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_STRING_LT_STRING_SHOULD_ERROR",
			operator:       LesserThan,
			objectValue:    "130",
			conditionValue: "132",
			want:           false,
			wantErr:        true,
		},
		{
			name:           "TEST_ARRAY_EQ_ARRAY",
			operator:       EqualTo,
			objectValue:    []interface{}{1, 2, 3},
			conditionValue: []interface{}{1, 2, 3},
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_ARRAY_EQ_STRING",
			operator:       EqualTo,
			objectValue:    []interface{}{1, 2, 3},
			conditionValue: "[1,2,3]",
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_ARRAY_LT_SHOULD_ERROR",
			operator:       LesserThan,
			objectValue:    []interface{}{1, 2, 3},
			conditionValue: "[1,2,3]",
			want:           false,
			wantErr:        true,
		},
		{
			name:           "TEST_MAP_EQ_MAP",
			operator:       EqualTo,
			objectValue:    map[interface{}]interface{}{"a": 1, "b": 2, "c": 3},
			conditionValue: map[interface{}]interface{}{"a": 1, "b": 2, "c": 3},
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_MAP_EQ_STRING",
			operator:       EqualTo,
			objectValue:    map[interface{}]interface{}{"a": 1, "b": 2, "c": 3},
			conditionValue: `{"a":1,"b":2,"c":3}`,
			want:           true,
			wantErr:        false,
		},
		{
			name:           "TEST_MAP_LT_SHOULD_ERROR",
			operator:       LesserThan,
			objectValue:    map[interface{}]interface{}{"a": 1, "b": 2, "c": 3},
			conditionValue: "[1,2,3]",
			want:           false,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculate(tt.operator, tt.objectValue, tt.conditionValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleCondition_String(t *testing.T) {
	type fields struct {
		operator int
		key      string
		value    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TEST_SIMPLE_EQUAL",
			fields: fields{
				operator: EqualTo,
				key:      "a.b.c",
				value:    3,
			},
			want: "a.b.c == 3",
		},
		{
			name: "TEST_SIMPLE_NOT_EQUAL",
			fields: fields{
				operator: NotEqual,
				key:      "a.b.c",
				value:    3,
			},
			want: "a.b.c != 3",
		},
		{
			name: "TEST_SIMPLE_NOT_EQUAL_OBJECT",
			fields: fields{
				operator: NotEqual,
				key:      "a.b.c",
				value:    map[string]string{},
			},
			want: "a.b.c != <object>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &simpleCondition{
				operator: tt.fields.operator,
				key:      tt.fields.key,
				value:    tt.fields.value,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notCondition_String(t *testing.T) {
	type fields struct {
		condition Condition
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TEST1",
			fields: fields{
				condition: New().ValueOf("a.b.c").EqualTo("3"),
			},
			want: "!(a.b.c == 3)",
		},
		{
			name: "TEST2",
			fields: fields{
				condition: New().Not(New().ValueOf("a.b.c").EqualTo("3")),
			},
			want: "!(!(a.b.c == 3))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &notCondition{
				condition: tt.fields.condition,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombinationCondition_String(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		want      string
	}{
		{
			name: "TEST1",
			condition: New().ValueOf("a.b.c").NotEqual(map[interface{}]interface{}{}).
				Or(New().ValueOf("a.b.c").NotEqual("7")).
				Or(New().Not(New().ValueOf("a.b.c").LesserThan("30"))),
			want: "(a.b.c != <object>) || (a.b.c != 7) || (!(a.b.c < 30))",
		},
		{
			name: "TEST2",
			condition: New().ValueOf("a.b.c").NotEqual(map[interface{}]interface{}{}).
				Or(New().ValueOf("a.b.c").NotEqual("7").And(New().ValueOf("x.y.z").LesserThan("80"))).
				Or(New().Not(New().ValueOf("a.b.c").LesserThan("30"))),
			want: "(a.b.c != <object>) || ((a.b.c != 7) && (x.y.z < 80)) || (!(a.b.c < 30))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.condition
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_existsCondition_Calculate(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		object  objects.StructuredObject
		want    bool
		wantErr bool
	}{
		{
			name:    "TEST1",
			key:     "d",
			object:  objects.FromMap(map[interface{}]interface{}{"a": 1, "b": 2, "c": 3}),
			want:    false,
			wantErr: false,
		},
		{
			name:    "TEST2",
			key:     "c",
			object:  objects.FromMap(map[interface{}]interface{}{"a": 1, "b": 2, "c": "aa"}),
			want:    true,
			wantErr: false,
		},
		{
			name:    "TEST3",
			key:     "c",
			object:  objects.FromMap(map[interface{}]interface{}{"a": 1, "b": 2, "c": nil}),
			want:    false,
			wantErr: false,
		},
		{
			name:    "TEST3",
			key:     "c",
			object:  objects.FromMap(map[interface{}]interface{}{"a": 1, "b": 2, "c": "null"}),
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &existsCondition{
				key: tt.key,
			}
			got, err := e.Calculate(tt.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
