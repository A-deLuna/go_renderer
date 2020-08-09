package main

import (
	"github.com/a-deluna/gorenderer/v2/mat4"
	"github.com/a-deluna/gorenderer/v2/vec3"
	"github.com/a-deluna/gorenderer/v2/vec4"
)

type Simulation struct {
	input           Input
	camera          Camera
	graphController GraphController
	tileCompletion  chan bool
	screen          *Screen
	resources       *Resources
}

func NewSimulation(screen *Screen, resources *Resources) (s Simulation) {
	s.input = NewInput()
	s.camera = NewCamera()
	s.screen = screen
	s.resources = resources
	s.graphController = NewGraphController(screen.tileCount)
	s.tileCompletion = make(chan bool, screen.tileCount)

	for i := int32(0); i < screen.tileCount; i++ {
		graphStateChannel := make(chan GraphState)
		drawCommandsChannel := make(chan *DrawCommand, 100)
		s.graphController.AddTileChans(graphStateChannel, drawCommandsChannel)

		tr := TileRenderer{
			screen:              s.screen,
			resources:           s.resources,
			id:                  i,
			drawCommandsChannel: drawCommandsChannel,
			graphStateChannel:   graphStateChannel,
			completion:          s.tileCompletion}

		go tr.Render()
	}
	return
}

func (s *Simulation) RenderNextFrame() {
	s.camera.Update(&s.input)

	s.ClearBuffers()
	s.SendDrawCommands()
	s.WaitForTiles()
}

func (s *Simulation) ClearBuffers() {
	// clear buffers
	for i := int32(0); i < s.screen.height*s.screen.width; i++ {
		s.screen.framebuffer[i] = Color{0, 0, 0, 0}
		s.screen.depthbuffer[i] = s.screen.depth * 2
	}
}

func (s *Simulation) WaitForTiles() {
	completionCount := int32(0)
	for completionCount < s.screen.tileCount {
		<-s.tileCompletion
		completionCount++
	}
}

func (s *Simulation) SendDrawCommands() {
	mesh := s.resources.obj.Objects[0]
	modelVerts := s.resources.obj.Vertices

	vp := s.camera.VP()

	s.graphController.ChangeState(RUNNING)

	for _, face := range mesh.Faces {
		v1Orig := toVec4(modelVerts, face.Vertices[0]*3)
		v2Orig := toVec4(modelVerts, face.Vertices[1]*3)
		v3Orig := toVec4(modelVerts, face.Vertices[2]*3)

		v1_homo := mat4.VectorMultiply(vp, v1Orig)
		v2_homo := mat4.VectorMultiply(vp, v2Orig)
		v3_homo := mat4.VectorMultiply(vp, v3Orig)

		if shouldClip(v1_homo, v2_homo, v3_homo) {
			continue
		}

		v1 := v1_homo.Wdivide()
		v2 := v2_homo.Wdivide()
		v3 := v3_homo.Wdivide()

		v1Screen := vec3.Vec3{
			(v1[0] + 1.0) * float32(s.screen.width) / 2.0,
			(v1[1] + 1.0) * float32(s.screen.height) / 2.0,
			(v1[2]) * float32(s.screen.depth),
		}
		v2Screen := vec3.Vec3{
			(v2[0] + 1.0) * float32(s.screen.width) / 2.0,
			(v2[1] + 1.0) * float32(s.screen.height) / 2.0,
			(v2[2]) * float32(s.screen.depth),
		}
		v3Screen := vec3.Vec3{
			(v3[0] + 1.0) * float32(s.screen.width) / 2.0,
			(v3[1] + 1.0) * float32(s.screen.height) / 2.0,
			(v3[2]) * float32(s.screen.depth),
		}

		normal := vec3.Normalize(vec3.Cross(v3.Sub(v1), v2.Sub(v1)))

		t1 := toVec2(s.resources.obj.Uvs, face.Uvs[0]*2)
		t2 := toVec2(s.resources.obj.Uvs, face.Uvs[1]*2)
		t3 := toVec2(s.resources.obj.Uvs, face.Uvs[2]*2)

		lightDir := vec3.Vec3{0, 0, 1}
		magnitude := -vec3.Dot(lightDir, normal)
		if magnitude < 0 {
			continue
		}

		box := boundingBox(v1Screen, v2Screen, v3Screen)

		dc := DrawCommand{v1Screen, v2Screen, v3Screen, t1, t2, t3}
		tilesIds := getTileIds(s.screen, &box)
		for _, id := range tilesIds {
			s.graphController.NotifyTile(id, &dc)
		}
	}
	s.graphController.ChangeState(DONE)
}

func toVec4(verts []float32, base int) vec4.Vec4 {
	return vec4.Vec4{verts[base], verts[base+1], verts[base+2], 1}
}

func shouldClip(vectors ...vec4.Vec4) bool {
	for _, v := range vectors {
		w := v[3]
		acceptable := -w <= v[0] && v[0] <= w &&
			-w <= v[1] && v[1] <= w &&
			-w <= v[2] && v[2] <= w
		if !acceptable {
			return true
		}
	}
	return false
}

func getTileIds(screen *Screen, box *Box) []int32 {
	mask := screen.tileSize - 1
	tilesPerRow := screen.width / screen.tileSize

	minx := (int32(box.minx) & ^mask) / screen.tileSize
	miny := (int32(box.miny) & ^mask) / screen.tileSize
	maxx := (int32(box.maxx) & ^mask) / screen.tileSize
	maxy := (int32(box.maxy) & ^mask) / screen.tileSize

	list := make([]int32, 0)
	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			id := x + tilesPerRow*y
			list = append(list, id)
		}
	}
	return list
}
