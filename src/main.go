package main

import (
	"GameOfLife/src/GUI"
	"GameOfLife/src/Simulation"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math/rand/v2"
	"time"
)

const (
	width  = 200
	height = 200
	depth  = 1
)

func main() {

	sim := Simulation.NewSimulation(width, height, depth)

	sim.AutoParallelize = false
	sim.UseParallel = true

	inhibitPercentage := 0.5
	for x := range sim.Grid.Cells {
		for y := range sim.Grid.Cells[x] {
			for z := range sim.Grid.Cells[x][y] {
				if rand.Float64() < inhibitPercentage {
					sim.Grid.Cells[x][y][z].Inhibited = true
				}
			}
		}
	}

	time.Sleep(1 * time.Second)
	gui := GUI.NewGUI("Game of Life", width*3, height*3)

	// create a new cell array for the GUI
	uiCells := make([][][]bool, width)
	for x := range uiCells {
		uiCells[x] = make([][]bool, height)
		for y := range uiCells[x] {
			uiCells[x][y] = make([]bool, depth)
			for z := range uiCells[x][y] {
				tmp := sim.Grid.Cells[x][y][z].Inhibited
				uiCells[x][y][z] = tmp
			}
		}
	}
	gui.Cells = &uiCells

	for {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(gui.Program)

		sim.Cycle()

		// update the GUI cells from the simulation cells
		for x := range sim.Grid.Cells {
			for y := range sim.Grid.Cells[x] {
				for z := range sim.Grid.Cells[x][y] {
					uiCells[x][y][z] = sim.Grid.Cells[x][y][z].Inhibited
				}
			}
		}

		xInd := 0
		for x := range *gui.Cells {
			yInd := 0
			for _, y := range (*gui.Cells)[x] {
				zInd := 0
				for _, c := range y {
					if c {
						gl.BindVertexArray(GUI.NewVertices(xInd, yInd, width, height))
						gl.DrawArrays(gl.TRIANGLES, 0, int32(len(GUI.Square)/3))
					}
					zInd++
				}
				yInd++
			}
			xInd++
		}

		glfw.PollEvents()
		gui.Window.SwapBuffers()

	}

}
