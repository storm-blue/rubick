package objects

import (
	"reflect"
	"testing"
)

func TestParseNextKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    YamlKey
		want1   string
		wantErr bool
	}{
		{
			name: "TEST_PARSE_BRACED_SIMPLE_OBJECT",
			key:  "(a)",
			want: YamlKey{
				key: "a",
			},
			want1:   "",
			wantErr: false,
		},
		{
			name: "TEST_PARSE_ARRAY",
			key:  "a[1].b.c",
			want: YamlKey{
				isArray: true,
				key:     "a",
				index:   YamlIndex{indexType: IndexNormal, index: 1},
			},
			want1:   "b.c",
			wantErr: false,
		},
		{
			name:    "TEST_PARSE_ARRAY_ERROR",
			key:     "a[1].",
			want:    YamlKey{},
			want1:   "",
			wantErr: true,
		},
		{
			name:    "TEST_PARSE_ARRAY_ERROR_2",
			key:     "a[1]-",
			want:    YamlKey{},
			want1:   "",
			wantErr: true,
		},
		{
			name: "TEST_PARSE_BRACED_ARRAY",
			key:  "(a.b.c)[1].b.c",
			want: YamlKey{
				isArray: true,
				key:     "a.b.c",
				index:   YamlIndex{indexType: IndexNormal, index: 1},
			},
			want1:   "b.c",
			wantErr: false,
		},
		{
			name:    "TEST_PARSE_BRACED_ARRAY_ERROR",
			key:     "(a.b.c)[1].",
			want:    YamlKey{},
			want1:   "",
			wantErr: true,
		},
		{
			name:    "TEST_PARSE_BRACED_ARRAY_ERROR_2",
			key:     "(a.b.c)[1]x",
			want:    YamlKey{},
			want1:   "",
			wantErr: true,
		},
		{
			name: "TEST_PARSE_BRACED_MAP",
			key:  "(a.b.c).d",
			want: YamlKey{
				key: "a.b.c",
			},
			want1:   "d",
			wantErr: false,
		},
		{
			name: "TEST_PARSE_MASS_1",
			key:  "c[1].key4[0].trouble.size",
			want: YamlKey{
				key:     "c",
				isArray: true,
				index:   YamlIndex{indexType: IndexNormal, index: 1},
			},
			want1:   "key4[0].trouble.size",
			wantErr: false,
		},
		{
			name: "TEST_PARSE_ARRAY_MAX",
			key:  "c[+].key4[0].trouble.size",
			want: YamlKey{
				key:     "c",
				isArray: true,
				index:   YamlIndex{indexType: IndexMax},
			},
			want1:   "key4[0].trouble.size",
			wantErr: false,
		},
		{
			name: "TEST_PARSE_ARRAY_APPEND",
			key:  "c[++].key4[0].trouble.size",
			want: YamlKey{
				key:     "c",
				isArray: true,
				index:   YamlIndex{indexType: IndexAppend},
			},
			want1:   "key4[0].trouble.size",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseNextKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNextKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNextKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseNextKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFromYAML(t *testing.T) {
	tests := []struct {
		name    string
		yamlStr string

		want    StructuredObject
		wantErr bool
	}{
		{
			name: "TEST1",
			yamlStr: `---
a: 1
`,
			want:    _object{"a": 1},
			wantErr: false,
		},
		{
			name: "TEST2",
			yamlStr: `---
a: 1.2
`,
			want:    _object{"a": 1.2},
			wantErr: false,
		},
		{
			name: "TEST3",
			yamlStr: `---
a: true
`,
			want:    _object{"a": true},
			wantErr: false,
		},
		{
			name: "TEST4",
			yamlStr: `---
a: 9223372036854775807
`,
			want:    _object{"a": 9223372036854775807},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromYAML(tt.yamlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromYAML() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lenientEqual(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		s    string
		want bool
	}{
		{
			name: "TEST1",
			v:    1.0,
			s:    "1.0",
			want: true,
		},
		{
			name: "TEST2",
			v:    1,
			s:    "1.0",
			want: true,
		},
		{
			name: "TEST3",
			v:    1.0,
			s:    "1",
			want: true,
		},
		{
			name: "TEST4",
			v:    "1.0",
			s:    "1.0",
			want: true,
		},
		{
			name: "TEST5",
			v:    "1",
			s:    "1.0",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lenientEqual(tt.v, tt.s); got != tt.want {
				t.Errorf("lenientEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
