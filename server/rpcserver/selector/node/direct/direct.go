package direct

import (
	"context"
	"sync/atomic"
	"time"

	selector2 "github.com/chaos-ma/chaos/server/rpcserver/selector"
)

const (
	defaultWeight = 100
)

var (
	_ selector2.WeightedNode        = &Node{}
	_ selector2.WeightedNodeBuilder = &Builder{}
)

// Node is endpoint instance
type Node struct {
	selector2.Node

	// last lastPick timestamp
	lastPick int64
}

// Builder is direct node builder
type Builder struct{}

// Build create node
func (*Builder) Build(n selector2.Node) selector2.WeightedNode {
	return &Node{Node: n, lastPick: 0}
}

func (n *Node) Pick() selector2.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di selector2.DoneInfo) {}
}

// Weight is node effective weight
func (n *Node) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return defaultWeight
}

func (n *Node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

func (n *Node) Raw() selector2.Node {
	return n.Node
}
