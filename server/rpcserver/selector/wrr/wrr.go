package wrr

import (
	"context"
	"sync"

	selector2 "github.com/chaos-ma/chaos/server/rpcserver/selector"
	"github.com/chaos-ma/chaos/server/rpcserver/selector/node/direct"
)

const (
	// Name is wrr balancer name
	Name = "wrr"
)

var _ selector2.Balancer = &Balancer{} // Name is balancer name

// Option is random builder option.

// Balancer is a random balancer.
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
}

// New random a selector.
func New() selector2.Selector {
	return NewBuilder().Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector2.WeightedNode) (selector2.WeightedNode, selector2.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector2.ErrNoAvailable
	}
	var totalWeight float64
	var selected selector2.WeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is wrr builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector2.Balancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
