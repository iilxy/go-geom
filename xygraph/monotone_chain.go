package xygraph

type MonotoneChain struct {
	mce        *MonotoneChainEdge
	chainIndex int
}

func (mc *MonotoneChain) computeIntersections(other MonotoneChain, si SegmentIntersector) {
	mc.mce.computeIntersectsForChain(mc.chainIndex, other.mce, other.chainIndex, si)
}
