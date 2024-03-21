package Game

import "time"

type Game struct {
	Grid        *Grid
	UseParallel bool

	AutoParallelize     bool
	nonParallelDuration int64
	parallelDuration    int64
}

func NewGame(x int, y int) *Game {
	game := &Game{Grid: NewGrid(x, y), UseParallel: false, AutoParallelize: true, nonParallelDuration: 0, parallelDuration: 0}
	return game
}

func (g *Game) Cycle() {
	if g.AutoParallelize {
		if g.nonParallelDuration == 0 {
			g.UseParallel = false
		} else if g.parallelDuration == 0 {
			g.UseParallel = true
		} else if g.nonParallelDuration < g.parallelDuration {
			g.UseParallel = false
			g.AutoParallelize = false
		} else {
			g.UseParallel = true
			g.AutoParallelize = false
		}
	}

	startTime := time.Now()
	if g.UseParallel {
		g.Grid.ParallelUpdate()
	} else {
		g.Grid.Update()
	}
	duration := time.Now().Sub(startTime).Nanoseconds()

	if g.AutoParallelize {
		if g.UseParallel {
			g.parallelDuration = duration
		} else {
			g.nonParallelDuration = duration
		}
	}
}
