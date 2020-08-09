package main

import (
	"fmt"
	"github.com/a-deluna/gorenderer/v2/mat4"
	"github.com/a-deluna/gorenderer/v2/vec3"
	"github.com/a-deluna/gorenderer/v2/vec4"
	_ "github.com/ftrvxmtrx/tga"
	"github.com/g3n/engine/loader/obj"
	"github.com/veandco/go-sdl2/sdl"
	"image"
	"os"
)

var (
	width  int32 = 1024
	height int32 = 1024
)

type Screen struct {
	width       int32
	height      int32
	depth       float32
	framebuffer []Color
	depthbuffer []float32
  tileSize    int32
  tileCount    int32
}

type Color struct {
	a byte
	b byte
	g byte
	r byte
}

type Resources struct {
	obj   *obj.Decoder
	image image.Image
}

func main() {
	decoder, err := obj.Decode("african_head/african_head.obj", "")
	if err != nil {
		panic(err)
	}

	imageFile, err := os.Open("african_head/african_head_diffuse.tga")
	if err != nil {
		panic(err)
	}

	image, _, err := image.Decode(imageFile)

	resources := Resources{decoder, image}

	if err != nil {
		panic(err)
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

  tileSize:= int32(16)
  tileCount := width /tileSize * height/tileSize
	bufferSize := width * height * 4
	screen := Screen{width: width, height: height, depth: 10000,
		framebuffer: make([]Color, bufferSize),
		depthbuffer: make([]float32, bufferSize),
    tileSize: tileSize, tileCount: tileCount }

	mainLoop(&screen, &resources, renderer, texture)
}

func mainLoop(screen *Screen, resources *Resources,
	renderer *sdl.Renderer, texture *sdl.Texture) {
	running := true
	var input Input
	input.Init()
	keyStates := input.KeyStates()
	var camera Camera
  camera.Init()

  graphController := NewGraphController(screen.tileCount)
  completion := make(chan bool, screen.tileCount)

  for i := int32(0); i < screen.tileCount; i++ {
    graphStateChannel := make(chan GraphState)
    drawCommandsChannel := make(chan *DrawCommand, 100)
    graphController.AddTileChans(graphStateChannel, drawCommandsChannel)

    tr := TileRenderer{
      screen:screen,
      resources: resources,
      id: i,
      drawCommandsChannel : drawCommandsChannel,
      graphStateChannel : graphStateChannel,
      completion: completion }

    go tr.Render()
  }

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				for _, keystate := range keyStates {
					if keystate.scancode == e.Keysym.Scancode {
						keystate.held = e.Type == sdl.KEYDOWN
					}
				}
				break
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_RIGHT {
					input.mouse.moving = e.Type == sdl.MOUSEBUTTONDOWN
				}
				break
			case *sdl.MouseMotionEvent:
				input.mouse.dx = float32(e.XRel)
				input.mouse.dy = float32(e.YRel)
			}
		}

    msAtStart := sdl.GetTicks()

		camera.Update(&input)

		// clear buffers
		for i := int32(0); i < screen.height*screen.width; i++ {
			screen.framebuffer[i] = Color{0, 0, 0, 0}
			screen.depthbuffer[i] = screen.depth * 2
		}

		Draw(screen, resources, &camera, &graphController)

    completionCount := int32(0)
    for completionCount < screen.tileCount {
      <-completion
      completionCount++
    }

		pixels, _, err := texture.Lock(nil)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(pixels); i += 4 {
			color := screen.framebuffer[i/4]
			pixels[i+0] = color.a
			pixels[i+1] = color.b
			pixels[i+2] = color.g
			pixels[i+3] = color.r
		}

		texture.Unlock()

		err = renderer.CopyEx(texture, nil, nil, 0, nil, sdl.FLIP_VERTICAL)
		if err != nil {
			panic(err)
		}

		renderer.Present()
    //panic("wow")
    msAtEnd := sdl.GetTicks()
    difference := msAtEnd-msAtStart
    fmt.Printf("Time: %d ms\n",difference)
	}
}

