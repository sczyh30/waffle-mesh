package k8s

import (
	"log"
	"reflect"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControllerOptions struct {
	// Namespace the controller watches. If set to metaV1.NamespaceAll (""), controller watches all namespaces.
	WatchedNamespace string
	ResyncPeriod     time.Duration
}

type WatcherEvent string

const (
	Add    WatcherEvent = "add"
	Update WatcherEvent = "update"
	Delete WatcherEvent = "delete"
)

type WatcherEventHandler func(obj interface{}, key string, event WatcherEvent) error

type WatchEventHandlerChain struct {
	eventHandlers []WatcherEventHandler
}

type EventTask struct {
	obj   interface{}
	key   string
	event WatcherEvent
	f     WatcherEventHandler
}

func (chain *WatchEventHandlerChain) AddHandler(handler WatcherEventHandler) {
	chain.eventHandlers = append(chain.eventHandlers, handler)
}

func (chain *WatchEventHandlerChain) Go(obj interface{}, key string, event WatcherEvent) error {
	for _, f := range chain.eventHandlers {
		if err := f(obj, key, event); err != nil {
			return err
		}
	}
	return nil
}

type WrappedIndexInformer struct {
	informer          cache.SharedIndexInformer
	eventHandlerChain *WatchEventHandlerChain
}

type Controller struct {
	client kubernetes.Interface
	queue  *RateLimitingWorkingQueue

	services  WrappedIndexInformer
	endpoints WrappedIndexInformer
	nodes     WrappedIndexInformer
	pods      *PodCache
}

func NewController(client kubernetes.Interface, options ControllerOptions) *Controller {
	log.Printf("Service controller watching Kubernetes namespace %s\n", options.WatchedNamespace)

	controller := &Controller{
		client: client,
		queue:  NewQueue(1 * time.Second),
	}

	controller.services = controller.createWrappedInformer(&v1.Service{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Services(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Services(options.WatchedNamespace).Watch(opts)
		})

	controller.endpoints = controller.createWrappedInformer(&v1.Endpoints{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).Watch(opts)
		})

	controller.nodes = controller.createWrappedInformer(&v1.Node{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Nodes().List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Nodes().Watch(opts)
		})

	controller.pods = createPodCache(controller.createWrappedInformer(&v1.Pod{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).Watch(opts)
		}))

	return controller
}

func (c *Controller) Run(stopCh chan struct{}) {
	go c.services.informer.Run(stopCh)
	go c.endpoints.informer.Run(stopCh)
	go c.pods.informer.Run(stopCh)
	go c.nodes.informer.Run(stopCh)

	go c.queue.Run(stopCh)

	<-stopCh
	log.Println("Waffle: Kubernetes controller terminated")
}

func (c *Controller) createWrappedInformer(o runtime.Object, resyncPeriod time.Duration,
	lf cache.ListFunc, wf cache.WatchFunc) WrappedIndexInformer {

	chain := &WatchEventHandlerChain{eventHandlers: []WatcherEventHandler{}}
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{ListFunc: lf, WatchFunc: wf}, o, resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					c.queue.Push(EventTask{obj: obj, key: key, event: Add, f: chain.Go})
				}
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					key, err := cache.MetaNamespaceKeyFunc(cur)
					if err == nil {
						c.queue.Push(EventTask{obj: cur, key: key, event: Update, f: chain.Go})
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					c.queue.Push(EventTask{obj: obj, key: key, event: Delete, f: chain.Go})
				}
			},
		})

	return WrappedIndexInformer{informer: informer, eventHandlerChain: chain}
}

func (c *Controller) GetEndpoints() []*v1.Endpoints {
	var endpoints []*v1.Endpoints
	list := c.endpoints.informer.GetStore().List()
	for _, v := range list {
		endpoints = append(endpoints, v.(*v1.Endpoints))
	}
	return endpoints
}

func (c *Controller) GetServices() []*v1.Service {
	var services []*v1.Service
	list := c.services.informer.GetStore().List()
	for _, v := range list {
		services = append(services, v.(*v1.Service))
	}
	return services
}

func (c *Controller) GetPodCache() *PodCache {
	return c.pods
}

func (c *Controller) GetServiceByName(serviceName string) (*v1.Service, bool) {
	svc, exists, err := c.services.informer.GetStore().GetByKey("default/" + serviceName)
	if err != nil {
		return nil, false
	}
	if !exists {
		return nil, false
	}
	return svc.(*v1.Service), true
}

func createPodCache(informer WrappedIndexInformer) *PodCache {
	podCache := &PodCache{keyByAddress: make(map[string]string), informer: informer.informer}
	informer.eventHandlerChain.AddHandler(func(obj interface{}, key string, event WatcherEvent) error {
		podCache.rwMutex.Lock()
		defer podCache.rwMutex.Unlock()

		pod := obj.(*v1.Pod)
		ipAddress := pod.Status.PodIP

		if len(ipAddress) > 0 {
			switch event {
			case Add:
				switch pod.Status.Phase {
				case v1.PodPending, v1.PodRunning:
					podCache.keyByAddress[ipAddress] = key
				}
			case Update:
				switch pod.Status.Phase {
				case v1.PodPending, v1.PodRunning:
					podCache.keyByAddress[ipAddress] = key
				default:
					if podCache.keyByAddress[ipAddress] == key {
						delete(podCache.keyByAddress, ipAddress)
					}
				}
			case Delete:
				if podCache.keyByAddress[ipAddress] == key {
					delete(podCache.keyByAddress, ipAddress)
				}
			}
		}
		return nil
	})
	return podCache
}
