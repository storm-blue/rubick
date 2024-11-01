package exporter

import (
	"gitlab.meetwhale.com/rubick/pkg/constants"
	"gitlab.meetwhale.com/rubick/pkg/modifier/objects"
)

func AllGetNamespaces(kubeconfig string) ([]string, error) {
	allNamespacesYaml, err := GetAllResources(kubeconfig, "namespace")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allNamespacesYaml)
	if err != nil {
		return nil, err
	}

	namespaceObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, namespaceObject := range namespaceObjects {
		namespace, err := namespaceObject.GetString("metadata.name")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		namespaceYAML, err := namespaceObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, namespaceYAML)
	}

	return result, nil
}

func AllGetDeployments(kubeconfig string) ([]string, error) {
	allDeploymentsYaml, err := GetAllResources(kubeconfig, "deployment")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allDeploymentsYaml)
	if err != nil {
		return nil, err
	}

	deploymentObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, deploymentObject := range deploymentObjects {
		namespace, err := deploymentObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		deploymentYAML, err := deploymentObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, deploymentYAML)
	}

	return result, nil
}

func AllGetServices(kubeconfig string) ([]string, error) {
	allServicesYaml, err := GetAllResources(kubeconfig, "service")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allServicesYaml)
	if err != nil {
		return nil, err
	}

	serviceObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, serviceObject := range serviceObjects {
		namespace, err := serviceObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		serviceYAML, err := serviceObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, serviceYAML)
	}

	return result, nil
}

func AllGetIngresses(kubeconfig string) ([]string, error) {
	allIngressesYaml, err := GetAllResources(kubeconfig, "ingress")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allIngressesYaml)
	if err != nil {
		return nil, err
	}

	ingressObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, ingressObject := range ingressObjects {
		namespace, err := ingressObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		ingressYAML, err := ingressObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, ingressYAML)
	}

	return result, nil
}

func AllGetSecrets(kubeconfig string) ([]string, error) {
	allSecretsYaml, err := GetAllResources(kubeconfig, "secret")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allSecretsYaml)
	if err != nil {
		return nil, err
	}

	secretObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, secretObject := range secretObjects {
		namespace, err := secretObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		typ, err := secretObject.GetString("type")
		if err != nil {
			return nil, err
		}

		if typ != "kubernetes.io/dockerconfigjson" {
			continue
		}

		secretYAML, err := secretObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, secretYAML)
	}

	return result, nil
}

func AllGetConfigmaps(kubeconfig string) ([]string, error) {
	allConfigmapYaml, err := GetAllResources(kubeconfig, "configmap")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allConfigmapYaml)
	if err != nil {
		return nil, err
	}

	configmapObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, configmapObject := range configmapObjects {
		namespace, err := configmapObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		name, err := configmapObject.GetString("metadata.name")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeConfigmaps[name]; exist {
			continue
		}

		configmapYAML, err := configmapObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, configmapYAML)
	}

	return result, nil
}

func AllGetStatefulSets(kubeconfig string) ([]string, error) {
	allStatefulSetYaml, err := GetAllResources(kubeconfig, "statefulset")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allStatefulSetYaml)
	if err != nil {
		return nil, err
	}

	StatefulSetObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, StatefulSetObject := range StatefulSetObjects {
		namespace, err := StatefulSetObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		statefulSetYAML, err := StatefulSetObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, statefulSetYAML)
	}

	return result, nil
}

func AllCronJobs(kubeconfig string) ([]string, error) {
	allCronJobsYaml, err := GetAllResources(kubeconfig, "cronjob")
	if err != nil {
		return nil, err
	}
	object, err := objects.FromYAML(allCronJobsYaml)
	if err != nil {
		return nil, err
	}

	cronJobObjects, err := object.GetObjects("items[*]")
	if err != nil {
		return nil, err
	}

	var result []string

	for _, cronJobObject := range cronJobObjects {
		namespace, err := cronJobObject.GetString("metadata.namespace")
		if err != nil {
			return nil, err
		}

		if _, exist := constants.ExcludeNamespaces[namespace]; exist {
			continue
		}

		cronJobYAML, err := cronJobObject.ToYAML()
		if err != nil {
			return nil, err
		}

		result = append(result, cronJobYAML)
	}

	return result, nil
}

// Export_All export all resources
// return map: resource type -> resources
func Export_All(kubeconfig string) (map[string][]string, error) {
	resources := map[string][]string{}

	allNamespaceYAMLs, err := AllGetNamespaces(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Namespace] = allNamespaceYAMLs

	allDeploymentYAMLs, err := AllGetDeployments(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Deployment] = allDeploymentYAMLs

	allServiceYAMLs, err := AllGetServices(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Service] = allServiceYAMLs

	allIngressYAMLs, err := AllGetIngresses(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Ingress] = allIngressYAMLs

	allSecretYAMLs, err := AllGetSecrets(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Secret] = allSecretYAMLs

	allConfigmapYAMLs, err := AllGetConfigmaps(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Configmap] = allConfigmapYAMLs

	allStatefulSetYAMLs, err := AllGetStatefulSets(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.StatefulSet] = allStatefulSetYAMLs

	allCronJobYAMLs, err := AllCronJobs(kubeconfig)
	if err != nil {
		return nil, err
	}
	resources[constants.Cronjob] = allCronJobYAMLs

	return resources, nil
}
