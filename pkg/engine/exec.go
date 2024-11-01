package engine

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

func Exec(ctx action.Context, yaml string, scripts string) (string, error) {
	obj, err := objects.FromYAML(yaml)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(scripts, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_action, err := ParseAction(line)
		if err != nil {
			fmt.Printf("parse line error: %v\n", line)
			continue
		}
		if _action != nil {
			_action.DoAction(ctx, obj)
		}
	}

	return objects.ToYAML(obj)
}
