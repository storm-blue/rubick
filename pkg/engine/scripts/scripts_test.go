package scripts

import (
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/conditions"
	"reflect"
	"testing"
)

func Test_splitExpression(t *testing.T) {
	tests := []struct {
		name               string
		expression         string
		wantConditionPart  string
		wantPureActionPart string
		wantErr            bool
	}{
		{
			name:               "TEST1",
			expression:         "IF VALUE_OF(a.b.c)==true THEN DELETE(z.x.f)",
			wantConditionPart:  "VALUE_OF(a.b.c)==true",
			wantPureActionPart: "DELETE(z.x.f)",
			wantErr:            false,
		},
		{
			name:               "TEST2",
			expression:         "IF VALUE_OF(a.b.c)==true THEN",
			wantConditionPart:  "",
			wantPureActionPart: "",
			wantErr:            true,
		},
		{
			name:               "TEST3",
			expression:         "IF   VALUE_OF(a.b.c)==true     THEN     DELETE(a.b.c)",
			wantConditionPart:  "VALUE_OF(a.b.c)==true",
			wantPureActionPart: "DELETE(a.b.c)",
			wantErr:            false,
		},
		{
			name:               "TEST4",
			expression:         "     DELETE(a.b.c)",
			wantConditionPart:  "",
			wantPureActionPart: "DELETE(a.b.c)",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConditionPart, gotPureActionPart, err := splitExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConditionPart != tt.wantConditionPart {
				t.Errorf("splitExpression() gotConditionPart = %v, want %v", gotConditionPart, tt.wantConditionPart)
			}
			if gotPureActionPart != tt.wantPureActionPart {
				t.Errorf("splitExpression() gotPureActionPart = %v, want %v", gotPureActionPart, tt.wantPureActionPart)
			}
		})
	}
}

func Test_splitCombinationConditionExpression(t *testing.T) {
	tests := []struct {
		name             string
		expression       string
		wantConditionStr string
		wantOperator     string
		wantRest         string
		wantErr          bool
	}{
		{
			name:             "TEST1",
			expression:       "(VALUE_OF(a.b.c)==12)||VALUE_OF(z.x)==why",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "||",
			wantRest:         "VALUE_OF(z.x)==why",
			wantErr:          false,
		},
		{
			name:             "TEST2",
			expression:       "VALUE_OF(a.b.c)==12||VALUE_OF(z.x)==why",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "||",
			wantRest:         "VALUE_OF(z.x)==why",
			wantErr:          false,
		},
		{
			name:             "TEST3",
			expression:       "VALUE_OF(a.b.c)==12||VALUE_OF(z.x)==why&&VALUE_OF(x.c)==shit",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "||",
			wantRest:         "VALUE_OF(z.x)==why&&VALUE_OF(x.c)==shit",
			wantErr:          false,
		},
		{
			name:             "TEST4",
			expression:       "VALUE_OF(a.b.c)==12&&VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "&&",
			wantRest:         "VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit",
			wantErr:          false,
		},
		{
			name:             "TEST5",
			expression:       "VALUE_OF(a.b.c)==12  &&    VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit    ",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "&&",
			wantRest:         "VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit",
			wantErr:          false,
		},
		{
			name:             "TEST6",
			expression:       "(VALUE_OF(a.b.c)==12   ) &&   (    VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit)    ",
			wantConditionStr: "VALUE_OF(a.b.c)==12",
			wantOperator:     "&&",
			wantRest:         "(    VALUE_OF(z.x)==why||VALUE_OF(x.c)==shit)",
			wantErr:          false,
		},
		{
			name:             "TEST7",
			expression:       "(VALUE_OF(a.b.c)==12   ) &&   ",
			wantConditionStr: "",
			wantOperator:     "",
			wantRest:         "",
			wantErr:          true,
		},
		{
			name:             "TEST8",
			expression:       "(||VALUE_OF(a.b.c)==12   )    ",
			wantConditionStr: "||VALUE_OF(a.b.c)==12",
			wantOperator:     "",
			wantRest:         "",
			wantErr:          false,
		},
		{
			name:             "TEST9",
			expression:       "||(||VALUE_OF(a.b.c)==12   )    ",
			wantConditionStr: "",
			wantOperator:     "",
			wantRest:         "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConditionStr, gotOperator, gotRest, err := splitCombinationConditionExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitCombinationConditionExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConditionStr != tt.wantConditionStr {
				t.Errorf("splitCombinationConditionExpression() gotConditionStr = %v, want %v", gotConditionStr, tt.wantConditionStr)
			}
			if gotOperator != tt.wantOperator {
				t.Errorf("splitCombinationConditionExpression() gotOperator = %v, want %v", gotOperator, tt.wantOperator)
			}
			if gotRest != tt.wantRest {
				t.Errorf("splitCombinationConditionExpression() gotRest = %v, want %v", gotRest, tt.wantRest)
			}
		})
	}
}

