package main

import (
	"GameOfLife/src/Game"
	"math/rand/v2"
	"strconv"
	"testing"
)

func BenchmarkUpdateAlgorithms(b *testing.B) {
	// create two equal grids with 100000 cells and random inhibited cells
	g := Game.Grid{}
	grid := g.New(1000, 1000)
	parallelGrid := g.New(1000, 1000)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if rand.IntN(2) == 1 {
				grid.Cells[i][j].Inhibited = true
				parallelGrid.Cells[i][j].Inhibited = true
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
	createAutoBenchmark(100, b)
	createAutoBenchmark(1000, b)
	createAutoBenchmark(2000, b)
	createAutoBenchmark(3000, b)
	createAutoBenchmark(4000, b)
	createAutoBenchmark(5000, b)
	createAutoBenchmark(10000, b)
	createAutoBenchmark(20000, b)
}

func createAutoBenchmark(n int, b *testing.B) {

	// create a game with a grid of 10000x10000 cells
	game2 := Game.Game{}.New(n, n)
	game2.AutoParallelize = false
	b.Run(strconv.Itoa(n)+" AP Off Cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < 4; j++ {
				game2.Cycle()
			}
		}
	})

	game4 := Game.Game{}.New(n, n)
	game4.AutoParallelize = true

	b.Run(strconv.Itoa(n)+" AP On Cycle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < 4; j++ {
				game4.Cycle()
			}
		}
	})
}
