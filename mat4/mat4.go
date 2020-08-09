package mat4

import (
	"github.com/a-deluna/gorenderer/v2/vec3"
	"github.com/a-deluna/gorenderer/v2/vec4"
	"math"
)

type Mat4 [4][4]float32

func Multiply(a, b Mat4) (out Mat4) {

	out[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0] + a[0][3]*b[3][0]
	out[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0] + a[1][3]*b[3][0]
	out[2][0] = a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0] + a[2][3]*b[3][0]
	out[3][0] = a[3][0]*b[0][0] + a[3][1]*b[1][0] + a[3][2]*b[2][0] + a[3][3]*b[3][0]

	out[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1] + a[0][3]*b[3][1]
	out[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1] + a[1][3]*b[3][1]
	out[2][1] = a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1] + a[2][3]*b[3][1]
	out[3][1] = a[3][0]*b[0][1] + a[3][1]*b[1][1] + a[3][2]*b[2][1] + a[3][3]*b[3][1]

	out[0][2] = a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2] + a[0][3]*b[3][2]
	out[1][2] = a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2] + a[1][3]*b[3][2]
	out[2][2] = a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2] + a[2][3]*b[3][2]
	out[3][2] = a[3][0]*b[0][2] + a[3][1]*b[1][2] + a[3][2]*b[2][2] + a[3][3]*b[3][2]

	out[0][3] = a[0][0]*b[0][3] + a[0][1]*b[1][3] + a[0][2]*b[2][3] + a[0][3]*b[3][3]
	out[1][3] = a[1][0]*b[0][3] + a[1][1]*b[1][3] + a[1][2]*b[2][3] + a[1][3]*b[3][3]
	out[2][3] = a[2][0]*b[0][3] + a[2][1]*b[1][3] + a[2][2]*b[2][3] + a[2][3]*b[3][3]
	out[3][3] = a[3][0]*b[0][3] + a[3][1]*b[1][3] + a[3][2]*b[2][3] + a[3][3]*b[3][3]

	return
}

func VectorMultiply(a Mat4, v vec4.Vec4) (out vec4.Vec4) {
	out = vec4.Vec4{0, 0, 0, 0}
	out[0] = a[0][0]*v[0] + a[0][1]*v[1] + a[0][2]*v[2] + a[0][3]*v[3]
	out[1] = a[1][0]*v[0] + a[1][1]*v[1] + a[1][2]*v[2] + a[1][3]*v[3]
	out[2] = a[2][0]*v[0] + a[2][1]*v[1] + a[2][2]*v[2] + a[2][3]*v[3]
	out[3] = a[3][0]*v[0] + a[3][1]*v[1] + a[3][2]*v[2] + a[3][3]*v[3]
	return
}

func Translation(v vec3.Vec3) Mat4 {
	return Mat4{
		{1.0, 0.0, 0.0, v[0]},
		{0.0, 1.0, 0.0, v[1]},
		{0.0, 0.0, 1.0, v[2]},
		{0.0, 0.0, 0.0, 1.0},
	}
}

func Transpose(a Mat4) (out Mat4) {
	out[0][0] = a[0][0]
	out[0][1] = a[1][0]
	out[0][2] = a[2][0]
	out[0][3] = a[3][0]
	out[1][0] = a[0][1]
	out[1][1] = a[1][1]
	out[1][2] = a[2][1]
	out[1][3] = a[3][1]
	out[2][0] = a[0][2]
	out[2][1] = a[1][2]
	out[2][2] = a[2][2]
	out[2][3] = a[3][2]
	out[3][0] = a[0][3]
	out[3][1] = a[1][3]
	out[3][2] = a[2][3]
	out[3][3] = a[3][3]
	return
}

func Rotation(x, y, z float64) Mat4 {
	sinx := float32(math.Sin(x))
	cosx := float32(math.Cos(x))
	rotX := Mat4{
		{1, 0, 0, 0},
		{0, cosx, -sinx, 0},
		{0, sinx, cosx, 0},
		{0, 0, 0, 1},
	}

	siny := float32(math.Sin(y))
	cosy := float32(math.Cos(y))
	rotY := Mat4{
		{cosy, 0, siny, 0},
		{0, 1, 0, 0},
		{-siny, 0, cosy, 0},
		{0, 0, 0, 1},
	}

	sinz := float32(math.Sin(z))
	cosz := float32(math.Cos(z))
	rotZ := Mat4{
		{cosz, -sinz, 0, 0},
		{sinz, cosz, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	return Multiply(rotZ, Multiply(rotY, rotX))
}

func Projection(n, f float32, fov float64) Mat4 {

	fustrum := n * float32(math.Tan(fov*(math.Pi/180.0*.5)))
	l := -fustrum
	r := fustrum
	b := -fustrum
	t := fustrum

	return Mat4{
		{(2 * n) / (r - l), 0, 0, 0},
		{0, (2 * n) / (t - b), 0, 0},
		{(r + l) / (r - l), (t + b) / (t - b), -((f + n) / (f - n)), -1},
		{0, 0, -((2 * f * n) / (f - n)), 0},
	}
}
