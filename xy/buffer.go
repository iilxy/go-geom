package xy

import "github.com/twpayne/go-geom"

type EndCapStyle int

const (
	RoundCap EndCapStyle = iota
	FlatCap
	SquareCap
)

type JointStyle int

const (
	RoundJoint JointStyle = iota
	MitreJoint
	BevelJoint
)

type BufferParams struct {
	Geometry geom.T
	Distance float64
	// QuadrantSegments is the number of line segments used to approximate an angle fillet.
	// Default is 8
	QuadrantSegments int
	// Default is Round
	EndCapStyle EndCapStyle
	// Default is Round
	JoinStyle JointStyle
	// MitreLimit
	// Default is 5
	MitreLimit    float64
	IsSingleSided bool
}

//func Buffer(g geom.T, distance float64) geom.T {
//	return BufferWithParams(BufferParams{Geometry:g, Distance:distance})
//}
//func BufferWithParams(params BufferParams) (buffered geom.T) {
//	if params.QuadrantSegments == 0 {
//		params.QuadrantSegments = 8;
//	}
//
//	buffered = bufferOriginalPrecision(params)
//
//	if buffered == nil {
//		// TODO implement precision model and execute it in precision model
//	}
//	return nil
//}
//func bufferOriginalPrecision(params BufferParams) geom.T {
//	curves := buildCurves(params)
//}
//func buildCurves(params BufferParams) interface{} {
//}
