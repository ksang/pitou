package board

import "testing"

func TestDistToSegmentSquared(t *testing.T) {
	var tests = []struct {
		p, v, w  Vector
		expected float64
	}{
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 1.0},
			w:        Vector{x: 1.0, y: 0.0},
			expected: 1.0,
		},
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 4.0},
			w:        Vector{x: 1.0, y: 0.0},
			expected: 0.0,
		},
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 0.0},
			w:        Vector{x: 1.0, y: 1.0},
			expected: 1.0,
		},
		{
			p:        Vector{x: 1.0, y: 1.0},
			v:        Vector{x: 0.0, y: 2.0},
			w:        Vector{x: 2.0, y: 0.0},
			expected: 0.0,
		},
	}
	for caseid, c := range tests {
		ret := distToSegmentSquared(c.p, c.v, c.w)
		if ret != c.expected {
			t.Errorf("case #%d failed, result: %f, expected: %f\n",
				caseid+1, ret, c.expected)
		}
	}
}

func TestDistToLineSquared(t *testing.T) {
	var tests = []struct {
		p, v, w  Vector
		expected float64
	}{
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 1.0},
			w:        Vector{x: 1.0, y: 0.0},
			expected: 0.0,
		},
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 4.0},
			w:        Vector{x: 1.0, y: 0.0},
			expected: 0.0,
		},
		{
			p:        Vector{x: 1.0, y: 2.0},
			v:        Vector{x: 1.0, y: 0.0},
			w:        Vector{x: 1.0, y: 1.0},
			expected: 0.0,
		},
		{
			p:        Vector{x: 1.0, y: 1.0},
			v:        Vector{x: 0.0, y: 2.0},
			w:        Vector{x: 2.0, y: 0.0},
			expected: 0.0,
		},
		{
			p:        Vector{x: 3.0, y: 0.0},
			v:        Vector{x: 2.0, y: 2.0},
			w:        Vector{x: 2.0, y: 0.0},
			expected: 1.0,
		},
		{
			p:        Vector{x: 4.0, y: 2.0},
			v:        Vector{x: 2.0, y: 2.0},
			w:        Vector{x: 2.0, y: 2.0},
			expected: 4.0,
		},
	}
	for caseid, c := range tests {
		ret := distToLineSquared(c.p, c.v, c.w)
		if ret != c.expected {
			t.Errorf("case #%d failed, result: %f, expected: %f\n",
				caseid+1, ret, c.expected)
		}
	}
}
