package noding

import "github.com/twpayne/go-geom/transform"

type SegmentNodeList struct {
	// parent edge
	edge    NodeableSegmentString
	nodeMap *transform.TreeSet
}

func (nl SegmentNodeList) getNodeMap() {
	if nl.nodeMap == nil {
		nl.nodeMap = transform.NewTreeMap()
	}
}
