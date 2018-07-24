package node

import "fmt"

type Node struct {
	Size, Used uint16
}

func New(size, used uint16) *Node {
	return &Node{
		Size: size,
		Used: used,
	}
}

func (n Node) String() string {
	return fmt.Sprintf("%v/%v", n.Used, n.Size)
}
