package reg

import (
	"fmt"
	"sort"
)

type File struct {
	f map[Reg]uint32
}

func NewFile() *File {
	return &File{f: map[Reg]uint32{}}
}

func (f *File) Get(reg Reg) uint32 {
	return f.f[reg]
}

func (f *File) Set(reg Reg, val uint32) {
	f.f[reg] = val
}

func (f *File) Print() {
	names := make([]int, len(f.f))
	i := -1
	for name, _ := range f.f {
		i++
		names[i] = int(name)
	}
	sort.Ints(names)

	for _, name := range names {
		r := Reg(name)
		fmt.Printf("reg[%s] = %v\n", r, f.f[r])
	}
}
