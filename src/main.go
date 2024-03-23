package main

import (
	"GameOfLife/src/GUI"
	"GameOfLife/src/Simulation"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {

	sim := Simulation.NewSimulation(10, 10, 0)
	gui := GUI.NewGUI()

	for {
		sim.Cycle()

		// get all cells with [x][y][0]
		var sGrid = sim.Grid[0:10][0:10][0]

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(*g.program)

		gui.
		for x := range cells {
			for _, y := range cells[x] {
				for _, c := range y {
					if c.Inhibited {
						gl.BindVertexArray(c.drawable)
						gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
					}
				}
			}
		}

		glfw.PollEvents()
		g.window.SwapBuffers()

	}

}
