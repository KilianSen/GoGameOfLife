package Game

type Cell struct {
	x         int
	y         int
	Inhibited bool
}

func resolveCell(c Cell, g Grid) bool {
	// apply the rules of the game of life to a cell
	neighbours := g.GetNeighbours(c)
	var inhibitedNeighbours int
	for _, n := range neighbours {
		if n.Inhibited {
			inhibitedNeighbours++
		}
	}

	if c.Inhibited {
		if inhibitedNeighbours < 2 || inhibitedNeighbours > 3 {
			return false
		}
		return true
	} else {
		if inhibitedNeighbours == 3 {
			return true
		}
		return false
	}
}
