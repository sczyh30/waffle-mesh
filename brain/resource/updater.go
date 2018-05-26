package resource

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sczyh30/waffle-mesh/brain/k8s"
	"github.com/sczyh30/waffle-mesh/brain/k8s/crd"
)

type XdsResourceUpdater struct {
	k8sController *k8s.Controller
	crdController *crd.RouteRuleController

	converter *ResourceConverter
}

func NewXdsResourceUpdater(kc *k8s.Controller, cc *crd.RouteRuleController) *XdsResourceUpdater {
	return &XdsResourceUpdater{
		converter:     NewResourceConverter(kc, cc),
		k8sController: kc,
		crdController: cc,
	}
}

func (updater *XdsResourceUpdater) fetchAndUpdate() error {
	routeRules := updater.crdController.GetRouteRules()
	routeConfigs, selectors := updater.converter.AggregateProxyRouteConfigs(routeRules)
	clusters := updater.converter.BuildOutboundClusters(routeConfigs, selectors)
	endpoints := updater.converter.BuildClusterEndpoints(selectors)

	XdsCache.updateCache(clusters, endpoints, routeConfigs)

	data, _ := json.Marshal(endpoints)
	log.Println("Xds endpoints updated: " + string(data))
	data, _ = json.Marshal(routeConfigs)
	log.Println("Xds route configs updated: " + string(data))

	return nil
}

func (updater *XdsResourceUpdater) Start(stop chan struct{}) error {
	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-stop:
			log.Println("Stopping the xDS resource updater")
			break
		case <-ticker.C:
			go updater.fetchAndUpdate()
		}
	}

	return nil
}
