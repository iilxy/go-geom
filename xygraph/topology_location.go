package xygraph

import (
	"bytes"
	"github.com/twpayne/go-geom/xy/location"
)

type topologyLocation []location.Type

func newTopologyLocationFromTemplate(template topologyLocation) topologyLocation {
	topoLoc := make(topologyLocation, len(template))
	copy(topoLoc, template)
	return topoLoc
}

func newOnTopologyLocation(on location.Type) topologyLocation {
	return topologyLocation{on}
}

func newTopologyLocation(on, left, right location.Type) topologyLocation {
	return topologyLocation{on, left, right}
}

func (topoLoc topologyLocation) isNull() bool {
	for _, l := range topoLoc {
		if l != location.None {
			return false
		}
	}

	return true
}

func (topoLoc topologyLocation) isEqualOnSide(le topologyLocation, locIndex int) bool {
	return topoLoc[locIndex] == le[locIndex]
}

func (topoLoc topologyLocation) isArea() bool {
	return len(topoLoc) > 1
}

func (topoLoc topologyLocation) isLine() bool {
	return len(topoLoc) == 1
}

func (topoLoc topologyLocation) flip() {
	if len(topoLoc) <= 1 {
		return
	}
	temp := topoLoc[LEFT]
	topoLoc[LEFT] = topoLoc[RIGHT]
	topoLoc[RIGHT] = temp
}

func (topoLoc topologyLocation) setAllLocations(locValue location.Type) {
	for i := 0; i < len(topoLoc); i++ {
		topoLoc[i] = locValue
	}
}

func (topoLoc topologyLocation) setAllLocationsIfNull(locValue location.Type) {
	for i := 0; i < len(topoLoc); i++ {
		if topoLoc[i] == location.None {
			topoLoc[i] = locValue
		}
	}
}
func (topoLoc topologyLocation) allPositionsEqual(loc location.Type) bool {
	for _, l := range topoLoc {
		if l != loc {
			return false
		}
	}
	return true
}
func (topoLoc topologyLocation) merge(gl topologyLocation) topologyLocation {
	// if the src is an Area label & and the dest is not, increase the dest to be an Area
	if len(gl) > len(topoLoc) {
		newLoc := make([]location.Type, 3)
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

func (topoLoc topologyLocation) String() string {
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
