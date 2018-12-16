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
	Op   int
	F    func(file *reg.File, op, a, b, c int)
	Name string
}

var (
	All = []Desc{
		{5, Addr, "Addr"},
		{14, Addi, "Addi"},
		{3, Mulr, "Mulr"},
		{10, Muli, "Muli"},
		{12, Banr, "Banr"},
		{9, Bani, "Bani"},
		{1, Borr, "Borr"},
		{0, Bori, "Bori"},
		{4, Setr, "Setr"},
		{2, Seti, "Seti"},
		{6, Gtir, "Gtir"},
		{8, Gtri, "Gtri"},
		{11, Gtrr, "Gtrr"},
		{7, Eqir, "Eqir"},
		{13, Eqri, "Eqri"},
		{15, Eqrr, "Eqrr"},
	}
)
