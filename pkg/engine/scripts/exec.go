package scripts

import (
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"strings"
)

func ExecYAMLs(ctx action.Context, multiYaml string, scripts string) (string, error) {
	_objects, err := objects.FromYAMLs(multiYaml)
	if err != nil {
		return "", err
	}

	if __objects, err := ExecObjects(ctx, _objects, scripts); err != nil {
		return "", err
	} else {
		return objects.ToYAMLs(__objects)
	}
}

func ExecYAML(ctx action.Context, yaml string, scripts string) (string, error) {
	obj, err := objects.FromYAML(yaml)
	if err != nil {
		return "", err
	}

	if err = ExecObject(ctx, obj, scripts); err != nil {
		return "", err
	}

	return objects.ToYAML(obj)
}

func ExecObject(ctx action.Context, object objects.StructuredObject, scripts string) error {
	actions, err := ParseScripts(scripts)
	if err != nil {
		return err
	}

	for _, _action := range actions {
		_action.DoAction(ctx, object)
	}

	return nil
}

func ExecObjects(ctx action.Context, _objects []objects.StructuredObject, scripts string) ([]objects.StructuredObject, error) {
	actions, err := ParseScripts(scripts)
	if err != nil {
		return _objects, err
	}

	var result []objects.StructuredObject

	for _, _object := range _objects {
		for _, _action := range actions {
			_action.DoAction(ctx, _object)
		}
		if !_object.Metadata().Removed() {
			result = append(result, _object)
		}
	}

	return result, nil
}

func ValidateScripts(scripts string) error {
	if _, err := ParseScripts(scripts); err != nil {
		return err
	}
	return nil
}

func ParseScripts(scripts string) ([]action.Action, error) {
	var actions []action.Action
	for _, line := range strings.Split(scripts, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if IsComment(line) {
			continue
		}

		_action, err := ParseAction(line)
		if err != nil {
			return nil, fmt.Errorf("parse scripts line error: \nscripts = [ %s ] \nerr = %v", line, err)
		}
		actions = append(actions, _action)
	}
	return actions, nil
}

func IsComment(line string) bool {
	return strings.HasPrefix(line, "#")
}
