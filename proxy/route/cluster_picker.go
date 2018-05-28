package route

import (
	"sync"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type ClusterPicker interface {
	NextCluster() string
}

type ClusterWeightPair struct {
	name string
	weight uint32
	effectiveWeight uint32
	currentWeight int
}

type SmoothWeightedClusterPicker struct {
	weightedPairs []*ClusterWeightPair

	mutex sync.Mutex
}

func NewSmoothWeightedClusterPicker(wc *api.WeightedCluster) *SmoothWeightedClusterPicker {
	wp := make([]*ClusterWeightPair, len(wc.Clusters))
	for i, pair := range wc.Clusters {
		wp[i] = &ClusterWeightPair{
			name: pair.Name,
			weight: pair.Weight,
			effectiveWeight: pair.Weight,
			currentWeight: 0,
		}
	}
	return &SmoothWeightedClusterPicker{
		weightedPairs: wp,
		mutex: sync.Mutex{},
	}
}

func (p *SmoothWeightedClusterPicker) NextCluster() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// A smooth load balancing algorithm for weighted round-robin.
	var total = 0
	for _, pair := range p.weightedPairs {
		pair.currentWeight += int(pair.effectiveWeight)
		total += int(pair.effectiveWeight)
	}
	max := p.weightedPairs[0]
	for _, pair := range p.weightedPairs {
		if pair.currentWeight > max.currentWeight {
			max = pair
		}
	}

	max.currentWeight -= total
	return max.name
}

type SingleClusterPicker struct {
	Name string
}

func (p *SingleClusterPicker) NextCluster() string {
	return p.Name
}