package graph

type NodeID string

type Graph interface {
	Neighbors(id NodeID) []NodeID
	NeighborDistance(from, to NodeID) int
}
