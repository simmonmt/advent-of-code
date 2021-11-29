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

package instr

import (
	"reg"
)

func Addr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] + file[b]
}

func Addi(file *reg.File, op, a, b, c int) {
	file[c] = file[a] + b
}

func Mulr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] * file[b]
}

func Muli(file *reg.File, op, a, b, c int) {
	file[c] = file[a] * b
}

func Banr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] & file[b]
}

func Bani(file *reg.File, op, a, b, c int) {
	file[c] = file[a] & b
}

func Borr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] | file[b]
}

func Bori(file *reg.File, op, a, b, c int) {
	file[c] = file[a] | b
}

func Setr(file *reg.File, op, a, b, c int) {
	file[c] = file[a]
}

func Seti(file *reg.File, op, a, b, c int) {
	file[c] = a
}

func Gtir(file *reg.File, op, a, b, c int) {
	if a > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtri(file *reg.File, op, a, b, c int) {
	if file[a] > b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtrr(file *reg.File, op, a, b, c int) {
	if file[a] > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqir(file *reg.File, op, a, b, c int) {
	if a == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqri(file *reg.File, op, a, b, c int) {
	if file[a] == b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqrr(file *reg.File, op, a, b, c int) {
	if file[a] == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

type Desc struct {
	F    func(file *reg.File, op, a, b, c int)
	Name string
}

var (
	All = []Desc{
		{Addr, "Addr"},
		{Addi, "Addi"},
		{Mulr, "Mulr"},
		{Muli, "Muli"},
		{Banr, "Banr"},
		{Bani, "Bani"},
		{Borr, "Borr"},
		{Bori, "Bori"},
		{Setr, "Setr"},
		{Seti, "Seti"},
		{Gtir, "Gtir"},
		{Gtri, "Gtri"},
		{Gtrr, "Gtrr"},
		{Eqir, "Eqir"},
		{Eqri, "Eqri"},
		{Eqrr, "Eqrr"},
	}
)
