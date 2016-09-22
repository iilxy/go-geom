package noding

type SegmentNodeComparer struct{}

func (comp SegmentNodeComparer) IsEquals(o1, o2 interface{}) bool {
	n1 := o1.(SegmentNode)
	n2 := o1.(SegmentNode)

}
func (comp SegmentNodeComparer) IsLess(o1, o2 interface{}) bool {

}
