package main

import(
	"github.com/a-deluna/gorenderer/v2/mat4"
	"github.com/a-deluna/gorenderer/v2/vec3"
	"github.com/a-deluna/gorenderer/v2/vec4"
)

type Camera struct {
	position vec3.Vec3
	yaw      float32
	pitch    float32
}

func (c *Camera) Init() {
	c.position = vec3.Vec3{0, 0, 2}
}

func (c *Camera) Update(input *Input) {
  v := vec4.Vec4{0,0,0,0}
  delta := float32(.15)
	if input.keyboard.forward.held && !input.keyboard.backward.held {
		v[2] = -delta
	}
	if !input.keyboard.forward.held && input.keyboard.backward.held {
		v[2] = delta
	}
	if input.keyboard.left.held && !input.keyboard.right.held {
		v[0] = -delta
	}
	if !input.keyboard.left.held && input.keyboard.right.held {
		v[0] = delta
	}
	if !input.keyboard.up.held && input.keyboard.down.held {
		v[1] = delta
	}
	if input.keyboard.up.held && !input.keyboard.down.held {
		v[1] = -delta
	}
  rotation := mat4.Rotation(float64(c.pitch), float64(c.yaw), 0)
  v = mat4.VectorMultiply(rotation, v)

	c.position[0] += v[0]
	c.position[1] += v[1]
	c.position[2] += v[2]

	if input.mouse.moving {
		c.yaw += input.mouse.dx * .003
		c.pitch += input.mouse.dy * .003
	}
}

func (c *Camera) VP() mat4.Mat4 {
	center := vec3.Vec3{0, 0, 0}

	trans := mat4.Translation(center.Sub(c.position))
	rot := mat4.Transpose(mat4.Rotation(float64(c.pitch), float64(c.yaw), 0))

  var near float32 = 1
  var far float32  = 10
  var fov float64 =  30
	projection := mat4.Projection(near,far,fov)

  return mat4.Multiply(projection, mat4.Multiply(rot, trans))
}
