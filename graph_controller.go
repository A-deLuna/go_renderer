package main

// import "fmt"

type GraphController struct {
	graphStateChannels []chan<- GraphState
	commandChannels    []chan<- *DrawCommand
}

func NewGraphController(tileCount int32) (gc GraphController) {
	gc.graphStateChannels = make([]chan<- GraphState, 0, tileCount)
	gc.commandChannels = make([]chan<- *DrawCommand, 0, tileCount)
	return gc
}

func (gc *GraphController) AddTileChans(gs chan GraphState, dc chan *DrawCommand) {
	gc.graphStateChannels = append(gc.graphStateChannels, gs)
	gc.commandChannels = append(gc.commandChannels, dc)
}

func (gc *GraphController) NotifyTile(id int32, dc *DrawCommand) {
	gc.commandChannels[id] <- dc
}

func (gc *GraphController) ChangeState(state GraphState) {
	//fmt.Printf("Notifying %d channels of state change\n", len(gc.graphStateChannels))
	for _, channel := range gc.graphStateChannels {
		channel <- state
	}
}
