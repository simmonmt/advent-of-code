package parse

type NodeType int

const (
	TYPE_IMM NodeType = iota
	TYPE_ADD
	TYPE_MULT
	TYPE_EXPR
)

func (t NodeType) String() string {
	switch t {
	case TYPE_IMM:
		return "imm"
	case TYPE_ADD:
		return "add"
	case TYPE_MULT:
		return "mult"
	case TYPE_EXPR:
		return "expr"
	default:
		return "???"
	}
}

type Node struct {
	Type NodeType
	Imm  int
	Expr []*Node
}
