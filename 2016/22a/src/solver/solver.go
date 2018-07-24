package solver

type aStarHelper struct {
	width, height int
}

func NewHelper(width, height int) *aStarHelper {
	return &aStarHelper{
		width:  width,
		height: height,
	}
}

func findNewGoals(cur *grid.Grid, curGoalX, curGoalY uint8) []*grid.Grid {
	grids := []*grid.Grid{}
	grids = append(grids, cur.MoveGoal(curGoalX-1, curGoalY))
	grids = append(grids, cur.MoveGoal(curGoalX, curGoalY-1))
	grids = append(grids, cur.MoveGoal(curGoalX+1, curGoalY))
	grids = append(grids, cur.MoveGoal(curGoalX, curGoalY+1))
	return grids
}

func findRecipients(cur *grid.Grid, x, y uint8) []*grid.Grid {
	grids := []*grid.Grid{}
	curNode := cur.Get(x, y)

	if n := cur.Get(x-1, y); n != nil && n.Avail() > curNode.Size {
		grids = append(grids, grid.Transfer(cur, n))
	}
	if n := cur.Get(x+1, y); n != nil && n.Avail() > curNode.Size {
		grids = append(grids, grid.Transfer(cur, n))
	}
	if n := cur.Get(x, y-1); n != nil && n.Avail() > curNode.Size {
		grids = append(grids, grid.Transfer(cur, n))
	}
	if n := cur.Get(x, y+1); n != nil && n.Avail() > curNode.Size {
		grids = append(grids, grid.Transfer(cur, n))
	}

	return grids
}

func (h *aStarHelper) AllNeighbors(curStr string) []string {
	cur, err := grid.Deserialize(w, h, []byte(curStr))
	if err != nil {
		panic("failed to deserialize")
	}

	neighbors := []*grid.Grid{}

	if h.goalX != 0 && h.goalY != 0 {
		neighbors = append(neighbors, h.findNewGoals(cur, h.goalX, h.goalY))
	}

	for y := 0; y < h.height; y++ {
		for x := 0; x < h.width; y++ {
			neighbors = append(neighbors, findRecipients(cur, x, y))
		}
	}

	nStrs := make([]string, len(neighbors))
	for i, n := range neighbors {
		nStrs[i] = string(n.Serialize())
	}
	return nStrs
}
