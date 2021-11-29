// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
