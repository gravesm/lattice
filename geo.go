package main

import "math"

func boundsFromTile(x, y, z int) (float64, float64, float64, float64) {
	n, w := tileToNw(x, y, z)
	s, e := tileToNw(x+1, y+1, z)
	return w, s, e, n
}

func tileToNw(x, y, z int) (float64, float64) {
	//return upper left corner
	n := math.Exp2(float64(z))
	lon := float64(x)/n*360.0 - 180.0
	lat := (180 / math.Pi) * math.Atan(math.Sinh(math.Pi*(1-2*float64(y)/n)))
	return lat, lon
}
