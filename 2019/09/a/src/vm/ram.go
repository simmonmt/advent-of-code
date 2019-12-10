package vm

import "fmt"

type Ram interface {
	Read(addr int) int
	Write(addr, val int)
	Clone() Ram
	Dump()
}

type ramImpl struct {
	r map[int]int
}

func NewRam(vals ...int) Ram {
	r := &ramImpl{
		r: map[int]int{},
	}

	for a, v := range vals {
		r.r[a] = v
	}

	return r
}

func (r *ramImpl) Read(addr int) int {
	return r.r[addr]
}

func (r *ramImpl) Write(addr, val int) {
	r.r[addr] = val
}

func (r *ramImpl) Clone() Ram {
	nr := NewRam().(*ramImpl)
	for a, v := range r.r {
		nr.r[a] = v
	}
	return nr
}

func (r *ramImpl) Dump() {
	fmt.Println(r.r)
}