func Draw(screen *Screen, resources *Resources, camera *Camera,
          graphController *GraphController) {
	mesh := resources.obj.Objects[0]
	// Each face consists of 3 vertices
	modelVerts := resources.obj.Vertices

	vp := camera.VP()

  //fmt.Printf("setting state to running\n")
  graphController.ChangeState(RUNNING)

	for _, face := range mesh.Faces {
		v1Orig := toVec(modelVerts, face.Vertices[0]*3)
		v2Orig := toVec(modelVerts, face.Vertices[1]*3)
		v3Orig := toVec(modelVerts, face.Vertices[2]*3)

		v1_homo := mat4.VectorMultiply(vp, vec4.FromVec3(v1Orig))
		v2_homo := mat4.VectorMultiply(vp, vec4.FromVec3(v2Orig))
		v3_homo := mat4.VectorMultiply(vp, vec4.FromVec3(v3Orig))

		if shouldClip(v1_homo, v2_homo, v3_homo) {
			continue
		}

    v1 := v1_homo.Wdivide()
    v2 := v2_homo.Wdivide()
    v3 := v3_homo.Wdivide()

		v1Screen := vec3.Vec3{
			(v1[0] + 1.0) * float32(screen.width) / 2.0,
			(v1[1] + 1.0) * float32(screen.height) / 2.0,
			(v1[2]) * float32(screen.depth),
		}
		v2Screen := vec3.Vec3{
			(v2[0] + 1.0) * float32(screen.width) / 2.0,
			(v2[1] + 1.0) * float32(screen.height) / 2.0,
			(v2[2]) * float32(screen.depth),
		}
		v3Screen := vec3.Vec3{
			(v3[0] + 1.0) * float32(screen.width) / 2.0,
			(v3[1] + 1.0) * float32(screen.height) / 2.0,
			(v3[2]) * float32(screen.depth),
		}

		normal := vec3.Normalize(vec3.Cross(v3.Sub(v1), v2.Sub(v1)))

		t1 := toVec2(resources.obj.Uvs, face.Uvs[0]*2)
		t2 := toVec2(resources.obj.Uvs, face.Uvs[1]*2)
		t3 := toVec2(resources.obj.Uvs, face.Uvs[2]*2)

    lightDir := vec3.Vec3{0,0,1}
    magnitude := -vec3.Dot(lightDir,normal)
    if magnitude < 0 {
      continue
    }


   box := boundingBox(v1Screen, v2Screen, v3Screen)

   dc := DrawCommand{v1Screen,v2Screen,v3Screen,t1,t2,t3}
   tilesIds := getTileIds(screen, &box)
   for _, id := range tilesIds {
     //fmt.Printf("sending draw command to %d\n", id)
     graphController.NotifyTile(id, &dc)
   }

	}
  graphController.ChangeState(DONE)
}

func getTileIds(screen *Screen, box *Box) []int32 {
  mask := screen.tileSize-1
  tilesPerRow := screen.width / screen.tileSize

  minx := (int32(box.minx) & ^mask) / screen.tileSize
  miny := (int32(box.miny) & ^mask) / screen.tileSize
  maxx := (int32(box.maxx) & ^mask) / screen.tileSize
  maxy := (int32(box.maxy) & ^mask) / screen.tileSize

  list := make([]int32, 0)
  for x := minx; x <= maxx; x++ {
    for y := miny; y <= maxy; y++ {
      id := x + tilesPerRow * y
      list = append(list, id)
    }
  }
  //fmt.Printf("box: %#v\n", box)
  //fmt.Printf("mask: %d minx: %d miny: %d maxx: %d maxy: %d\n",
  //          mask,minx,miny,maxx,maxy)
  //fmt.Printf("list: %v\n", list)
  return list
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

