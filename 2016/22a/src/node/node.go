package node

import "fmt"

type Node struct {
	Used, Size uint16
}

func New(size, used uint16) *Node {
	return &Node{
		Used: used,
		Size: size,
	}
}

func (n Node) String() string {
	return fmt.Sprintf("%v/%v", n.Used, n.Size)
}

func (n Node) Avail() uint16 {
	return n.Size - n.Used
}
