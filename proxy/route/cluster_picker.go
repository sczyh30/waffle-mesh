package route

import "github.com/sczyh30/waffle-mesh/api/gen"

type ClusterPicker interface {
	NextCluster() string
}

type ClusterWeightPair struct {
	name string
	weight uint32
	effectiveWeight uint32
	currentWeight uint32
}

type SmoothWeightedClusterPicker struct {
	weightedPairs []*ClusterWeightPair
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
	}
}

func (p *SmoothWeightedClusterPicker) NextCluster() string {
	// A smooth load balancing algorithm for weighted round-robin.
	var total uint32 = 0
	for _, pair := range p.weightedPairs {
		pair.currentWeight += pair.effectiveWeight
		total += pair.effectiveWeight
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