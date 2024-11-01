package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/storm-blue/rubick/exporter"
	"github.com/storm-blue/rubick/pkg/common"
	"io/fs"
	"os"
	"strings"
)

type Schema string

var (
	configFile *string
	kubeconfig *string

	kubeconfigsP  *[]string
	namespacesP   *[]string
	applicationsP *[]string
	domainsP      *[]string

	rootCmd = &cobra.Command{
		Use:   "help",
		Short: "资源导出",
		Long:  `将部署应用所需要的资源从目标环境（k8s）中导出`,
	}

	allCmd = &cobra.Command{
		Use:   "all",
		Short: "资源导出-all",
		Long:  `将部署应用所需要的资源从目标环境（k8s）中导出，会导出所有deployment、service、ingress、secrets资源`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("导出参数：")
			fmt.Println(*kubeconfig)

			allResources, err := exporter.Export_All(*kubeconfig)
			if err != nil {
				return err
			}

			resultResources := map[string]string{}

			for resourceType, resources := range allResources {
				for i, resource := range resources {
					resultResources[fmt.Sprintf("%v_%v.yaml", resourceType, i)] = resource
				}
			}

			zipBytes, err := common.Zip(resultResources)
			if err != nil {
				return err
			}

			err = os.WriteFile("resource.zip", zipBytes, fs.ModePerm)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	//all
	kubeconfig = allCmd.Flags().String("kubeconfig", "", "导出的目标集群连接配置")
	err := allCmd.MarkFlagRequired("kubeconfig")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(allCmd)
}

func main() {
	_ = Execute()
}

// ParseConfigFile parse file like:
// [kubeconfig]
// /User/xyz/kubeconfig1
// /User/xyz/kubeconfig2
// [namespace]
// default
// ns1
// ns2
// [application]
// alivia2-root
// alivia2-root2
// alivia2-root3
// [domain]
// www.meetwhale.com
// xyz.meetwhale.com
func ParseConfigFile(fileName string) (kubeconfigs, namespaces, applications, domains []string, err error) {
	var bs []byte
	bs, err = os.ReadFile(fileName)
	if err != nil {
		return
	}

	var mode string

	const (
		K = "[kubeconfig]"
		N = "[namespace]"
		A = "[application]"
		D = "[domain]"
	)

	content := string(bs)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")
		line = strings.Trim(line, "\r")

		if line == "" {
			continue
		}

		switch line {
		case K:
			mode = line
			continue
		case N:
			mode = line
			continue
		case A:
			mode = line
			continue
		case D:
			mode = line
			continue
		default:
			switch mode {
			case K:
				kubeconfigs = append(kubeconfigs, line)
			case N:
				namespaces = append(namespaces, line)
			case A:
				applications = append(applications, line)
			case D:
				domains = append(domains, line)
			}
		}
	}
	return
}
