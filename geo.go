package main

import "math"

func boundsFromTile(x, y, z int) (float64, float64, float64, float64) {
	m := 20037508.34
	r := m * 2 / math.Exp2(float64(z))
	w := -m + r*float64(x)
	s := m - r*float64(y)
	e := -m + r*float64(x) + r
	n := m - r*float64(y) - r
	return w, s, e, n
}
