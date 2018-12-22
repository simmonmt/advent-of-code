package rope

type Rope [][]byte

type Cursor struct {
	rope   Rope
	arrIdx int
	off    int
}

func (c *Cursor) Val() byte {
	return c.rope[c.arrIdx][c.off]
}

func (c *Cursor) Advance() bool {
	c.off++
	if c.off >= len(c.rope[c.arrIdx]) {
		c.arrIdx++
		if c.arrIdx >= len(c.rope) {
			return false
		}
	}
	return true
}

func NewRope(vals [][]byte) hRope {
	return Rope(vals)
}

func (r Rope) Begin() *Cursor {
	return &Cursor{r, 0, 0}
}
