package board

// vector represents a point on 2D board
type Vector struct {
	x float64
	y float64
}

func sqr(x float64) float64 {
	return x * x
}

// distance squared between two points
func distSquared(v, w Vector) float64 {
	return sqr(v.x-w.x) + sqr(v.y-w.y)
}

// distance squared between a point and a segment
func distToSegmentSquared(p, v, w Vector) float64 {
	lenSq := distSquared(v, w)
	dot := (p.x-v.x)*(w.x-v.x) + (p.y-v.y)*(w.y-v.y)
	var (
		param, xx, yy float64
	)
	param = -1
	if lenSq != 0 {
		param = dot / lenSq
	}
	if param < 0 {
		xx = v.x
		yy = v.y
	} else if param > 1 {
		xx = w.x
		yy = w.y
	} else {
		xx = v.x + param*(w.x-v.x)
		yy = v.y + param*(w.y-v.y)
	}
	return distSquared(p, Vector{x: xx, y: yy})
}

// distance squared between a point and a line
func distToLineSquared(p, v, w Vector) float64 {
	lenSq := distSquared(v, w)
	if lenSq == 0 {
		return distSquared(p, v)
	}
	dot := (p.x-v.x)*(w.x-v.x) + (p.y-v.y)*(w.y-v.y)
	param := dot / lenSq
	xx := v.x + param*(w.x-v.x)
	yy := v.y + param*(w.y-v.y)
	return distSquared(p, Vector{x: xx, y: yy})
}
