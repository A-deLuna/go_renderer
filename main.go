package main

import (
	"fmt"
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
	tileCount   int32
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

	tileSize := int32(16)
	tileCount := width / tileSize * height / tileSize
	bufferSize := width * height * 4
	screen := Screen{width: width, height: height, depth: 10000,
		framebuffer: make([]Color, bufferSize),
		depthbuffer: make([]float32, bufferSize),
		tileSize:    tileSize, tileCount: tileCount}

	mainLoop(&screen, &resources, renderer, texture)
}

func mainLoop(screen *Screen, resources *Resources,
	renderer *sdl.Renderer, texture *sdl.Texture) {

	simulation := NewSimulation(screen, resources)
	keyStates := simulation.input.KeyStates()

	running := true
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
					simulation.input.mouse.moving = e.Type == sdl.MOUSEBUTTONDOWN
				}
				break
			case *sdl.MouseMotionEvent:
				simulation.input.mouse.dx = float32(e.XRel)
				simulation.input.mouse.dy = float32(e.YRel)
			}
		}

		msAtStart := sdl.GetTicks()

		simulation.RenderNextFrame()

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
		msAtEnd := sdl.GetTicks()
		difference := msAtEnd - msAtStart
		fmt.Printf("Time: %d ms\n", difference)
	}
}
