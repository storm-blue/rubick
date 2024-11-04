package modifier

import (
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"testing"
)

func TestYamlModifier_ModifyYAML(t *testing.T) {
	tests := []struct {
		name    string
		yamlStr string
		config  []action.Action
		want    string
		wantErr bool
	}{
		{
			name: "TEST_DELETE_ARRAY_ELE",
			yamlStr: `a:
  b: 300
  c:
  - key1: shit
    key2: git
    key3: mud
  - key1: shit2
    key2: git2
    key3: mud2
b: gala
`,
			config: []action.Action{
				action.NewDeleteAction("a.c[0]"),
			},
			want: `a:
  b: 300
  c:
  - key1: shit2
    key2: git2
    key3: mud2
b: gala
`,
			wantErr: false,
		},
		{
			name: "TEST_ADD_ARRAY",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewSetAction("a.c[++].key1", action.Original("shut")),
			},
			want: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
  - key1: shut
`,
			wantErr: false,
		},
		{
			name: "TEST_MODIFY_ARRAY",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewSetAction("a.c[0].key1", action.Original("shut")),
				action.NewSetAction("a.c[0].key4", action.Original("card")),
			},
			want: `a:
  c:
  - key1: shut
    key2: git
    key3: mud
    key4: card
`,
			wantErr: false,
		},
		{
			name: "TEST_INDEX_GREAT_THAN_ARRAY_SIZE_SHOULD_ERROR",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewSetAction("a.c[2].key1", action.Original("shut")),
				action.NewSetAction("a.c[2].key4", action.Original("card")),
			},
			want: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			wantErr: true,
		},
		{
			name: "TEST_RECURSIVE_SET_1",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewSetAction("a.c[1].key1", action.Original("shut")),
				action.NewSetAction("a.c[1].key4[0].trouble.size", action.Original(36)),
			},
			want: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
  - key1: shut
    key4:
    - trouble:
        size: 36
`,
			wantErr: false,
		},
		{
			name: "TEST_RECURSIVE_SET_2",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewSetAction("a.c[++].key1", action.Original("shut")),
				action.NewSetAction("a.c[1].key4[0].trouble.size", action.Original(36)),
			},
			want: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
  - key1: shut
    key4:
    - trouble:
        size: 36
`,
			wantErr: false,
		},
		{
			name: "TEST_SEARCH_DELETE",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewDeleteAction("a.c[key2=git]"),
			},
			want: `a:
  c: []
`,
			wantErr: false,
		},
		{
			name: "TEST_SEARCH_DELETE_2",
			yamlStr: `a:
  c:
  - key1: shit
    key2: git
    key3: mud
`,
			config: []action.Action{
				action.NewDeleteAction("a.c[key2=git].key3"),
			},
			want: `a:
  c:
  - key1: shit
    key2: git
`,
			wantErr: false,
		},
		{
			name: "TEST_LOOP_DELETE_1",
			yamlStr: `a:
  c:
  - key1: shit1
    key2: git1
    key3: mud1
  - key1: shit2
    key2: git2
    key3: mud2
  - key1: shit3
    key2: git3
    key3: mud3
`,
			config: []action.Action{
				action.NewDeleteAction("a.c[*].key3"),
			},
			want: `a:
  c:
  - key1: shit1
    key2: git1
  - key1: shit2
    key2: git2
  - key1: shit3
    key2: git3
`,
			wantErr: false,
		},
		{
			name: "TEST_LOOP_DELETE_2",
			yamlStr: `a:
  c:
  - key1: shit1
    key2: git1
    key3: mud1
  - key1: shit2
    key2: git2
    key3: mud2
  - key1: shit3
    key2: git3
    key3: mud3
`,
			config: []action.Action{
				action.NewDeleteAction("a.c[*]"),
			},
			want: `a:
  c: []
`,
			wantErr: false,
		},
		{
			name: "TEST_LOOP_SET_1",
			yamlStr: `a:
  c:
  - key1: shit1
    key2: git1
    key3: mud1
  - key1: shit2
    key2: git2
    key3: mud2
  - key1: shit3
    key2: git3
    key3: mud3
`,
			config: []action.Action{
				action.NewSetAction("a.c[*].key3", action.Original("cat")),
			},
			want: `a:
  c:
  - key1: shit1
    key2: git1
    key3: cat
  - key1: shit2
    key2: git2
    key3: cat
  - key1: shit3
    key2: git3
    key3: cat
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := action.NewContext([]string{})
			got, err := ModifyYAML(context, tt.yamlStr, tt.config)
			if ((err != nil) != tt.wantErr) && ((len(context.Logs()) > 0) != tt.wantErr) {
				t.Errorf("Modify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Modify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
