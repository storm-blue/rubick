package modifier

import (
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

func ModifyYAML(ctx action.Context, yamlStr string, actions []action.Action) (string, error) {
	object, err := objects.FromYAML(yamlStr)
	if err != nil {
		return "", err
	}

	for _, a := range actions {
		a.DoAction(ctx, object)
	}

	return object.ToYAML()
}

func ModifyObject(ctx action.Context, object objects.StructuredObject, actions []action.Action) {
	for _, a := range actions {
		a.DoAction(ctx, object)
	}
}
