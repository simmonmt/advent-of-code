package main

import (
	"fmt"
	"xyzpos"
)

type SearchCube struct {
	Min  xyzpos.Pos
	Size int
}

func (c SearchCube) String() string {
	return fmt.Sprintf("%d,%d,%d:%d", c.Min.X, c.Min.Y, c.Min.Z, c.Size)
}

func distToRange(val, lo, hi int) int {
	if val > hi {
		return val - hi
	} else if val < lo {
		return lo - val
	}
	return 0
}

func (s *SearchCube) InRange(bot *Bot) bool {
	// The Manhattan distance from the bot's location to a point
	// just inside the cube.
	dist := distToRange(bot.Pos.X, s.Min.X, s.Min.X+s.Size-1)
	dist += distToRange(bot.Pos.Y, s.Min.Y, s.Min.Y+s.Size-1)
	dist += distToRange(bot.Pos.Z, s.Min.Z, s.Min.Z+s.Size-1)

	// If the bot includes the cube in its broadcast radius, the
	// radius can reach the point.
	return dist <= bot.Radius
}

func (s *SearchCube) Divide() []*SearchCube {
	out := make([]*SearchCube, 8)
	i := 0
	newSize := s.Size / 2
	for z := 0; z < 2; z++ {
		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				out[i] = &SearchCube{
					Min:  xyzpos.Pos{s.Min.X + newSize*x, s.Min.Y + newSize*y, s.Min.Z + newSize*z},
					Size: newSize,
				}
				i++
			}
		}
	}
	return out
}

type AllocatedCube struct {
	Cube *SearchCube
	Bots []*Bot
}

func NewAllocatedCube(bots []*Bot, searchCube *SearchCube) *AllocatedCube {
	inRange := []*Bot{}
	for _, bot := range bots {
		if searchCube.InRange(bot) {
			inRange = append(inRange, bot)
		} else {
			// logger.LogF("not in range: %v of %v", bot, *searchCube)
		}
	}

	return &AllocatedCube{
		Cube: searchCube,
		Bots: inRange,
	}
}

func (c *AllocatedCube) Divide() []*AllocatedCube {
	out := make([]*AllocatedCube, 8)
	subs := c.Cube.Divide()
	for i, sub := range subs {
		out[i] = NewAllocatedCube(c.Bots, sub)
	}
	return out
}

type cubeQueueItem struct {
	cube  *AllocatedCube
	index int
}

type CubeQueue []*cubeQueueItem

func (q CubeQueue) Len() int { return len(q) }

func (q CubeQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return len(q[i].cube.Bots) > len(q[j].cube.Bots)
}

func (q CubeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *CubeQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*cubeQueueItem)
	item.index = n
	*q = append(*q, item)
}

func (q *CubeQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}
