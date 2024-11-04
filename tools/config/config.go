package config

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/engine/match"
	"github.com/storm-blue/rubick/pkg/engine/scripts"
	"os"
	"regexp"
	"strings"
)

const (
	KubeconfigHead = "__kubeconfig__"
	ScriptsHead    = "__scripts__"
)

func ParseConfig(config string) (kubeconfig string, _scripts string, resourceMatchers map[string]match.Matcher, err error) {
	lines := strings.Split(config, "\n")

	var head string
	resources := map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")
		line = strings.Trim(line, "\r")

		if line == "" {
			continue
		}

		if scripts.IsComment(line) {
			continue
		}

		if isHead(line) {
			head = parseHead(line)
			if head != KubeconfigHead && head != ScriptsHead {
				if !isValidResourceTypeString(head) {
					return "", "", nil, fmt.Errorf("invalid resource type: %s", head)
				}
			}
			continue
		}

		if head == KubeconfigHead {
			if !isValidPath(line) {
				return "", "", nil, fmt.Errorf("invalid kubeconfig path: %s", kubeconfig)
			}
			kubeconfig = line
		} else if head == ScriptsHead {
			if _scripts == "" {
				_scripts = line
			} else {
				_scripts = _scripts + "\n" + line
			}
		} else {
			if !isValidResourceExpression(line) {
				return "", "", nil, fmt.Errorf("invalid resource expression: %s", line)
			}
			resources[head] = append(resources[head], line)
		}
	}

	if err = scripts.ValidateScripts(_scripts); err != nil {
		return "", "", nil, fmt.Errorf("validate scripts failed: %v", err)
	}

	if resourceMatchers, err = buildMatchers(resources); err != nil {
		return "", "", nil, fmt.Errorf("build matcher failed: %v", err)
	}

	return
}

func ParseConfigFile(fileName string) (kubeconfig string, scripts string, matchers map[string]match.Matcher, err error) {
	var bs []byte
	bs, err = os.ReadFile(fileName)
	if err != nil {
		return
	}

	content := string(bs)
	return ParseConfig(content)
}

func isHead(line string) bool {
	return strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")
}

func parseHead(line string) string {
	head := line[1 : len(line)-1]
	return strings.TrimSpace(head)
}

var resourceTypeRegex = regexp.MustCompile("^[a-zA-Z0-9.\\-]+$")

func isValidResourceTypeString(s string) bool {
	return resourceTypeRegex.MatchString(s)
}

var pathRegex = regexp.MustCompile("^[a-zA-Z0-9_\\-/.${}]+$")

func isValidPath(path string) bool {
	return pathRegex.MatchString(path)
}

var resourceNameRegex = regexp.MustCompile("^(\\*|([a-zA-Z0-9\\-])+)/(\\*|([a-zA-Z0-9\\-])+)$")

func isValidResourceExpression(s string) bool {
	return resourceNameRegex.MatchString(s)
}

func buildMatchers(resources map[string][]string) (map[string]match.Matcher, error) {
	result := make(map[string]match.Matcher)

	for resource, expressions := range resources {
		var singleLineMatchers []match.Matcher
		for _, expression := range expressions {
			namespace, name, err := splitResourceExpression(expression)
			if err != nil {
				return nil, err
			}
			namespaceMatcher := match.NewStringMatcher(namespace, "metadata.namespace")
			nameMatcher := match.NewStringMatcher(name, "metadata.name")
			singleLineMather := match.NewAndMatcher(namespaceMatcher, nameMatcher)

			singleLineMatchers = append(singleLineMatchers, singleLineMather)
		}
		resourceMatcher := match.NewOrMather(singleLineMatchers...)
		result[resource] = resourceMatcher
	}
	return result, nil
}

func splitResourceExpression(expression string) (namespace, name string, err error) {
	parts := strings.SplitN(expression, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid expression: %s", expression)
	}
	namespace = parts[0]
	name = parts[1]
	return
}
