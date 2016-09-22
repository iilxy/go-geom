package graph

import (
	"bytes"
	"github.com/twpayne/go-geom/xy/location"
)

type Label [2]topologyLocation

func NewLabel(geomIndex int, onLoc, leftLoc, rightLoc location.Type) *Label {
	label := NewNullLabel()
	label[geomIndex][OnLabel] = onLoc
	label[geomIndex][LeftOfLabel] = leftLoc
	label[geomIndex][RightOfLabel] = rightLoc
	return label
}
func NewLabelFromTemplate(template *Label) *Label {
	return &Label{
		newTopologyLocationFromTemplate(template[0]),
		newTopologyLocationFromTemplate(template[1]),
	}
}

func NewHomogeneousLabel(loc location.Type) *Label {
	return &Label{
		newOnTopologyLocation(loc),
		newOnTopologyLocation(loc),
	}
}
func NewNullLabel() *Label {
	return &Label{
		newOnTopologyLocation(location.None),
		newOnTopologyLocation(location.None),
	}
}

func (l *Label) flip() {
	l[0].flip()
	l[1].flip()
}
func (l *Label) setAllLocationsIfNull(loc location.Type) {
	l[0].setAllLocationsIfNull(loc)
	l[1].setAllLocationsIfNull(loc)
}

// Merge this label with another one.
// Merging updates any null attributes of this label with the attributes from lbl
func (l *Label) merge(lbl *Label) {
	for i := 0; i < 2; i++ {
		if l[i] == nil && lbl[i] != nil {
			l[i] = newOnTopologyLocation(lbl[i][OnLabel])
		} else {
			l[i].merge(lbl[i])
		}
	}
}

func (l *Label) getGeometryCount() int {
	count := 0
	if !l[0].isNull() {
		count++
	}
	if !l[1].isNull() {
		count++
	}
	return count
}

func (l Label) isArea() bool {
	return l[0].isArea() || l[1].isArea()
}

func (l Label) isEqualOnSide(lbl Label, side int) bool {
	return l[0].isEqualOnSide(lbl[0], side) && l[1].isEqualOnSide(lbl[1], side)
}

func (l *Label) toLine(geomIndex int) {
	if l[geomIndex].isArea() {
		l[geomIndex] = newOnTopologyLocation(l[geomIndex][0])
	}
}

func (l Label) toLineLabel() *Label {
	lineLabel := NewHomogeneousLabel(location.None)
	for i := 0; i < 2; i++ {
		lineLabel[i] = l[i]
	}

	return lineLabel
}

func (l Label) String() string {
	buf := bytes.Buffer{}
	if l[0] != nil {
		buf.WriteString("A:")
		buf.WriteString(l[0].String())
	}
	if l[1] != nil {
		buf.WriteString(" B:")
		buf.WriteString(l[1].String())
	}
	return buf.String()
}
