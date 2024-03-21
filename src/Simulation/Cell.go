package Simulation

type Cell struct {
	x         int
	y         int
	z         int
	Inhibited bool
}

func resolveCell(c *Cell, g *Grid) bool {
	// apply the rules of the game of life to a cell
	// TODO: check if these rules apply to the 3D version of the game of life
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
