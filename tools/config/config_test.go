package config

import (
	"github.com/storm-blue/rubick/pkg/engine/match"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         string
		wantKubeconfig string
		wantScripts    string
		wantResources  map[string]match.Matcher
		wantErr        bool
	}{
		{
			name: "TEST1",
			config: `
[__kubeconfig__]
/.kube/config
[deployment]

[service]
java-dev/*
java-qa1/*
java-qa2/*
java-sit/*

[__scripts__]
DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)

# service scripts
DELETE(metadata.uid)
DELETE(status)
`,
			wantKubeconfig: "/.kube/config",
			wantScripts: `DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)
# service scripts
DELETE(metadata.uid)
DELETE(status)`,
			wantResources: map[string]match.Matcher{
				"service": match.NewOrMather(
					match.NewAndMatcher(
						match.NewStringMatcher("java-dev", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-qa1", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-qa2", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-sit", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
				),
			},
			wantErr: false,
		},
		{
			name: "TEST2",
			config: `
[__kubeconfig__]
/.kube/config
[deployment]
*/*
[service]
java-dev/*
java-qa1/*
java-qa2/*
java-sit/*

[__scripts__]
DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)

# service scripts
DELETE(metadata.uid)
DELETE(status)
`,
			wantKubeconfig: "/.kube/config",
			wantScripts: `DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)
# service scripts
DELETE(metadata.uid)
DELETE(status)`,
			wantResources: map[string]match.Matcher{
				"deployment": match.NewOrMather(
					match.NewAndMatcher(
						match.NewStringMatcher("*", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
				),
				"service": match.NewOrMather(
					match.NewAndMatcher(
						match.NewStringMatcher("java-dev", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-qa1", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-qa2", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
					match.NewAndMatcher(
						match.NewStringMatcher("java-sit", "metadata.namespace"),
						match.NewStringMatcher("*", "metadata.name"),
					),
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKubeconfig, gotScripts, gotResources, err := ParseConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKubeconfig != tt.wantKubeconfig {
				t.Errorf("ParseConfig() gotKubeconfig = %v, want %v", gotKubeconfig, tt.wantKubeconfig)
			}
			if !reflect.DeepEqual(gotScripts, tt.wantScripts) {
				t.Errorf("ParseConfig() gotScripts = %v, want %v", gotScripts, tt.wantScripts)
			}
			if !reflect.DeepEqual(gotResources, tt.wantResources) {
				t.Errorf("ParseConfig() gotResources = %v, want %v", gotResources, tt.wantResources)
			}
		})
	}
}

func Test_isValidResourceTypeString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "TEST1",
			s:    "abc.github.io",
			want: true,
		},
		{
			name: "TEST2",
			s:    "abc./github.io",
			want: false,
		},
		{
			name: "TEST3",
			s:    "abc.github123.io",
			want: true,
		},
		{
			name: "TEST3",
			s:    "abc.github-123.io",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidResourceTypeString(tt.s); got != tt.want {
				t.Errorf("isValidResourceTypeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "TEST1",
			path: "/a",
			want: true,
		},
		{
			name: "TEST2",
			path: "/   /",
			want: false,
		},
		{
			name: "TEST3",
			path: "a/",
			want: true,
		},
		{
			name: "TEST4",
			path: "/a/b c/.kube/config",
			want: false,
		},
		{
			name: "TEST5",
			path: "/a/b/c/.kube/config",
			want: true,
		},
		{
			name: "TEST6",
			path: "/.kube/config",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPath(tt.path); got != tt.want {
				t.Errorf("isValidPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidResourceExpression(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "TEST1",
			s:    "*/*",
			want: true,
		},
		{
			name: "TEST2",
			s:    "*a/*",
			want: false,
		},
		{
			name: "TEST3",
			s:    "abc/*",
			want: true,
		},
		{
			name: "TEST4",
			s:    "abc/abc",
			want: true,
		},
		{
			name: "TEST5",
			s:    "abc/abc*",
			want: false,
		},
		{
			name: "TEST6",
			s:    "abc/abc-123",
			want: true,
		},
		{
			name: "TEST6",
			s:    "*/abc-123",
			want: true,
		},
		{
			name: "TEST6",
			s:    "123/abc-123",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidResourceExpression(tt.s); got != tt.want {
				t.Errorf("isValidResourceExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitResourceExpression(t *testing.T) {
	tests := []struct {
		name          string
		expression    string
		wantNamespace string
		wantName      string
		wantErr       bool
	}{
		{
			name:          "TEST1",
			expression:    "a/b",
			wantNamespace: "a",
			wantName:      "b",
			wantErr:       false,
		},
		{
			name:          "TEST1",
			expression:    "a/b/c",
			wantNamespace: "a",
			wantName:      "b/c",
			wantErr:       false,
		},
		{
			name:          "TEST1",
			expression:    "ac",
			wantNamespace: "",
			wantName:      "",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotName, err := splitResourceExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitResourceExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNamespace != tt.wantNamespace {
				t.Errorf("splitResourceExpression() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotName != tt.wantName {
				t.Errorf("splitResourceExpression() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
