package xygraph

type monotoneChain struct {
	mce        *monotoneChainEdge
	chainIndex int
}

func (mc *monotoneChain) computeIntersections(other monotoneChain, si SegmentIntersector) {
	mc.mce.computeIntersectsForChain(mc.chainIndex, other.mce, other.chainIndex, si)
}
