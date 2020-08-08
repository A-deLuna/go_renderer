package main

import (
	"github.com/veandco/go-sdl2/sdl"
  "github.com/g3n/engine/loader/obj"
  "image"
  "os"
  _ "github.com/ftrvxmtrx/tga"
  "fmt"
  "github.com/a-deluna/gorenderer/v2/vec3"
  "github.com/a-deluna/gorenderer/v2/vec4"
  "github.com/a-deluna/gorenderer/v2/mat4"
)

var (
	width  int32 = 1024
	height int32 = 1024
)

type Screen struct {
  width int32
  height int32
  depth float32
  framebuffer []Color
  depthbuffer []float32
}

type Color struct {
  a byte
  b byte
  g byte
  r byte
}


type Resources struct {
  obj *obj.Decoder
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

  image, _, err := image.Decode(imageFile);

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

  bufferSize := width * height * 4
  screen := Screen{width: width, height: height, depth:10000,
  framebuffer: make([]Color, bufferSize),
  depthbuffer: make([]float32, bufferSize)}

  for i := int32(0); i < width; i++ {
    base := (i * width * 4) + i * 4
    screen.framebuffer[base] = Color{0xff,0,0,0xff}
  }

  mainLoop(&screen, &resources, renderer, texture)
}


type Input struct {
  f int
  h int

  forward KeyState
  backward KeyState
  left KeyState
  right KeyState
  up KeyState
  down KeyState
}

func (i *Input) Init() {
  i.forward = KeyState{sdl.SCANCODE_W, false}
  i.backward = KeyState{sdl.SCANCODE_S, false}
  i.left = KeyState{sdl.SCANCODE_A, false}
  i.right = KeyState{sdl.SCANCODE_D, false}
  i.up = KeyState{sdl.SCANCODE_Q, false}
  i.down = KeyState{sdl.SCANCODE_E, false}
}

type KeyState struct {
  scancode sdl.Scancode
  held bool
}

type Camera struct {
  position vec3.Vec3
  yaw float32
  pitch float32
}

func (c *Camera) Init() {
  c.position = vec3.Vec3{0,0,3}
}

func (c *Camera) Update(input *Input) {
  var dx,dy,dz float32
  if input.forward.held && !input.backward.held {
    dz = -.3
  }
  if !input.forward.held && input.backward.held {
    dz = .3
  }
  if input.left.held && !input.right.held {
    dx = -.3
  }
  if !input.left.held && input.right.held {
    dx = .3
  }
  if !input.up.held && input.down.held {
    dy = .3
  }
  if input.up.held && !input.down.held {
    dy = -.3
  }
  c.position[0] += dx
  c.position[1] += dy
  c.position[2] += dz
}

func mainLoop(screen *Screen, resources *Resources,
              renderer *sdl.Renderer, texture *sdl.Texture) {
	running := true
  var input Input
  input.Init()
  keyStates := []*KeyState{&input.forward, &input.backward,
      &input.left, &input.right, &input.up, &input.down}

  var camera Camera
  camera.Init()

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch e:= event.(type) {
      case *sdl.QuitEvent:
        running = false
        break
      case *sdl.KeyboardEvent:
        for _, keystate := range keyStates {
          if keystate.scancode == e.Keysym.Scancode {
            keystate.held = e.Type == sdl.KEYDOWN
          }
        }
      }
		}

    // update camera
    camera.Update(&input)

    // clear buffers
    for i := int32(0); i < screen.height * screen.width; i++ {
      screen.framebuffer[i] = Color{0,0,0,0}
      screen.depthbuffer[i] = 0
    }

    Draw(screen, resources, &camera)


    pixels, _, err := texture.Lock(nil)
    if err != nil {
      panic(err)
    }

    for i := 0; i < len(pixels); i+=4 {
      color := screen.framebuffer[i/4]
      pixels[i + 0] = color.a
      pixels[i + 1] = color.b
      pixels[i + 2] = color.g
      pixels[i + 3] = color.r
    }

    texture.Unlock()

    err = renderer.CopyEx(texture, nil, nil, 0, nil, sdl.FLIP_VERTICAL)
    if err != nil {
      panic (err)
    }

    renderer.Present()

	}
}

