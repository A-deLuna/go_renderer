package main

import (
	"github.com/a-deluna/gorenderer/v2/vec3"
)

type Box struct {
	minx float32
	miny float32
	maxx float32
	maxy float32
}

func boundingBox(v1, v2, v3 vec3.Vec3) Box {
	minx := min(v1[0], min(v2[0], v3[0]))
	miny := min(v1[1], min(v2[1], v3[1]))
	maxx := max(v1[0], max(v2[0], v3[0]))
	maxy := max(v1[1], max(v2[1], v3[1]))
	return Box{minx, miny, maxx, maxy}
}

func toVec2(verts []float32, base int) vec3.Vec3 {
	return vec3.Vec3{verts[base], verts[base+1], 0}
}

func toVec(verts []float32, base int) vec3.Vec3 {
	return vec3.Vec3{verts[base], verts[base+1], verts[base+2]}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
