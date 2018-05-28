package cmd

import (
	"fmt"
	"os"
	"io"
	"bufio"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

type injectArgs struct {
	output              string
	ignoreInboundPorts  []uint
	ignoreOutboundPorts []uint
}

var injectCmdArgs = injectArgs{}

var injectCommand = &cobra.Command{
	Use:   "inject",
	Short: "Inject the Waffle Proxy as sidecar proxy to Kubernetes application deployment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) < 1 {
			return fmt.Errorf("please specify the deployment config file")
		}

		if in, err := os.Open(args[0]); err != nil {
			return err
		} else {
			reader := in
			defer func() {
				if errClose := in.Close(); errClose != nil {
					fmt.Errorf("error when closing file %s: %s", args[0], errClose)

					if err == nil {
						err = errClose
					}
				}
			}()
			var writer io.Writer
			if injectCmdArgs.output == "" {
				writer = os.Stdout
			} else {
				var out *os.File
				if out, err = os.Create(injectCmdArgs.output); err != nil {
					return err
				}
				writer = out
				defer func() {
					if errClose := out.Close(); errClose != nil {
						fmt.Errorf("error when closing file %s: %s", args[0], errClose)

						if err == nil {
							err = errClose
						}
					}
				}()
			}
			return doInject(reader, writer)
		}
	},
}

func doInject(in io.Reader, out io.Writer) error {
	reader := yamlDecoder.NewYAMLReader(bufio.NewReaderSize(in, 4096))
	for {
		raw, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var meta metaV1.TypeMeta
		if err := yaml.Unmarshal(raw, &meta); err != nil {
			return err
		}

		var originObj interface{}
		var podTemplateSpec *v1.PodTemplateSpec

		// Currently only support Deployment.
		switch meta.Kind {
		case "Deployment":
			var deployment v1beta1.Deployment
			err = yaml.Unmarshal(raw, &deployment)
			if err != nil {
				return err
			}
			originObj = &deployment
			podTemplateSpec = &deployment.Spec.Template
		}
		output := raw
		if podTemplateSpec != nil && injectPodTemplateSpec(podTemplateSpec) {
			output, err = yaml.Marshal(originObj)
			if err != nil {
				return err
			}
		}

		out.Write(output)
		out.Write([]byte("---\n"))
	}
	return nil
}

func injectPodTemplateSpec(t *v1.PodTemplateSpec) bool {
	if t.Spec.HostNetwork {
		return false
	}

	proxyUid := int64(2186)
	proxyInboundPort := int32(9081)
	proxyMetricsPort := int32(19802)
	// Waffle proxy sidecar container.
	sidecarProxyContainer := v1.Container{
		Name:            "waffle-proxy",
		Image:           "waffle.io/waffle-proxy:latest",
		ImagePullPolicy: v1.PullIfNotPresent,
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &proxyUid,
		},
		Ports: []v1.ContainerPort{
			{
				Name:          "proxy-inbound",
				ContainerPort: int32(proxyInboundPort),
			},
			{
				Name:          "proxy-metrics",
				ContainerPort: int32(proxyMetricsPort),
			},
		},
	}
	// Waffle proxy init container
	initContainer := v1.Container{
		Name:            "waffle-proxy-init",
		Image:           "waffle.io/waffle-proxy-init:latest",
		ImagePullPolicy: v1.PullIfNotPresent,
		SecurityContext: &v1.SecurityContext{
			Capabilities: &v1.Capabilities{
				Add: []v1.Capability{v1.Capability("NET_ADMIN")},
			},
		},
	}

	t.Spec.Containers = append(t.Spec.Containers, sidecarProxyContainer)
	t.Spec.InitContainers = append(t.Spec.InitContainers, initContainer)

	return true
}

func init() {
	injectCommand.PersistentFlags().StringVarP(&injectCmdArgs.output, "output", "o", "", "Inject output file")
	injectCommand.PersistentFlags().UintSliceVar(&injectCmdArgs.ignoreInboundPorts, "skip-inbound-ports", nil, "Ports that should skip the proxy and send directly to the application")
	injectCommand.PersistentFlags().UintSliceVar(&injectCmdArgs.ignoreOutboundPorts, "skip-outbound-ports", nil, "Outbound ports that should skip the proxy")
	WaffleCommand.AddCommand(injectCommand)
}
