package utils

import (
	"errors"
	"fmt"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

func GetYAML(kubeConfig, namespace, resourceType, resourceName string) (string, error) {
	cmd := exec.Command("kubectl", "--kubeconfig="+kubeConfig, "-n", namespace, "get", resourceType, resourceName, "-o", "yaml")

	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return "", err
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}

	if opBytes, err := io.ReadAll(stdout); err != nil {
		return "", err
	} else {
		yamlStr := string(opBytes)
		if strings.HasPrefix(yamlStr, "Error from server") {
			return "", errors.New(yamlStr)
		}
		return yamlStr, nil
	}
}

var ParseLineRegex = regexp.MustCompile("\\S+")

func ParseLine(line string) []string {
	return ParseLineRegex.FindAllString(line, -1)
}

// GetAllResourceNames get all resource from api server
// return map: namespace -> name -> struct{}
func GetAllResourceNames(kubeconfig string, resourceType string) (map[string]map[string]struct{}, error) {
	cmd := exec.Command("kubectl", "--kubeconfig="+kubeconfig, "get", resourceType, "-A")
	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	if opBytes, err := io.ReadAll(stdout); err != nil {
		return nil, err
	} else {
		cmdOutput := string(opBytes)
		lines := strings.Split(cmdOutput, "\n")
		lines = lines[1:]

		resourceMap := map[string]map[string]struct{}{}

		for _, line := range lines {
			if line == "" {
				continue
			}
			words := ParseLine(line)
			if len(words) < 2 {
				return nil, fmt.Errorf("GetAllResourceNames error: parse line error: %v", line)
			}
			namespace := words[0]
			name := words[1]
			var resourceSet map[string]struct{}
			if _, exist := resourceMap[namespace]; exist {
				resourceSet = resourceMap[namespace]
			} else {
				resourceSet = map[string]struct{}{}
				resourceMap[namespace] = resourceSet
			}
			resourceSet[name] = struct{}{}
		}

		return resourceMap, nil
	}
}

// GetAllNamespaces get all resource from api server
// return namespace list
func GetAllNamespaces(kubeconfig string) ([]string, error) {
	cmd := exec.Command("kubectl", "--kubeconfig="+kubeconfig, "get", "namespace", "-A")
	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	if opBytes, err := io.ReadAll(stdout); err != nil {
		return nil, err
	} else {
		cmdOutput := string(opBytes)
		if strings.HasPrefix(cmdOutput, "Error from server") {
			return nil, errors.New(cmdOutput)
		}
		lines := strings.Split(cmdOutput, "\n")
		lines = lines[1:]

		var namespaces []string
		for _, line := range lines {
			if line == "" {
				continue
			}
			words := ParseLine(line)
			if len(words) < 1 {
				return nil, fmt.Errorf("GetAllNamespaces error: parse line error: %v", line)
			}
			namespaces = append(namespaces, words[0])
		}

		return namespaces, nil
	}
}

// GetResources get resource from api server
// blank namespaces means all namespace
func GetResources(kubeconfig string, namespaces []string, resourceType string) ([]objects.StructuredObject, error) {
	if len(namespaces) == 0 {
		return GetResourcesFromAllNamespace(kubeconfig, resourceType)
	} else {
		var result []objects.StructuredObject
		for _, namespace := range namespaces {
			resources, err := GetResourcesFromNamespace(kubeconfig, namespace, resourceType)
			if err != nil {
				return nil, err
			}
			result = append(result, resources...)
		}
		return result, nil
	}
}

func GetResourcesFromAllNamespace(kubeconfig string, resourceType string) ([]objects.StructuredObject, error) {
	var cmdArguments []string
	kubeconfig = strings.TrimSpace(kubeconfig)
	if kubeconfig != "" {
		cmdArguments = append(cmdArguments, "--kubeconfig="+kubeconfig)
	}
	cmdArguments = append(cmdArguments, "get", resourceType, "-A", "-o", "yaml")

	cmd := exec.Command("kubectl", cmdArguments...)
	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	opBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	output := string(opBytes)
	if strings.HasPrefix(output, "Error from server") {
		return nil, errors.New(output)
	}

	o, err := objects.FromYAML(output)
	if err != nil {
		return nil, err
	}
	return o.GetObjects("items")
}

func GetResourcesFromNamespace(kubeconfig string, namespace string, resourceType string) ([]objects.StructuredObject, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is empty")
	}
	var cmdArguments []string
	kubeconfig = strings.TrimSpace(kubeconfig)
	if kubeconfig != "" {
		cmdArguments = append(cmdArguments, "--kubeconfig="+kubeconfig)
	}
	cmdArguments = append(cmdArguments, "-n", namespace, "get", resourceType, "-o", "yaml")

	cmd := exec.Command("kubectl", cmdArguments...)
	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	opBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	output := string(opBytes)
	if strings.HasPrefix(output, "Error from server") {
		return nil, errors.New(output)
	}

	o, err := objects.FromYAML(output)
	if err != nil {
		return nil, err
	}
	return o.GetObjects("items")
}
