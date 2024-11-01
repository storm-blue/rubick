package constants

var ExcludeConfigmaps = map[string]struct{}{
	"kube-root-ca.crt": {},
}

var ExcludeNamespaces = map[string]struct{}{
	"istio-system":                {},
	"kube-system":                 {},
	"kube-public":                 {},
	"kube-node-lease":             {},
	"cattle-dashboards":           {},
	"cattle-fleet-system":         {},
	"cattle-impersonation-system": {},
	"cattle-monitoring-system":    {},
	"cattle-system":               {},
	"cattle-logging-system":       {},
	"ingress-nginx":               {},
	"local":                       {},
}

const (
	Namespace   = "namespace"
	Deployment  = "deployment"
	Service     = "service"
	Ingress     = "ingress"
	Secret      = "secret"
	Configmap   = "configmap"
	StatefulSet = "statefulset"
	Cronjob     = "cronjob"
)
