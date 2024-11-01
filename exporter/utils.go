package exporter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

func SearchResources(kubeConfig string, resourceType string, resourceName string, namespaces []string) (string, string, error) {
	if len(namespaces) == 0 {
		cmd := exec.Command("kubectl", "--kubeconfig="+kubeConfig, "get", resourceType, "-A")
		stdout, err := cmd.StdoutPipe()
		defer func() { _ = stdout.Close() }()

		if err != nil {
			return "", "", err
		}

		if err = cmd.Start(); err != nil {
			return "", "", err
		}

		if opBytes, err := ioutil.ReadAll(stdout); err != nil {
			return "", "", err
		} else {
			return string(opBytes), "", nil
		}
	} else {
		return "", "", nil
	}
}

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

	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
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

	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
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

	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
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

// GetAllResources get all resource from api server
func GetAllResources(kubeconfig string, resourceType string) (string, error) {
	cmd := exec.Command("kubectl", "--kubeconfig="+kubeconfig, "get", resourceType, "-A", "-o", "yaml")
	stdout, err := cmd.StdoutPipe()
	defer func() { _ = stdout.Close() }()

	if err != nil {
		return "", err
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}

	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		return "", err
	} else {
		return string(opBytes), nil
	}
}
