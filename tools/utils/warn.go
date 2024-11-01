package utils

import (
	"github.com/storm-blue/rubick/pkg/engine/match"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

var WarningNamespaces = map[string]struct{}{
	"istio-system":    {},
	"kube-system":     {},
	"kube-public":     {},
	"kube-ops":        {},
	"kube-node-lease": {},
	"cattle-system":   {},
	"ingress-nginx":   {},
	"ack-onepilot":    {},
	"ack-csi-fuse":    {},
}

var WarningMather match.Matcher

func init() {
	var matchers []match.Matcher
	for namespace := range WarningNamespaces {
		matchers = append(matchers, match.NewStringMatcher(namespace, "metadata.namespace"))
	}
	WarningMather = match.NewOrMather(matchers...)
}

func IsWarningResource(resource objects.StructuredObject) bool {
	return WarningMather.Match(resource)
}
