package k8s

import (
	"time"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ControllerOptions struct {
	// Namespace the controller watches. If set to metaV1.NamespaceAll (""), controller watches all namespaces.
	WatchedNamespace string
	ResyncPeriod     time.Duration
}

type Controller struct {
	client kubernetes.Interface
	queue  workqueue.RateLimitingInterface

	services  cache.SharedIndexInformer
	endpoints cache.SharedIndexInformer
	nodes     cache.SharedIndexInformer
	pods      cache.SharedIndexInformer
}

func NewController(client kubernetes.Interface, options ControllerOptions) *Controller {
	log.Printf("Service controller watching Kubernetes namespace %s\n", options.WatchedNamespace)

	controller := &Controller{
		client: client,
		queue:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	controller.services = controller.createInformer(&v1.Service{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Services(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Services(options.WatchedNamespace).Watch(opts)
		})

	controller.endpoints = controller.createInformer(&v1.Endpoints{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).Watch(opts)
		})

	controller.nodes = controller.createInformer(&v1.Node{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Nodes().List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Nodes().Watch(opts)
		})

	controller.pods = controller.createInformer(&v1.Pod{}, options.ResyncPeriod,
		func(opts metaV1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).List(opts)
		},
		func(opts metaV1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).Watch(opts)
		})

	return controller
}

func (c *Controller) runQueue() {
	for c.handleQueueItem() {}
}

func (c *Controller) handleQueueItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)

	// TODO: Add custom handler.

	return true
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer c.queue.ShutDown()

	go c.services.Run(stopCh)
	go c.endpoints.Run(stopCh)
	go c.pods.Run(stopCh)
	go c.nodes.Run(stopCh)

	go wait.Until(c.runQueue, time.Second, stopCh)

	<-stopCh
	log.Println("Kubernetes Controller terminated")
}

func (c *Controller) createInformer(o runtime.Object, resyncPeriod time.Duration,
	lf cache.ListFunc, wf cache.WatchFunc) cache.SharedIndexInformer {

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{ListFunc: lf, WatchFunc: wf}, o,
		resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					c.queue.AddRateLimited(key)
				}
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					key, err := cache.MetaNamespaceKeyFunc(cur)
					if err == nil {
						c.queue.AddRateLimited(key)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					c.queue.AddRateLimited(key)
				}
			},
		})

	return informer
}
