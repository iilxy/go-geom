package xygraph

import (
	"bytes"
	"github.com/twpayne/go-geom/xy/location"
)

type TopologyLocation []location.Type

func NewTopologyLocationFromTemplate(template TopologyLocation) TopologyLocation {
	topoLoc := make(TopologyLocation, len(template))
	copy(topoLoc, template)
	return topoLoc
}

func NewOnTopologyLocation(on location.Type) TopologyLocation {
	return TopologyLocation{on}
}

func NewTopologyLocation(on, left, right location.Type) TopologyLocation {
	return TopologyLocation{on, left, right}
}

func (topoLoc TopologyLocation) isNull() bool {
	for _, l := range topoLoc {
		if l != location.None {
			return false
		}
	}

	return true
}

func (topoLoc TopologyLocation) isEqualOnSide(le TopologyLocation, locIndex int) bool {
	return topoLoc[locIndex] == le[locIndex]
}

func (topoLoc TopologyLocation) isArea() bool {
	return len(topoLoc) > 1
}

func (topoLoc TopologyLocation) isLine() bool {
	return len(topoLoc) == 1
}

func (topoLoc *TopologyLocation) flip() {
	if len(topoLoc) <= 1 {
		return
	}
	temp := topoLoc[LEFT]
	topoLoc[LEFT] = topoLoc[RIGHT]
	topoLoc[RIGHT] = temp
}

func (topoLoc *TopologyLocation) setAllLocations(locValue location.Type) {
	for i := 0; i < len(topoLoc); i++ {
		topoLoc[i] = locValue
	}
}

func (topoLoc *TopologyLocation) setAllLocationsIfNull(locValue location.Type) {
	for i := 0; i < len(topoLoc); i++ {
		if topoLoc[i] == location.None {
			topoLoc[i] = locValue
		}
	}
}
func (topoLoc *TopologyLocation) allPositionsEqual(loc location.Type) bool {
	for _, l := range topoLoc {
		if l != loc {
			return false
		}
	}
	return true
}
func (topoLoc *TopologyLocation) merge(gl TopologyLocation) TopologyLocation {
	// if the src is an Area label & and the dest is not, increase the dest to be an Area
	if len(gl) > len(topoLoc) {
		newLoc := make([]int, 3)
		newLoc[ON] = topoLoc[ON]
		newLoc[LEFT] = location.None
		newLoc[RIGHT] = location.None
		topoLoc = newLoc
	}

	for i := 0; i < len(topoLoc); i++ {
		if topoLoc[i] == location.None && i < len(gl) {
			topoLoc[i] = gl[i]
		}
	}
	return topoLoc
}

func (topoLoc TopologyLocation) String() string {
	buffer := bytes.Buffer{}
	if len(topoLoc) > 1 {
		buffer.WriteRune(topoLoc[LEFT].Symbol())
	}
	buffer.WriteRune(topoLoc[ON].Symbol())
	if len(topoLoc) > 1 {
		buffer.WriteRune(topoLoc[RIGHT].Symbol())
	}
	return buffer.String()

}
