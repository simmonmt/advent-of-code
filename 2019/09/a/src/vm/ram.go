package vm

import "fmt"

type Ram interface {
	Read(addr int64) int64
	Write(addr, val int64)
	Clone() Ram
	Dump()
}

type ramImpl struct {
	r map[int64]int64
}

func NewRam(vals ...int64) Ram {
	r := &ramImpl{
		r: map[int64]int64{},
	}

	for a, v := range vals {
		r.r[int64(a)] = v
	}

	return r
}

func (r *ramImpl) Read(addr int64) int64 {
	return r.r[addr]
}

func (r *ramImpl) Write(addr, val int64) {
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
