package instr

import (
	"reg"
)

func Addr(file *reg.File, a, b, c int) {
	file[c] = file[a] + file[b]
}

func Addi(file *reg.File, a, b, c int) {
	file[c] = file[a] + b
}

func Mulr(file *reg.File, a, b, c int) {
	file[c] = file[a] * file[b]
}

func Muli(file *reg.File, a, b, c int) {
	file[c] = file[a] * b
}

func Banr(file *reg.File, a, b, c int) {
	file[c] = file[a] & file[b]
}

func Bani(file *reg.File, a, b, c int) {
	file[c] = file[a] & b
}

func Borr(file *reg.File, a, b, c int) {
	file[c] = file[a] | file[b]
}

func Bori(file *reg.File, a, b, c int) {
	file[c] = file[a] | b
}

func Setr(file *reg.File, a, b, c int) {
	file[c] = file[a]
}

func Seti(file *reg.File, a, b, c int) {
	file[c] = a
}

func Gtir(file *reg.File, a, b, c int) {
	if a > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtri(file *reg.File, a, b, c int) {
	if file[a] > b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtrr(file *reg.File, a, b, c int) {
	if file[a] > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqir(file *reg.File, a, b, c int) {
	if a == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqri(file *reg.File, a, b, c int) {
	if file[a] == b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqrr(file *reg.File, a, b, c int) {
	if file[a] == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

type Desc struct {
	Op   int
	F    func(file *reg.File, a, b, c int)
	Name string
}

var (
	All = []Desc{
		{5, Addr, "addr"},
		{14, Addi, "addi"},
		{3, Mulr, "mulr"},
		{10, Muli, "muli"},
		{12, Banr, "banr"},
		{9, Bani, "bani"},
		{1, Borr, "borr"},
		{0, Bori, "bori"},
		{4, Setr, "setr"},
		{2, Seti, "seti"},
		{6, Gtir, "gtir"},
		{8, Gtri, "gtri"},
		{11, Gtrr, "gtrr"},
		{7, Eqir, "eqir"},
		{13, Eqri, "eqri"},
		{15, Eqrr, "eqrr"},
	}
)
