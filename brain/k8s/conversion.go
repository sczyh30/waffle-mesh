package k8s

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertLabels(obj metaV1.ObjectMeta) map[string]string {
	out := make(map[string]string, len(obj.Labels))
	for k, v := range obj.Labels {
		out[k] = v
	}
	return out
}

func ConvertKey(name, namespace string) string {
	if len(namespace) == 0 {
		return name
	}
	return namespace + "/" + name
}