func Test_indexInQuota(t *testing.T) {
	tests := []struct {
		name       string
		index      int
		expression string
		want       bool
	}{
		{
			name:       "TEST1",
			index:      3,
			expression: `my "baby"`,
			want:       true,
		},
		{
			name:       "TEST2",
			index:      4,
			expression: `my "baby"`,
			want:       true,
		},
		{
			name:       "TEST3",
			index:      5,
			expression: `my "baby"`,
			want:       true,
		},
		{
			name:       "TEST4",
			index:      2,
			expression: `my "baby"`,
			want:       false,
		},
		{
			name:       "TEST5",
			index:      10,
			expression: `my "baby",i love "you"`,
			want:       false,
		},
		{
			name:       "TEST6",
			index:      17,
			expression: `my "baby",i love "you"`,
			want:       true,
		},
		{
			name:       "TEST7",
			index:      22,
			expression: `my "baby",i love "you"`,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexInQuota(tt.index, tt.expression); got != tt.want {
				t.Errorf("indexInQuota() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitRelationalSimpleConditionExpression(t *testing.T) {
	tests := []struct {
		name         string
		expression   string
		wantLeft     string
		wantRight    string
		wantOperator string
		wantErr      bool
	}{
		{
			name:         "TEST1",
			expression:   `a==b`,
			wantLeft:     "a",
			wantRight:    "b",
			wantOperator: "==",
			wantErr:      false,
		},
		{
			name:         "TEST2",
			expression:   `a  <==   b d`,
			wantLeft:     "a  <",
			wantRight:    "b d",
			wantOperator: "==",
			wantErr:      false,
		},
		{
			name:         "TEST3",
			expression:   `a  <="==   b d"`,
			wantLeft:     "a",
			wantRight:    `"==   b d"`,
			wantOperator: "<=",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeft, gotRight, gotOperator, err := splitRelationalSimpleConditionExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitRelationalSimpleConditionExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLeft != tt.wantLeft {
				t.Errorf("splitRelationalSimpleConditionExpression() gotLeft = %v, want %v", gotLeft, tt.wantLeft)
			}
			if gotRight != tt.wantRight {
				t.Errorf("splitRelationalSimpleConditionExpression() gotRight = %v, want %v", gotRight, tt.wantRight)
			}
			if gotOperator != tt.wantOperator {
				t.Errorf("splitRelationalSimpleConditionExpression() gotOperator = %v, want %v", gotOperator, tt.wantOperator)
			}
		})
	}
}

func Test_splitPureActionExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantMethod string
		wantArgs   []string
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "ABC(a,b,c)",
			wantMethod: "ABC",
			wantArgs:   []string{"a", "b", "c"},
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: "ABC(a , b,    c)",
			wantMethod: "ABC",
			wantArgs:   []string{"a", "b", "c"},
			wantErr:    false,
		},
		{
			name:       "TEST3",
			expression: `ABC(a , b,"    c")`,
			wantMethod: "ABC",
			wantArgs:   []string{"a", "b", `"    c"`},
			wantErr:    false,
		},
		{
			name:       "TEST4",
			expression: `ABC()`,
			wantMethod: "ABC",
			wantArgs:   nil,
			wantErr:    false,
		},
		{
			name:       "TEST4",
			expression: `ABC(a,"adb,")`,
			wantMethod: "ABC",
			wantArgs:   []string{"a", `"adb`, `"`},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotArgs, err := splitPureActionExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitPureActionExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("splitPureActionExpression() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("splitPureActionExpression() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func Test_isSimpleCondition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{
			name:       "TEST1",
			expression: "(VALUE_OF(z.x.b)==3)",
			want:       false,
		},
		{
			name:       "TEST2",
			expression: "VALUE_OF(z.x.b)==3)",
			want:       true,
		},
		{
			name:       "TEST3",
			expression: "VALUE_OF(z.x.b)==3)||VALUE_OF(a)>5",
			want:       false,
		},
		{
			name:       "TEST4",
			expression: "(VALUE_OF(z.x.b)==3)|| VALUE_OF(a)>5",
			want:       false,
		},
		{
			name:       "TEST5",
			expression: "(VALUE_OF(z.x.b)==3)VALUE_OF(a)>5",
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSimpleCondition(tt.expression); got != tt.want {
				t.Errorf("isSimpleCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePureAction(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       action.Action
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "DELETE(a.b.c)",
			want:       action.NewDeleteAction("a.b.c"),
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: "DELETE(a.b.c,z.yz)",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "TEST3",
			expression: "SET(a.b.c,   z.yz)",
			want:       action.NewSetAction("a.b.c", action.Original("z.yz")),
			wantErr:    false,
		},
		{
			name:       "TEST4",
			expression: `SET(a.b.c, VALUE_OF("z.yz"))`,
			want:       action.NewSetAction("a.b.c", action.ValueOf("z.yz")),
			wantErr:    false,
		},
		{
			name:       "TEST5",
			expression: `SET(a.b.c, VALUE_OF("  z.yz"))`,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "TEST6",
			expression: `SET(a.b.c, VALUE_OF(  z.yz))`,
			want:       action.NewSetAction("a.b.c", action.ValueOf("z.yz")),
			wantErr:    false,
		},
		{
			name:       "TEST7",
			expression: `SET(a.b.c, VALUE_OF(z.yz))`,
			want:       action.NewSetAction("a.b.c", action.ValueOf("z.yz")),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePureAction(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePureAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePureAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseCondition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       conditions.Condition
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "(VALUE_OF(x.y.z)==3 && VALUE_OF(a.b.c)>=18)",
			want:       conditions.New().ValueOf("x.y.z").EqualTo("3").And(conditions.New().ValueOf("a.b.c").GreaterThanOrEqual("18")),
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: "((((VALUE_OF(x.y.z)==3)) && (VALUE_OF(a.b.c)>=18)))",
			want:       conditions.New().ValueOf("x.y.z").EqualTo("3").And(conditions.New().ValueOf("a.b.c").GreaterThanOrEqual("18")),
			wantErr:    false,
		},
		{
			name:       "TEST3",
			expression: "((((VALUE_OF(x.y.z)==3)) && (VALUE_OF(a.b.c)>=18 || VALUE_OF(a.b.c[0]) == shit )))",
			want:       conditions.New().ValueOf("x.y.z").EqualTo("3").And(conditions.New().ValueOf("a.b.c").GreaterThanOrEqual("18").Or(conditions.New().ValueOf("a.b.c[0]").EqualTo("shit"))),
			wantErr:    false,
		},
		{
			name:       "TEST4",
			expression: "((((VALUE_OF(x.y.z)==3)) && (VALUE_OF(a.b.c)>=18 || VALUE_OF(a.b.c[0]) == shit ) )))",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "TEST5",
			expression: "((((VALUE_OF(x.y.z)==3)) && (VALUE_OF(a.b.c)>=18 || VALUE_OF(a.b.c[0]) == shit ))) || (VALUE_OF(z)<1.2 && (VALUE_OF(xxx) <= shit))",
			want: conditions.New().ValueOf("x.y.z").EqualTo("3").And(conditions.New().ValueOf("a.b.c").GreaterThanOrEqual("18").Or(conditions.New().ValueOf("a.b.c[0]").EqualTo("shit"))).
				Or(conditions.New().ValueOf("z").LesserThan("1.2").And(conditions.New().ValueOf("xxx").LesserThanOrEqual("shit"))),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCondition(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseValueOfSimpleCondition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       conditions.Condition
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "VALUE_OF(a.b.c)==3",
			want:       conditions.New().ValueOf("a.b.c").EqualTo("3"),
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: `VALUE_OF(a.b.c)=="shit"`,
			want:       conditions.New().ValueOf("a.b.c").EqualTo("shit"),
			wantErr:    false,
		},
		{
			name:       "TEST3",
			expression: `(VALUE_OF(a.b.c))=="shit"`,
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValueOfSimpleCondition(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSimpleCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSimpleCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseToNumberIfPossible(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want interface{}
	}{
		{
			name: "TEST1",
			arg:  "1",
			want: 1,
		},
		{
			name: "TEST2",
			arg:  "1.0",
			want: 1.0,
		},
		{
			name: "TEST3",
			arg:  `"1.0"`,
			want: "1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseToNumberIfPossible(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseToNumberIfPossible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExistsSimpleCondition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       conditions.Condition
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "EXISTS(a.b.c)",
			want:       conditions.New().Exists("a.b.c"),
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: "EXISTS( a.b.c )",
			want:       conditions.New().Exists("a.b.c"),
			wantErr:    false,
		},
		{
			name:       "TEST3",
			expression: "EXISTS( a.b.c)",
			want:       conditions.New().Exists("a.b.c"),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExistsSimpleCondition(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExistsSimpleCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExistsSimpleCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNotExistsSimpleCondition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       conditions.Condition
		wantErr    bool
	}{
		{
			name:       "TEST1",
			expression: "NOT_EXISTS(a.b.c)",
			want:       conditions.New().Not(conditions.New().Exists("a.b.c")),
			wantErr:    false,
		},
		{
			name:       "TEST2",
			expression: "NOT_EXISTS( a.b.c )",
			want:       conditions.New().Not(conditions.New().Exists("a.b.c")),
			wantErr:    false,
		},
		{
			name:       "TEST3",
			expression: "NOT_EXISTS( a.b.c)",
			want:       conditions.New().Not(conditions.New().Exists("a.b.c")),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNotExistsSimpleCondition(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNotExistsSimpleCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNotExistsSimpleCondition() got = %v, want %v", got, tt.want)
			}
		})
	}
}
