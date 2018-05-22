package v1

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RouteRule struct {
	metaV1.TypeMeta    `json:",inline"`
	metaV1.ObjectMeta  `json:"metadata"`
	Spec RouteRuleSpec `json:"spec"`
}

type RouteRuleSpec struct {
	Destination RouteRuleDestination  `json:"destination"`
	Order       int                   `json:"order, omitempty"`
	Route       []RouteSelectorWeight `json:"route"`
	Match       RouteMatchCondition   `json:"match, omitempty"`
}

type RouteSelectorWeight struct {
	Labels map[string]string `json:"labels"`
	Weight uint32               `json:"weight"`
}

type RouteMatchCondition struct {
	Request RequestMatchCondition `json:"request"`
}

type RequestMatchCondition struct {
	Headers map[string]StringMatchCondition `json:"headers"`
}

type StringMatchCondition struct {
	Exact  string `json:"exact, omitempty"`
	Prefix string `json:"prefix, omitempty"`
	Regex  string `json:"regex, omitempty"`
}

type RouteRuleDestination struct {
	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RouteRuleList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata"`

	Items []RouteRule `json:"items"`
}
