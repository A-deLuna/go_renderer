package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Input struct {
	keyboard KeyboardInput
	mouse    MouseInput
}

type KeyboardInput struct {
	forward  KeyState
	backward KeyState
	left     KeyState
	right    KeyState
	up       KeyState
	down     KeyState
}

type KeyState struct {
	scancode sdl.Scancode
	held     bool
}

type MouseInput struct {
	moving bool
	dx, dy float32
}


func (i *Input) Init() {
	i.keyboard.forward = KeyState{sdl.SCANCODE_W, false}
	i.keyboard.backward = KeyState{sdl.SCANCODE_S, false}
	i.keyboard.left = KeyState{sdl.SCANCODE_A, false}
	i.keyboard.right = KeyState{sdl.SCANCODE_D, false}
	i.keyboard.up = KeyState{sdl.SCANCODE_Q, false}
	i.keyboard.down = KeyState{sdl.SCANCODE_E, false}
}

func (i *Input) KeyStates() []*KeyState {
return []*KeyState{&i.keyboard.forward,
    &i.keyboard.backward,
		&i.keyboard.left,
    &i.keyboard.right,
    &i.keyboard.up,
    &i.keyboard.down}

}

