package main

import(
	"github.com/a-deluna/gorenderer/v2/vec3"
)
type DrawCommand struct {
  v1,v2,v3,t1,t2,t3 vec3.Vec3
}

type GraphState int32
const (
  INIT GraphState = iota
  RUNNING
  DONE
)

type TileRenderer struct {
  screen *Screen
  resources *Resources
  id int32
  drawCommandsChannel <-chan *DrawCommand
  graphStateChannel <-chan GraphState
  completion chan<- bool
  currentState GraphState
}


func (tr *TileRenderer) Render() {
  for {
    select {
    case state := <-tr.graphStateChannel:
      //fmt.Printf("Received graph state update in renderer %d\n", tr.id)
      tr.currentState = state;
      if len(tr.drawCommandsChannel) == 0 && tr.currentState == DONE {
        tr.completion <- true
      }
    case dc := <-tr.drawCommandsChannel:
      //fmt.Printf("Drawing on tile %d\n",tr.id)
        tr.drawTriangle(dc)
        if len(tr.drawCommandsChannel) == 0 && tr.currentState == DONE {
          tr.completion <- true
        }
    }
  }
}

func (tr *TileRenderer) GetCoords() (x int32, y int32) {
  tilesPerLine := tr.screen.width / tr.screen.tileSize;
  x = tr.id & (tilesPerLine-1)
  y = tr.id / tilesPerLine
  return;
}

func (tr *TileRenderer) drawTriangle(dc *DrawCommand) {
  tileX, tileY := tr.GetCoords()
  tileSize := tr.screen.tileSize


	for y := tileY *tileSize; y < (tileY+1) * tileSize; y++ {
		for x := tileX * tileSize; x < (tileX+1) * tileSize; x++ {
			b := baricenter(vec3.Vec3{float32(x), float32(y), 1}, dc.v1, dc.v2, dc.v3)
			if pointInTriangle(b) {
        //fmt.Printf("found point in triangle")
				pointz := dc.v1[2]*b[0] + dc.v2[2]*b[1] + dc.v3[2]*b[2]
				//fmt.Printf("pointz: %v\n", pointz)
				depthIndex := int32(x) + int32(y)*tr.screen.width
				if pointz < tr.screen.depthbuffer[depthIndex] {
					tr.screen.depthbuffer[depthIndex] = pointz
					textureCoords := vec3.Add(
						vec3.Scale(dc.t1, b[0]),
						vec3.Scale(dc.t2, b[1]),
						vec3.Scale(dc.t3, b[2]))
					screenTexCoords := []int{
						int(textureCoords[0] * float32(tr.resources.image.Bounds().Dx())),
						int((1.0 - textureCoords[1]) * float32(tr.resources.image.Bounds().Dy()))}
					color := tr.resources.image.At(screenTexCoords[0], screenTexCoords[1])
					r, g, b, a := color.RGBA()

          //fmt.Printf("writing o the colo buffer\n")
					tr.screen.framebuffer[int32(x)+int32(y)*tr.screen.width] =
						Color{byte(a), byte(b), byte(g), byte(r)}
				}
			}
		}
	}
}

func baricenter(p, v1, v2, v3 vec3.Vec3) vec3.Vec3 {
	v31 := v1.Sub(v3)
	v32 := v2.Sub(v3)
	pv3 := v3.Sub(p)

	x := vec3.Vec3{v31[0], v32[0], pv3[0]}
	y := vec3.Vec3{v31[1], v32[1], pv3[1]}

	cross := vec3.Cross(x, y)

	u := cross[0] / cross[2]
	v := cross[1] / cross[2]

	return vec3.Vec3{u, v, 1.0 - u - v}
}

func pointInTriangle(baricenter vec3.Vec3) bool {
	return baricenter[0] >= 0 && baricenter[1] >= 0 &&
		baricenter[0]+baricenter[1] <= 1.0
}
