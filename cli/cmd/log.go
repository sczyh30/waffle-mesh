package cmd

import (
	"fmt"
	"path"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

var kubeconfigPath = path.Join(homedir.HomeDir(), ".kube/config")

var logCommand = &cobra.Command {
	Use:   "log",
	Short: "Get logs of Waffle Brain",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		k8s, err := createInterface(kubeconfigPath)
		if err != nil {
			return err
		}
		podList, err := k8s.CoreV1().Pods("").List(v1.ListOptions{LabelSelector: "app=waffle-brain"})
		if err != nil {
			return err
		}
		for _, pod := range podList.Items {
			if strings.Contains(pod.Name, "waffle-brain") {
				c := exec.Command("kubectl", "log", pod.Name)
				output, err := c.Output()
				fmt.Printf("%s\n", string(output))
				return err
			}
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func createInterface(kubeconfig string) (kubernetes.Interface, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

func init() {
	WaffleCommand.AddCommand(logCommand)
}

