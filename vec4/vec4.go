package vec4

import (
	"fmt"
	"github.com/a-deluna/gorenderer/v2/vec3"
)

type Vec4 []float32

func FromVec3(v vec3.Vec3) Vec4 {
	return Vec4{v[0], v[1], v[2], 1.0}
}

func (v Vec4) Wdivide() vec3.Vec3 {
	return vec3.Vec3{v[0] / v[3], v[1] / v[3], v[2] / v[3]}
}

func (v Vec4) String() string {
	return fmt.Sprintf("{x: %f, y:%f, z:%f, w:%f}", v[0], v[1], v[2], v[3])
}
