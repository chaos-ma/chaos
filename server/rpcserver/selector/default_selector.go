package selector

import (
	"context"
	"sync/atomic"
)

// Default is composite selector.
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer

	nodes atomic.Value
}

// Select is select one node.
func (d *Default) Select(ctx context.Context) (selected Node, done DoneFunc, err error) {
	var (
		candidates []WeightedNode
	)
	nodes, ok := d.nodes.Load().([]WeightedNode)
	if !ok {
		return nil, nil, ErrNoAvailable
	}
	candidates = nodes

	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	p, ok := FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// Apply update nodes info.
func (d *Default) Apply(nodes []Node) {
	weightedNodes := make([]WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.NodeBuilder.Build(n))
	}
	// TODO: Do not delete unchanged nodes
	d.nodes.Store(weightedNodes)
}

// DefaultBuilder is de
type DefaultBuilder struct {
	Node     WeightedNodeBuilder
	Balancer BalancerBuilder
}

// Build create builder
func (db *DefaultBuilder) Build() Selector {
	return &Default{
		NodeBuilder: db.Node,
		Balancer:    db.Balancer.Build(),
	}
}
