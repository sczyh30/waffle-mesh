package crd

import (
	"log"

	"github.com/sczyh30/waffle-mesh/brain/k8s"
	crdV1 "github.com/sczyh30/waffle-mesh/brain/k8s/crd/v1"
	clientset "github.com/sczyh30/waffle-mesh/brain/k8s/crd/gen/clientset/versioned"
	listers "github.com/sczyh30/waffle-mesh/brain/k8s/crd/gen/listers/crd/v1"
	informers "github.com/sczyh30/waffle-mesh/brain/k8s/crd/gen/informers/externalversions"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/util/runtime"
)

type RouteRuleController struct {
	client clientset.Interface

	informer    cache.SharedIndexInformer
	lister      listers.RouteRuleLister
	routeSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

func newRouteRuleClientset() (*clientset.Clientset, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return clientset.NewForConfig(restConfig)
}

func NewRouteRuleController(options k8s.ControllerOptions) (*RouteRuleController, error) {
	log.Printf("RouteRule CRD controller watching Kubernetes namespace %s\n", options.WatchedNamespace)
	client, err := newRouteRuleClientset()
	if err != nil {
		return nil, err
	}
	routeRuleInformer := informers.NewSharedInformerFactory(client, options.ResyncPeriod).Config().V1().RouteRules()
	controller := &RouteRuleController{
		informer:    routeRuleInformer.Informer(),
		lister:      routeRuleInformer.Lister(),
		client:      client,
		routeSynced: routeRuleInformer.Informer().HasSynced,
		queue:       workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	return controller, nil
}

func (c *RouteRuleController) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	<-stopCh
	log.Println("Waffle CRD controller terminated")
}

func (c *RouteRuleController) GetRouteRules() []*crdV1.RouteRule {
	var rules []*crdV1.RouteRule
	for _, v := range c.informer.GetStore().List() {
		rules = append(rules, v.(*crdV1.RouteRule))
	}
	return rules
}
