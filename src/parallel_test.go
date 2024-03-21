package main

import (
	"GameOfLife/src/Simulation"
	"math/rand/v2"
	"strconv"
	"testing"
)

func BenchmarkUpdateAlgorithms(b *testing.B) {
	// create two equal grids with 100000 cells and random inhibited cells

	x, y, z := 100, 100, 100

	grid := Simulation.NewGrid(x, y, z)
	parallelGrid := Simulation.NewGrid(x, y, z)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			for k := 0; k < z; k++ {
				if rand.IntN(2) == 1 {
					grid.Cells[i][j][k].Inhibited = true
					parallelGrid.Cells[i][j][k].Inhibited = true
				}
			}
		}
	}

	// run the benchmark
	b.Run("Update", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			grid.Update()
		}
	})
	b.Run("ParallelUpdate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parallelGrid.ParallelUpdate()
		}
	})
}

func BenchmarkAutoParallelization(b *testing.B) {
	createAutoBenchmark(10, b)
	createAutoBenchmark(20, b)
	createAutoBenchmark(50, b)
	createAutoBenchmark(100, b)
	createAutoBenchmark(200, b)
	createAutoBenchmark(500, b)
	//createAutoBenchmark(1000, b) to much ram
}

func createAutoBenchmark(n int, b *testing.B) {

	// create a game with a grid of 10000x10000 cells
	game2 := Simulation.NewSimulation(n, n, n)
	game2.AutoParallelize = false
	b.Run(strconv.Itoa(n)+" AP Off Cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < 4; j++ {
				game2.Cycle()
			}
		}
	})

	game4 := Simulation.NewSimulation(n, n, n)
	game4.AutoParallelize = true

	b.Run(strconv.Itoa(n)+" AP On Cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < 4; j++ {
				game4.Cycle()
			}
		}
	})
}
