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

package vm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

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

func NewRamFromReader(r io.Reader) (Ram, error) {
	var line string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	ram := NewRam()
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram.Write(int64(i), int64(val))
	}

	return ram, nil
}

func NewRamFromFile(path string) (Ram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewRamFromReader(f)
}

type roRam struct {
	Ram
}

func (r *roRam) Write(addr, val int64) {
	panic("write to r/o ram")
}

func (r *roRam) Clone() Ram {
	return r
}

func NewReadOnlyRam(ram Ram) Ram {
	return &roRam{ram}
}
