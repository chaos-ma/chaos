package random

import (
	"context"
	"math/rand"

	selector2 "chaos/server/rpcserver/selector"
	"chaos/server/rpcserver/selector/node/direct"
)

const (
	// Name is random balancer name
	Name = "random"
)

var _ selector2.Balancer = &Balancer{} // Name is balancer name

// Balancer is a random balancer.
type Balancer struct{}

// New an random selector.
func New() selector2.Selector {
	return NewBuilder().Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector2.WeightedNode) (selector2.WeightedNode, selector2.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector2.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with random balancer
func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is random builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector2.Balancer {
	return &Balancer{}
}