func Draw(screen *Screen, resources *Resources, camera *Camera) {
  mesh := resources.obj.Objects[0]
  // Each face consists of 3 vertices
  modelVerts := resources.obj.Vertices


  center := vec3.Vec3{0,0,0}
  cameraPosition := camera.position
  fmt.Println(cameraPosition)

  camTrans := mat4.Translation(center.Sub(cameraPosition))

  projection := mat4.Projection(1, 50000, 30)

  mvp := mat4.Multiply(projection, camTrans)

  for _, face := range mesh.Faces {
    v1 := toVec(modelVerts, face.Vertices[0] * 3)
    v2 := toVec(modelVerts, face.Vertices[1] * 3)
    v3 := toVec(modelVerts, face.Vertices[2] * 3)

    v1_homo := mat4.VectorMultiply(mvp,vec4.FromVec3(v1))
    v2_homo := mat4.VectorMultiply(mvp,vec4.FromVec3(v2))
    v3_homo := mat4.VectorMultiply(mvp,vec4.FromVec3(v3))

    if shouldClip(v1_homo, v2_homo, v3_homo) {
      continue
    }

    v1 = v1_homo.Wdivide();
    v2 = v2_homo.Wdivide();
    v3 = v3_homo.Wdivide();


    v1Screen := vec3.Vec3 {
      (v1[0] + 1.0) * float32(screen.width) / 2.0,
      (v1[1] + 1.0) * float32(screen.height) / 2.0,
      (v1[2] + 1.0) * float32(screen.depth) / 2.0,
    };
    v2Screen := vec3.Vec3{
      (v2[0] + 1.0) * float32(screen.width) / 2.0,
      (v2[1] + 1.0) * float32(screen.height) / 2.0,
      (v2[2] + 1.0) * float32(screen.depth) / 2.0,
    };
    v3Screen := vec3.Vec3{
      (v3[0] + 1.0) * float32(screen.width) / 2.0,
      (v3[1] + 1.0) * float32(screen.height) / 2.0,
      (v3[2] + 1.0) * float32(screen.depth) / 2.0,
    };

    normal := vec3.Normalize(vec3.Cross(v3.Sub(v1), v2.Sub(v1)))


    tex1 := toVec2(resources.obj.Uvs, face.Uvs[0] * 2)
    tex2 := toVec2(resources.obj.Uvs, face.Uvs[1] * 2)
    tex3 := toVec2(resources.obj.Uvs, face.Uvs[2] * 2)

    drawTriangle(screen, v1Screen, v2Screen, v3Screen, normal, tex1, tex2, tex3,
      resources.image)

  }
}

func shouldClip(vectors... vec4.Vec4) bool {
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

type Box struct {
  minx float32
  miny float32
  maxx float32
  maxy float32
}

func drawTriangle(screen *Screen, v1, v2, v3, normal,
                  t1, t2, t3 vec3.Vec3, image image.Image) {
      box:= boundingBox(v1,v2,v3)
      lightDir := vec3.Vec3{0,0,1}

      magnitude := -vec3.Dot(lightDir, normal)
      if magnitude < 0 {
        return
      }

      for y := int(box.miny); y <= int(box.maxy); y++ {
        for x := int(box.minx); x <= int(box.maxx); x++ {
          b := baricenter(vec3.Vec3{float32(x),float32(y),1}, v1, v2, v3)
          if pointInTriangle(b) {
            pointz := v1[2] * b[0] + v2[2] * b [1] + v3[2] * b[2]
            //fmt.Printf("pointz: %v\n", pointz)
            depthIndex := int32(x) + int32(y) * screen.width

            if pointz > screen.depthbuffer[depthIndex] {
              screen.depthbuffer[depthIndex] = pointz
              textureCoords := vec3.Add(
               vec3.Scale(t1, b[0]),
               vec3.Scale(t2, b[1]),
               vec3.Scale(t3, b[2]))
              screenTexCoords := []int{
                int(textureCoords[0] * float32(image.Bounds().Dx())),
                int((1.0 - textureCoords[1]) * float32(image.Bounds().Dy()))}
              color:= image.At(screenTexCoords[0], screenTexCoords[1])
              r, g, b, a:= color.RGBA()
              screen.framebuffer[int32(x) + int32(y) * screen.width] =
                  Color{byte(a),byte(b),byte(g),byte(r)}
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

  return vec3.Vec3{u, v, 1.0-u-v}
}

func pointInTriangle(baricenter vec3.Vec3) bool {
  return baricenter[0] >= 0 && baricenter[1] >= 0 &&
      baricenter[0] + baricenter[1] <= 1.0
}

func boundingBox(v1, v2, v3 vec3.Vec3) Box {
  minx := min(v1[0], min(v2[0], v3[0]))
  miny := min(v1[1], min(v2[1], v3[1]))
  maxx := max(v1[0], max(v2[0], v3[0]))
  maxy := max(v1[1], max(v2[1], v3[1]))
  return Box{minx,miny,maxx,maxy}
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
