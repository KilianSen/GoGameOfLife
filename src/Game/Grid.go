package Game

import (
	"runtime"
	"sync"
)

type Grid struct {
	Cells [][]Cell
	pool  sync.Pool
	done  chan bool
}

func (g *Grid) New(x int, y int) *Grid {
	// initialize the grid with cells all set to not inhibited
	var cells = make([][]Cell, x)
	for i := range cells {
		cells[i] = make([]Cell, y)
		for j := range cells[i] {
			cells[i][j] = Cell{x: i, y: j, Inhibited: false}
		}
	}
	return &Grid{
		Cells: cells,
		pool: sync.Pool{
			New: func() interface{} {
				return make([][]Cell, x)
			},
		},
		done: make(chan bool, runtime.NumCPU()), // Buffered channel for synchronization
	}
}

func (g *Grid) GetNeighbours(c *Cell) []Cell {
	// get the neighbours of a cell
	var neighbours []Cell
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if c.x+i >= 0 && c.x+i < len(g.Cells) && c.y+j >= 0 && c.y+j < len(g.Cells[0]) {
				neighbours = append(neighbours, g.Cells[c.x+i][c.y+j])
			}
		}
	}
	return neighbours
}

func (g *Grid) Print() {
	// print the grid
	for i := 0; i < len(g.Cells); i++ {
		for j := 0; j < len(g.Cells[0]); j++ {
			if g.Cells[i][j].Inhibited {
				print("X")
			} else {
				print("O")
			}
		}
		println()
	}
}

func (g *Grid) Update() {
	// update the grid according to the rules of the game of life
	var newCells = make([][]Cell, len(g.Cells))
	for i := range newCells {
		newCells[i] = make([]Cell, len(g.Cells[0]))
		for j := range newCells[i] {
			newCells[i][j] = Cell{x: i, y: j, Inhibited: resolveCell(&g.Cells[i][j], g)}
		}
	}
	g.Cells = newCells
}

func (g *Grid) ParallelUpdate() {
	newCells := g.pool.Get().([][]Cell) // Get a slice from the pool

	// Divide the grid into sections
	sections := runtime.NumCPU()
	sectionSize := len(g.Cells) / sections

	for i := 0; i < sections; i++ {
		// Calculate the start and end indices for this section
		start := i * sectionSize
		end := start + sectionSize
		if i == sections-1 {
			end = len(g.Cells) // Make sure the last section goes to the end of the grid
		}

		// Launch a goroutine to handle this section
		go func(start, end int) {
			for i := start; i < end; i++ {
				newCells[i] = make([]Cell, len(g.Cells[0]))
				for j := range newCells[i] {
					newCells[i][j] = Cell{x: i, y: j, Inhibited: resolveCell(&g.Cells[i][j], g)}
				}
			}
			g.done <- true // Signal that this goroutine is done
		}(start, end)
	}

	// Wait for all goroutines to finish
	for i := 0; i < sections; i++ {
		<-g.done
	}

	// Update the grid
	g.Cells = newCells

	// Put the slice back in the pool for reuse
	g.pool.Put(newCells)
}
