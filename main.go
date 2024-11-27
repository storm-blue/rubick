package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/storm-blue/rubick/pkg/config"
	"github.com/storm-blue/rubick/pkg/engine/scripts"
	"github.com/storm-blue/rubick/pkg/modifier/action"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
	"github.com/storm-blue/rubick/pkg/utils"
	"os"
	"time"
)

const (
	TimeFormat = "20060102-150405"
)

var (
	kubeconfig       *string
	namespaces       *[]string
	resource         *string
	yamlFile         *string
	scriptsFile      *string
	exportOutputFile *string
	modifyOutputFile *string
	execOutputFile   *string
	configFile       *string

	rootCmd = &cobra.Command{
		Use:   "help",
		Short: "Rubick",
		Long: `Rubick是一个可以将k8s资源导出到文件的工具，也是一个可以通过自定义规则处理任意YAML的工具。它
可以用于需要批量处理k8s资源的场景：比如需要将一个集群的所有服务复制到另一个集群，但需要批量修
改服务的标签和部分属性以满足迁移需要；再比如，需要将一个命名空间的内容复制到另一个命名空间，并
且需要添加额外的注释。如果通过代码批量处理这些场景比较繁琐且无法复用，此时便可以使用Rubick来
进行相应资源的导出和清理工作。
Rubick is a tool that can export k8s resources to files, and it is also a tool
that can process arbitrary YAML through custom rules. Rubick can be used in
scenarios that require batch processing of k8s resources: for example, you need
to copy all services of one cluster to another cluster, but you need to modify
the labels and some attributes of the services in batches to meet migration
needs; for example, you need to copy the contents of a namespace  to another
namespace and additional comments need to be added. If batch processing of
these scenarios through code is cumbersome and cannot be reused, you can use
Rubick to export and clean up the corresponding resources.
`,
	}

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "导出k8s资源",
		Long:  `将指定的资源从目标k8s集群中导出到当前目录`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("processing...")

			resources, err := utils.GetResources(*kubeconfig, *namespaces, *resource)
			if err != nil {
				return err
			}

			yamls, err := objects.ToYAMLs(resources)
			if err != nil {
				return err
			}

			outputFileName := *exportOutputFile
			if outputFileName == "" {
				outputFileName = fmt.Sprintf("%vs-exported-%v.yaml", *resource, time.Now().Format(TimeFormat))
			}

			if len(yamls) == 0 {
				fmt.Println("no results return after process, skip output results to file!")
			} else {
				if err := os.WriteFile(outputFileName, []byte(yamls), os.ModePerm); err != nil {
					return err
				}
			}

			fmt.Println("success.")
			return nil
		},
	}

	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "修改YAML文件",
		Long:  `通过自定义清洗规则脚本，对指定的YAML文件进行修改`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("processing...")

			yamlBytes, err := os.ReadFile(*yamlFile)
			if err != nil {
				return err
			}
			scriptsBytes, err := os.ReadFile(*scriptsFile)
			if err != nil {
				return err
			}

			yaml, err := scripts.ExecYAMLs(action.NewContext(nil), string(yamlBytes), string(scriptsBytes))
			if err != nil {
				return err
			}

			outputFileName := *modifyOutputFile
			if outputFileName == "" {
				outputFileName = fmt.Sprintf("modified-%v.yaml", time.Now().Format(TimeFormat))
			}

			if len(yaml) == 0 {
				fmt.Println("no results return after process, skip output results to file!")
			} else {
				if err := os.WriteFile(outputFileName, []byte(yaml), os.ModePerm); err != nil {
					return err
				}
			}

			fmt.Println("success.")
			return nil
		},
	}

	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "根据配置处理资源并导出",
		Long: `通过自定义配置文件，选择想要导出的资源，并且执行配置中的脚本进行清洗。
config example:

[__kubeconfig__]
${HOME}/.kube/config

[deployment]
*/redis

[service]
java-dev/*
java-qa1/*
java-qa2/*
java-sit/*

[__scripts__]
# common scripts
DELETE(metadata.annotations.(kubectl.kubernetes.io/last-applied-configuration))
DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)
DELETE(metadata.uid)
DELETE(status)

# service scripts
IF VALUE_OF(kind)=="Service" THEN DELETE(spec.clusterIP)
IF VALUE_OF(kind)=="Service" THEN DELETE(spec.clusterIPs)
IF VALUE_OF(kind)=="Service" THEN SET(spec.ports[port=8080].port, 80)
IF (VALUE_OF(kind)=="Service" && EXISTS(metadata.labels.(github.io/app))) THEN SET_WITH_VALUE_OF(metadata.name, metadata.labels.(github.io/app))
IF NOT_EXISTS(metadata.labels.(github.io/app)) THEN SET_WITH_VALUE_OF(metadata.labels.(github.io/app), metadata.name)
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("processing...")

			configBytes, err := os.ReadFile(*configFile)
			if err != nil {
				return err
			}

			outputFileName := *execOutputFile
			if outputFileName == "" {
				outputFileName = fmt.Sprintf("executed-%v.yaml", time.Now().Format(TimeFormat))
			}

			kubeconfig, _scripts, resources, err := config.ParseConfig(string(configBytes))
			if err != nil {
				return err
			}

			var _objects []objects.StructuredObject

			for resource, matcher := range resources {
				__objects, err := utils.GetResources(kubeconfig, nil, resource)
				if err != nil {
					return err
				}

				for _, __object := range __objects {
					if matcher.Match(__object) {
						_objects = append(_objects, __object)
					}
				}
			}

			__objects, err := scripts.ExecObjects(action.NewContext(nil), _objects, _scripts)
			if err != nil {
				return err
			}

			yaml, err := objects.ToYAMLs(__objects)
			if err != nil {
				return err
			}

			if len(yaml) == 0 {
				fmt.Println("no results return after process, skip output results to file!")
			} else {
				if err := os.WriteFile(outputFileName, []byte(yaml), os.ModePerm); err != nil {
					return err
				}
			}

			fmt.Println("success.")
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// exporter
	kubeconfig = exportCmd.Flags().String("kubeconfig", "", "导出的目标集群连接配置, 默认值为: ${HOME}/.kube/config")
	namespaces = exportCmd.Flags().StringArrayP("namespace", "n", nil, "要导出资源的命名空间, 默认所有命名空间")
	resource = exportCmd.Flags().StringP("resource", "r", "", "要导出的资源类型")
	err := exportCmd.MarkFlagRequired("resource")
	if err != nil {
		panic(err)
	}
	exportOutputFile = exportCmd.Flags().StringP("output", "o", "", "指定输出的文件路径")
	rootCmd.AddCommand(exportCmd)

	// modify
	yamlFile = modifyCmd.Flags().StringP("file", "f", "", "要修改的文件路径")
	err = modifyCmd.MarkFlagRequired("file")
	if err != nil {
		panic(err)
	}
	scriptsFile = modifyCmd.Flags().StringP("scripts", "s", "", "脚本文件路径")
	err = modifyCmd.MarkFlagRequired("scripts")
	if err != nil {
		panic(err)
	}
	modifyOutputFile = modifyCmd.Flags().StringP("output", "o", "", "指定输出的文件路径")
	rootCmd.AddCommand(modifyCmd)

	// exec
	configFile = execCmd.Flags().String("config", "", "配置文件路径")
	err = execCmd.MarkFlagRequired("config")
	if err != nil {
		panic(err)
	}
	execOutputFile = execCmd.Flags().StringP("output", "o", "", "指定输出的文件路径")
	rootCmd.AddCommand(execCmd)
}

func main() {
	_ = Execute()
}
