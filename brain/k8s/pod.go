package k8s

import (
	"sync"

	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
)

type PodCache struct {
	rwMutex sync.RWMutex

	informer     cache.SharedIndexInformer
	keyByAddress map[string]string
}

func (pc *PodCache) GetPodKeyByIp(address string) (string, bool) {
	pc.rwMutex.RLock()
	defer pc.rwMutex.RUnlock()

	key, exists := pc.keyByAddress[address]
	return key, exists
}

func (pc *PodCache) GetPodByIP(address string) (*v1.Pod, bool) {
	pc.rwMutex.RLock()
	defer pc.rwMutex.RUnlock()

	key, exists := pc.keyByAddress[address]
	if !exists {
		return nil, false
	}
	item, exists, err := pc.informer.GetStore().GetByKey(key)
	if !exists || err != nil {
		return nil, false
	}
	return item.(*v1.Pod), true
}

func (pc *PodCache) LabelsByIP(address string) (map[string]string, bool) {
	pod, exists := pc.GetPodByIP(address)
	if !exists {
		return nil, false
	}
	return ConvertLabels(pod.ObjectMeta), true
}
